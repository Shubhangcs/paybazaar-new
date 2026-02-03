package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) GetAllDTHOperatorsQuery(
	ctx context.Context,
) ([]models.GetDTHOperatorsResponseModel, error) {
	query := `
		SELECT 
			operator_code, 
			operator_name
		FROM dth_recharge_operators;
	`
	res, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var operators []models.GetDTHOperatorsResponseModel
	for res.Next() {
		var operator models.GetDTHOperatorsResponseModel
		if err := res.Scan(
			&operator.OperatorCode,
			&operator.OperatorName,
		); err != nil {
			return nil, err
		}
		operators = append(operators, operator)
	}
	return operators, res.Err()
}

func (db *Database) CreateDTHRechargeSuccessOrPendingQuery(
	ctx context.Context,
	req models.CreateDTHRechargeRequestModel,
) error {
	if req.Amount <= 99 {
		return db.dthRechargeWithoutCommision(ctx, req)
	}
	return db.dthRechargeWithCommision(ctx, req)
}

func (db *Database) CreateDTHRechargeFailedQuery(
	ctx context.Context,
	req models.CreateDTHRechargeRequestModel,
) error {
	insertToDthRechargeTable := `
		INSERT INTO dth_recharge (
			retailer_id,
			partner_request_id,
			customer_id,
			operator_name,
			operator_code,
			amount,
			commision,
			status
		) VALUES (
			@retailer_id,
			@partner_request_id,
			@customer_id,
			@operator_name,
			@operator_code,
			@amount,
			@commision,
			@status
		)
		RETURNING dth_transaction_id AS transaction_id;
	`
	if _, err := db.pool.Exec(ctx, insertToDthRechargeTable, pgx.NamedArgs{
		"retailer_id":        req.RetailerID,
		"partner_request_id": req.PartnerRequestID,
		"customer_id":        req.CustomerID,
		"operator_name":      req.OperatorName,
		"operator_code":      req.OperatorCode,
		"amount":             req.Amount,
		"commision":          0,
		"status":             req.Status,
	}); err != nil {
		return err
	}
	return nil
}

func (db *Database) dthRechargeWithoutCommision(
	ctx context.Context,
	req models.CreateDTHRechargeRequestModel,
) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getRetailerBeforeBalanceQuery := `
		SELECT retailer_wallet_balance AS retailer_before_balance
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`
	var retailerBeforeBalance float64
	if err := tx.QueryRow(ctx, getRetailerBeforeBalanceQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(&retailerBeforeBalance); err != nil {
		return err
	}

	deductRetailerAmountAndGetAfterBalance := `
		UPDATE retailers
		SET retailer_wallet_balance = retailer_wallet_balance - @amount
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance AS retailer_after_balance;
	`
	var retailerAfterBalance float64
	if err := tx.QueryRow(ctx, deductRetailerAmountAndGetAfterBalance, pgx.NamedArgs{
		"amount":      req.Amount,
		"retailer_id": req.RetailerID,
	}).Scan(&retailerAfterBalance); err != nil {
		return err
	}

	insertToDthRechargeTable := `
		INSERT INTO dth_recharge (
			retailer_id,
			partner_request_id,
			customer_id,
			operator_name,
			operator_code,
			amount,
			commision,
			status
		) VALUES (
			@retailer_id,
			@partner_request_id,
			@customer_id,
			@operator_name,
			@operator_code,
			@amount,
			@commision,
			@status
		)
		RETURNING dth_transaction_id AS transaction_id;
	`
	var transactionID int
	if err := tx.QueryRow(ctx, insertToDthRechargeTable, pgx.NamedArgs{
		"retailer_id":        req.RetailerID,
		"partner_request_id": req.PartnerRequestID,
		"customer_id":        req.CustomerID,
		"operator_name":      req.OperatorName,
		"operator_code":      req.OperatorCode,
		"amount":             req.Amount,
		"commision":          0,
		"status":             req.Status,
	}).Scan(&transactionID); err != nil {
		return err
	}

	insertToRetailerWalletTransactionsTable := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			debit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@debit_amount,
			@before_balance,
			@after_balance,
			@transaction_reason,
			@remarks
		);
	`
	if _, err := tx.Exec(ctx, insertToRetailerWalletTransactionsTable, pgx.NamedArgs{
		"user_id":            req.RetailerID,
		"reference_id":       fmt.Sprintf("%d",transactionID),
		"debit_amount":       req.Amount,
		"before_balance":     retailerBeforeBalance,
		"after_balance":      retailerAfterBalance,
		"transaction_reason": "DTH_RECHARGE",
		"remarks":            fmt.Sprintf("DTH Recharge to: %s", req.CustomerID),
	}); err != nil {
		return err
	}

	// CRITICAL: Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (db *Database) dthRechargeWithCommision(
	ctx context.Context,
	req models.CreateDTHRechargeRequestModel,
) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getAdminBeforeBalanceQuery := `
		SELECT ad.admin_wallet_balance, ad.admin_id
		FROM retailers AS r
		LEFT JOIN distributors AS d
    		ON r.distributor_id = d.distributor_id
		LEFT JOIN master_distributors AS md
    		ON d.master_distributor_id = md.master_distributor_id
		LEFT JOIN admins AS ad
    		ON md.admin_id = ad.admin_id
		WHERE r.retailer_id = @retailer_id;
	`
	var adminBeforeBalance float64
	var adminID string
	if err := tx.QueryRow(ctx, getAdminBeforeBalanceQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(&adminBeforeBalance, &adminID); err != nil {
		return err
	}

	// Deduct commission from admin wallet
	deductAdminAmountAndGetAfterBalanceQuery := `
		UPDATE admins AS ad
		SET admin_wallet_balance = admin_wallet_balance - @commision
		FROM retailers AS r
		LEFT JOIN distributors AS d
    		ON r.distributor_id = d.distributor_id
		LEFT JOIN master_distributors AS md
    		ON d.master_distributor_id = md.master_distributor_id
		WHERE r.retailer_id = @retailer_id
    		AND md.admin_id = ad.admin_id
		RETURNING ad.admin_wallet_balance AS admin_after_balance;
	`
	var adminAfterBalance float64
	if err := tx.QueryRow(ctx, deductAdminAmountAndGetAfterBalanceQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
		"commision":   1,
	}).Scan(&adminAfterBalance); err != nil {
		return err
	}

	getRetailerBeforeBalanceQuery := `
		SELECT retailer_wallet_balance AS retailer_before_balance
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`
	var retailerBeforeBalance float64
	if err := tx.QueryRow(ctx, getRetailerBeforeBalanceQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(&retailerBeforeBalance); err != nil {
		return err
	}

	// Deduct (amount - commission) from retailer
	// Retailer pays less because they get commission benefit
	deductAmountFromRetailerQuery := `
		UPDATE retailers 
		SET retailer_wallet_balance = retailer_wallet_balance - @amount
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance AS retailer_after_balance;
	`
	var retailerAfterBalance float64
	if err := tx.QueryRow(ctx, deductAmountFromRetailerQuery, pgx.NamedArgs{
		"amount":      req.Amount - 1, // Amount minus commission
		"retailer_id": req.RetailerID,
	}).Scan(&retailerAfterBalance); err != nil {
		return err
	}

	insertToDthRechargeTableQuery := `
		INSERT INTO dth_recharge (
			retailer_id,
			partner_request_id,
			customer_id,
			operator_name,
			operator_code,
			amount,
			commision,
			status
		) VALUES (
			@retailer_id,
			@partner_request_id,
			@customer_id,
			@operator_name,
			@operator_code,
			@amount,
			@commision,
			@status
		)
		RETURNING dth_transaction_id AS transaction_id;
	`
	var transactionID int
	if err := tx.QueryRow(ctx, insertToDthRechargeTableQuery, pgx.NamedArgs{
		"retailer_id":        req.RetailerID,
		"partner_request_id": req.PartnerRequestID,
		"customer_id":        req.CustomerID,
		"operator_name":      req.OperatorName,
		"operator_code":      req.OperatorCode,
		"amount":             req.Amount,
		"commision":          1,
		"status":             req.Status,
	}).Scan(&transactionID); err != nil {
		return err
	}

	// Record admin's commission deduction
	insertToAdminWalletTransactionsQuery := `
		INSERT INTO wallet_transactions (
    		user_id,
    		reference_id,
    		debit_amount,
    		before_balance,
    		after_balance,
    		transaction_reason,
    		remarks
		) VALUES (
    		@user_id,
    		@reference_id,
    		@debit_amount,
    		@before_balance,
    		@after_balance,
    		@transaction_reason,
    		@remarks
		);
	`
	if _, err := tx.Exec(ctx, insertToAdminWalletTransactionsQuery, pgx.NamedArgs{
		"user_id":            adminID,
		"reference_id":       fmt.Sprintf("%d",transactionID),
		"debit_amount":       1,
		"before_balance":     adminBeforeBalance,
		"after_balance":      adminAfterBalance,
		"transaction_reason": "DTH_RECHARGE",
		"remarks":            fmt.Sprintf("Commission for Retailer: %s", req.RetailerID),
	}); err != nil {
		return err
	}

	// Record retailer's transaction (debit for recharge amount minus commission received)
	insertToRetailerWalletTransactionsQuery := `
		INSERT INTO wallet_transactions (
    		user_id,
    		reference_id,
    		debit_amount,
    		before_balance,
    		after_balance,
    		transaction_reason,
    		remarks
		) VALUES (
    		@user_id,
    		@reference_id,
    		@debit_amount,
    		@before_balance,
    		@after_balance,
    		@transaction_reason,
    		@remarks
		);
	`
	if _, err := tx.Exec(ctx, insertToRetailerWalletTransactionsQuery, pgx.NamedArgs{
		"user_id":            req.RetailerID,
		"reference_id":       fmt.Sprintf("%d",transactionID),
		"debit_amount":       req.Amount - 1, // Net amount paid by retailer
		"before_balance":     retailerBeforeBalance,
		"after_balance":      retailerAfterBalance,
		"transaction_reason": "DTH_RECHARGE",
		"remarks":            fmt.Sprintf("DTH Recharge to: %s (Commission: ₹1)", req.CustomerID),
	}); err != nil {
		return err
	}

	// CRITICAL: Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (db *Database) GetAllDTHRechargesQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetDTHRechargeHistoryResponseModel, error) {
	query := `
		SELECT 
			d.dth_transaction_id,
			d.retailer_id,
    		d.partner_request_id,
    		d.customer_id,
    		d.operator_name,
    		d.operator_code,
    		d.amount,
    		d.commision,
			d.status,
			w.before_balance,
			w.after_balance,
			d.created_at
		FROM dth_recharge d
		JOIN wallet_transactions w
			ON w.reference_id = d.dth_transaction_id
			AND w.retailer_id = d.retailer_id
			AND w.transaction_reason = 'MOBILE_RECHARGE'
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`
	res, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"offset": offset,
		"limit":  limit,
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var history []models.GetDTHRechargeHistoryResponseModel
	for res.Next() {
		var recharge models.GetDTHRechargeHistoryResponseModel
		if err := res.Scan(
			&recharge.DTHTransactionID,
			&recharge.RetailerID,
			&recharge.PartnerRequestID,
			&recharge.CustomerID,
			&recharge.OperatorName,
			&recharge.OperatorCode,
			&recharge.Amount,
			&recharge.Commision,
			&recharge.Status,
			&recharge.CreatedAt,
		); err != nil {
			return nil, err
		}
		if recharge.Status == "PENDING" {
			newStatus, err := db.DTHRechargeStatusCheck(recharge.PartnerRequestID)
			if err != nil {
				return nil, err
			}
			if newStatus != "PENDING" {
				if err := db.UpdateDTHRechargeStatus(ctx, newStatus, recharge.DTHTransactionID); err != nil {
					return nil, err
				}
				recharge.Status = newStatus
			}
		}
		history = append(history, recharge)
	}
	return history, res.Err()
}

func (db *Database) GetDTHRechargesByRetailerIDQuery(
	ctx context.Context,
	retailerId string,
	limit, offset int,
) ([]models.GetDTHRechargeHistoryResponseModel, error) {
	query := `
		SELECT 
			dth_transaction_id,
			retailer_id,
    		partner_request_id,
    		customer_id,
    		operator_name,
    		operator_code,
    		amount,
    		commision,
			status,
			created_at
		FROM dth_recharge
		WHERE retailer_id = @retailer_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`
	res, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerId,
		"limit":       limit,
		"offset":      offset,
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var history []models.GetDTHRechargeHistoryResponseModel
	for res.Next() {
		var recharge models.GetDTHRechargeHistoryResponseModel
		if err := res.Scan(
			&recharge.DTHTransactionID,
			&recharge.RetailerID,
			&recharge.PartnerRequestID,
			&recharge.CustomerID,
			&recharge.OperatorName,
			&recharge.OperatorCode,
			&recharge.Amount,
			&recharge.Commision,
			&recharge.Status,
			&recharge.CreatedAt,
		); err != nil {
			return nil, err
		}
		if recharge.Status == "PENDING" {
			newStatus, err := db.DTHRechargeStatusCheck(recharge.PartnerRequestID)
			if err != nil {
				return nil, err
			}
			if newStatus != "PENDING" {
				if err := db.UpdateDTHRechargeStatus(ctx, newStatus, recharge.DTHTransactionID); err != nil {
					return nil, err
				}
				recharge.Status = newStatus
			}
		}
		history = append(history, recharge)
	}
	return history, res.Err()
}

func (db *Database) UpdateDTHRechargeStatus(
	ctx context.Context,
	status string,
	transactionID int,
) error {
	query := `
		UPDATE dth_recharge 
		SET status = @status
		WHERE dth_transaction_id = @transaction_id;
	`
	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"status":         status,
		"transaction_id": transactionID,
	}); err != nil {
		return err
	}
	return nil
}