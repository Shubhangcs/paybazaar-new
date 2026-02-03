package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) GetPayoutCommisionQuery(
	ctx context.Context,
	retailerId string,
	amount float64,
) (*models.GetPayoutCommisionModel, error) {

	var (
		adminId             string
		masterDistributorId string
		distributorId       string
	)

	hierarchyQuery := `
		SELECT a.admin_id, md.master_distributor_id, d.distributor_id
		FROM retailers r
		JOIN distributors d ON d.distributor_id = r.distributor_id
		JOIN master_distributors md ON md.master_distributor_id = d.master_distributor_id
		JOIN admins a ON a.admin_id = md.admin_id
		WHERE r.retailer_id = @retailer_id;
	`

	if err := db.pool.QueryRow(ctx, hierarchyQuery, pgx.NamedArgs{
		"retailer_id": retailerId,
	}).Scan(&adminId, &masterDistributorId, &distributorId); err != nil {
		return nil, err
	}

	getCommission := func(userId string) (*models.GetPayoutCommisionModel, error) {
		query := `
			SELECT 
				total_commision,
				admin_commision,
				master_distributor_commision,
				distributor_commision,
				retailer_commision
			FROM commisions
			WHERE user_id=@user_id AND service='PAYOUT'
			LIMIT 1;
		`

		var c models.GetPayoutCommisionModel
		err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
			"user_id": userId,
		}).Scan(
			&c.TotalCommision,
			&c.AdminCommision,
			&c.MasterDistributorCommision,
			&c.DistributorCommision,
			&c.RetailerCommision,
		)

		if err == pgx.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return &c, nil
	}

	var commission *models.GetPayoutCommisionModel

	ids := []string{
		retailerId,
		distributorId,
		masterDistributorId,
	}

	for _, id := range ids {
		c, err := getCommission(id)
		if err != nil {
			return nil, err
		}
		if c != nil {
			commission = c
			break
		}
	}

	// Default commission if nothing found
	if commission == nil {
		commission = &models.GetPayoutCommisionModel{
			TotalCommision:             1.2, // %
			RetailerCommision:          0.5, // % of total
			DistributorCommision:       0.2,
			MasterDistributorCommision: 0.05,
			AdminCommision:             0.25,
		}
	}

	// Final calculation (percentage → amount)
	totalAmount := (commission.TotalCommision / 100) * amount

	return &models.GetPayoutCommisionModel{
		TotalCommision:             totalAmount,
		RetailerCommision:          totalAmount * commission.RetailerCommision,
		DistributorCommision:       totalAmount * commission.DistributorCommision,
		MasterDistributorCommision: totalAmount * commission.MasterDistributorCommision,
		AdminCommision:             totalAmount * commission.AdminCommision,
	}, nil
}

func (db *Database) VerifyRetailerForPayoutTransactionQuery(
	ctx context.Context,
	retailerId string,
	minWallerBalance float64,
) error {
	var (
		retailerWalletBalance float64
		retailerKYCStatus     bool
		retailerBlockStatus   bool
	)
	query := `
		SELECT retailer_wallet_balance, retailer_kyc_status,
		is_retailer_blocked
		FROM retailers
		WHERE retailer_id=@retailer_id;
	`
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerId,
	}).Scan(
		&retailerWalletBalance,
		&retailerKYCStatus,
		&retailerBlockStatus,
	); err != nil {
		return err
	}

	if retailerWalletBalance < minWallerBalance {
		return fmt.Errorf("insufficient wallet balance")
	}

	if retailerKYCStatus == false {
		return fmt.Errorf("retailer kyc is pending")
	}

	if retailerBlockStatus == true {
		return fmt.Errorf("retailer is blocked")
	}
	return nil
}

func (db *Database) CreatePayoutSuccessOrPendingQuery(
	ctx context.Context,
	req models.CreatePayoutRequestModel,
	commision models.GetPayoutCommisionModel,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userDetails struct {
		adminId               string
		adminBeforeBalance    float64
		adminAfterBalance     float64
		mdId                  string
		mdBeforeBalance       float64
		mdAfterBalance        float64
		disId                 string
		disBeforeBalance      float64
		disAfterBalance       float64
		retailerBeforeBalance float64
		retailerAfterBalance  float64
		transactionId         string
	}

	// 1️⃣ Insert payout transaction
	insertToPayoutTransactionQuery := `
		INSERT INTO payout_transactions (
			partner_request_id,
			operator_transaction_id,
			order_id,
			retailer_id,
			mobile_number,
			bank_name,
			beneficiary_name,
			account_number,
			ifsc_code,
			amount,
			transfer_type,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			payout_transaction_status
		) VALUES (
			@partner_request_id,
			@operator_transaction_id,
			@order_id,
			@retailer_id,
			@mobile_number,
			@bank_name,
			@beneficiary_name,
			@account_number,
			@ifsc_code,
			@amount,
			@transfer_type,
			@admin_commision,
			@md_commision,
			@dis_commision,
			@retailer_commision,
			@status
		)
		RETURNING payout_transaction_id::TEXT;
	`

	var transferType string
	if req.TransferType == 5 {
		transferType = "IMPS"
	}

	if req.TransferType == 6 {
		transferType = "NEFT"
	}

	if err := tx.QueryRow(ctx, insertToPayoutTransactionQuery, pgx.NamedArgs{
		"partner_request_id":      req.PartnerRequestId,
		"operator_transaction_id": req.OperatorTransactionId,
		"order_id":                req.OrderId,
		"retailer_id":             req.RetailerId,
		"mobile_number":           req.MobileNumber,
		"bank_name":               req.BankName,
		"beneficiary_name":        req.BeneficiaryName,
		"account_number":          req.AccountNumber,
		"ifsc_code":               req.IFSCCode,
		"amount":                  req.Amount,
		"transfer_type":           transferType,
		"admin_commision":         commision.AdminCommision,
		"md_commision":            commision.MasterDistributorCommision,
		"dis_commision":           commision.DistributorCommision,
		"retailer_commision":      commision.RetailerCommision,
		"status":                  req.TransactionStatus,
	}).Scan(&userDetails.transactionId); err != nil {
		return err
	}

	// 2️⃣ Fetch balances & hierarchy
	getUserDetailsQuery := `
		SELECT 
			r.retailer_wallet_balance,
			d.distributor_id,
			d.distributor_wallet_balance,
			m.master_distributor_id,
			m.master_distributor_wallet_balance,
			a.admin_id,
			a.admin_wallet_balance
		FROM retailers r
		JOIN distributors d ON d.distributor_id = r.distributor_id
		JOIN master_distributors m ON m.master_distributor_id = d.master_distributor_id
		JOIN admins a ON a.admin_id = m.admin_id
		WHERE r.retailer_id = @retailer_id;
	`

	if err := tx.QueryRow(ctx, getUserDetailsQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerId,
	}).Scan(
		&userDetails.retailerBeforeBalance,
		&userDetails.disId,
		&userDetails.disBeforeBalance,
		&userDetails.mdId,
		&userDetails.mdBeforeBalance,
		&userDetails.adminId,
		&userDetails.adminBeforeBalance,
	); err != nil {
		return err
	}

	// 3️⃣ Update admin wallet
	updateAdminWalletBalanceQuery := `
		UPDATE admins
		SET admin_wallet_balance = admin_wallet_balance + @commision
		WHERE admin_id = @admin_id
		RETURNING admin_wallet_balance;
	`

	if err := tx.QueryRow(ctx, updateAdminWalletBalanceQuery, pgx.NamedArgs{
		"admin_id":  userDetails.adminId,
		"commision": commision.AdminCommision,
	}).Scan(&userDetails.adminAfterBalance); err != nil {
		return err
	}

	// 4️⃣ Update master distributor wallet
	updateMasterDistributorWalletBalanceQuery := `
		UPDATE master_distributors
		SET master_distributor_wallet_balance =
			master_distributor_wallet_balance + @commision
		WHERE master_distributor_id = @md_id
		RETURNING master_distributor_wallet_balance;
	`

	if err := tx.QueryRow(ctx, updateMasterDistributorWalletBalanceQuery, pgx.NamedArgs{
		"md_id":     userDetails.mdId,
		"commision": commision.MasterDistributorCommision,
	}).Scan(&userDetails.mdAfterBalance); err != nil {
		return err
	}

	// 5️⃣ Update distributor wallet
	updateDistributorWalletBalanceQuery := `
		UPDATE distributors
		SET distributor_wallet_balance =
			distributor_wallet_balance + @commision
		WHERE distributor_id = @dis_id
		RETURNING distributor_wallet_balance;
	`

	if err := tx.QueryRow(ctx, updateDistributorWalletBalanceQuery, pgx.NamedArgs{
		"dis_id":    userDetails.disId,
		"commision": commision.DistributorCommision,
	}).Scan(&userDetails.disAfterBalance); err != nil {
		return err
	}

	// 6️⃣ Update retailer wallet
	updateRetailerWalletBalanceQuery := `
		UPDATE retailers
		SET retailer_wallet_balance = retailer_wallet_balance - @after_balance
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance;
	`

	if err := tx.QueryRow(ctx, updateRetailerWalletBalanceQuery, pgx.NamedArgs{
		"retailer_id":   req.RetailerId,
		"after_balance": req.Amount + (commision.TotalCommision - commision.RetailerCommision),
	}).Scan(&userDetails.retailerAfterBalance); err != nil {
		return err
	}

	// 7️⃣ Wallet transactions

	insertToWalletTransactions := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			credit_amount,
			debit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@credit_amount,
			@debit_amount,
			@before_balance,
			@after_balance,
			@transaction_reason,
			@remarks
		);
	`

	// Admin
	if _, err := tx.Exec(ctx, insertToWalletTransactions, pgx.NamedArgs{
		"user_id":            userDetails.adminId,
		"reference_id":       userDetails.transactionId,
		"credit_amount":      commision.AdminCommision,
		"debit_amount":       0,
		"before_balance":     userDetails.adminBeforeBalance,
		"after_balance":      userDetails.adminAfterBalance,
		"transaction_reason": "PAYOUT",
		"remarks":            "Payout commission credited",
	}); err != nil {
		return err
	}

	// Master Distributor
	if _, err := tx.Exec(ctx, insertToWalletTransactions, pgx.NamedArgs{
		"user_id":            userDetails.mdId,
		"reference_id":       userDetails.transactionId,
		"credit_amount":      commision.MasterDistributorCommision,
		"debit_amount":       0,
		"before_balance":     userDetails.mdBeforeBalance,
		"after_balance":      userDetails.mdAfterBalance,
		"transaction_reason": "PAYOUT",
		"remarks":            "Payout commission credited",
	}); err != nil {
		return err
	}

	// Distributor
	if _, err := tx.Exec(ctx, insertToWalletTransactions, pgx.NamedArgs{
		"user_id":            userDetails.disId,
		"reference_id":       userDetails.transactionId,
		"credit_amount":      commision.DistributorCommision,
		"debit_amount":       0,
		"before_balance":     userDetails.disBeforeBalance,
		"after_balance":      userDetails.disAfterBalance,
		"transaction_reason": "PAYOUT",
		"remarks":            "Payout commission credited",
	}); err != nil {
		return err
	}

	// Retailer (DEBIT)
	if _, err := tx.Exec(ctx, insertToWalletTransactions, pgx.NamedArgs{
		"user_id":            req.RetailerId,
		"reference_id":       userDetails.transactionId,
		"debit_amount":       req.Amount + (commision.TotalCommision - commision.RetailerCommision),
		"credit_amount":      0,
		"before_balance":     userDetails.retailerBeforeBalance,
		"after_balance":      userDetails.retailerAfterBalance,
		"transaction_reason": "PAYOUT",
		"remarks":            "Payout amount debited",
	}); err != nil {
		return err
	}

	// 8️⃣ Commit
	return tx.Commit(ctx)
}

func (db *Database) CreatePayoutFailureQuery(
	ctx context.Context,
	req models.CreatePayoutRequestModel,
	commision models.GetPayoutCommisionModel,
) error {

	insertToPayoutTransactionQuery := `
		INSERT INTO payout_transactions (
			partner_request_id,
			operator_transaction_id,
			order_id,
			retailer_id,
			mobile_number,
			bank_name,
			beneficiary_name,
			account_number,
			ifsc_code,
			amount,
			transfer_type,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			payout_transaction_status
		) VALUES (
			@partner_request_id,
			@operator_transaction_id,
			@order_id,
			@retailer_id,
			@mobile_number,
			@bank_name,
			@beneficiary_name,
			@account_number,
			@ifsc_code,
			@amount,
			@transfer_type,
			@admin_commision,
			@md_commision,
			@dis_commision,
			@retailer_commision,
			@status
		)
		RETURNING payout_transaction_id;
	`

	var transferType string

	if req.TransferType == 5 {
		transferType = "IMPS"
	}

	if req.TransferType == 6 {
		transferType = "NEFT"
	}

	if _, err := db.pool.Exec(ctx, insertToPayoutTransactionQuery, pgx.NamedArgs{
		"partner_request_id":      req.PartnerRequestId,
		"operator_transaction_id": req.OperatorTransactionId,
		"order_id":                req.OrderId,
		"retailer_id":             req.RetailerId,
		"mobile_number":           req.MobileNumber,
		"bank_name":               req.BankName,
		"beneficiary_name":        req.BeneficiaryName,
		"account_number":          req.AccountNumber,
		"ifsc_code":               req.IFSCCode,
		"amount":                  req.Amount,
		"transfer_type":           transferType,
		"admin_commision":         commision.AdminCommision,
		"md_commision":            commision.MasterDistributorCommision,
		"dis_commision":           commision.DistributorCommision,
		"retailer_commision":      commision.RetailerCommision,
		"status":                  req.TransactionStatus,
	}); err != nil {
		return err
	}
	return nil
}

func (db *Database) GetAllPayoutTransactionsQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetAllPayoutTransactionsResponseModel, error) {
	query := `
		SELECT
			p.payout_transaction_id,
			p.operator_transaction_id,
			p.order_id,
			p.partner_request_id,
			p.retailer_id,
			r.retailer_name,
			r.retailer_business_name,
			p.mobile_number,
			p.bank_name,
			p.beneficiary_name,
			p.account_number,
			p.ifsc_code,
			p.amount,
			p.transfer_type,
			p.admin_commision,
			p.master_distributor_commision,
			p.distributor_commision,
			p.retailer_commision,
			p.payout_transaction_status,
			p.created_at,
			p.updated_at
		FROM payout_transactions p
		JOIN retailers r
			ON r.retailer_id = p.retailer_id
		ORDER BY created_at DESC
		LIMIT @limit AND OFFSET @offset;
	`

	res, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var transactions []models.GetAllPayoutTransactionsResponseModel
	for res.Next() {
		var transaction models.GetAllPayoutTransactionsResponseModel
		if err := res.Scan(
			&transaction.PayoutTransactionId,
			&transaction.OperatorTransactionId,
			&transaction.OrderId,
			&transaction.PartnerRequestId,
			&transaction.RetailerId,
			&transaction.RetailerName,
			&transaction.RetailerBusinessName,
			&transaction.MobileNumber,
			&transaction.BankName,
			&transaction.BeneficiaryName,
			&transaction.AccountNumber,
			&transaction.IFSCCode,
			&transaction.Amount,
			&transaction.TransferType,
			&transaction.AdminCommision,
			&transaction.MasterDistributorCommision,
			&transaction.DistributorCommision,
			&transaction.RetailerCommision,
			&transaction.TransactionStatus,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, res.Err()
}

func (db *Database) GetPayoutTransactionsByRetailerIdQuery(
	ctx context.Context,
	retailerId string,
	limit, offset int,
) ([]models.GetRetailerPayoutTransactionsResponseModel, error) {
	query := `
		SELECT
    p.payout_transaction_id,
    p.operator_transaction_id,
    p.order_id,
    p.partner_request_id,
    p.retailer_id,
    r.retailer_name,
    r.retailer_business_name,
    p.mobile_number,
    p.bank_name,
    p.beneficiary_name,
    p.account_number,
    p.ifsc_code,
    p.amount,
    p.transfer_type,
    p.retailer_commision,
    w.before_balance,
    w.after_balance,
    p.payout_transaction_status,
    p.created_at,
    p.updated_at
FROM payout_transactions p
JOIN retailers r
    ON r.retailer_id = p.retailer_id
JOIN wallet_transactions w
    ON w.user_id = p.retailer_id
   AND w.reference_id = p.payout_transaction_id::TEXT
   AND w.transaction_reason = 'PAYOUT'
WHERE p.retailer_id = @retailer_id
ORDER BY p.created_at DESC
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

	var transactions []models.GetRetailerPayoutTransactionsResponseModel
	for res.Next() {
		var transaction models.GetRetailerPayoutTransactionsResponseModel
		if err := res.Scan(
			&transaction.PayoutTransactionId,
			&transaction.OperatorTransactionId,
			&transaction.OrderId,
			&transaction.PartnerRequestId,
			&transaction.RetailerId,
			&transaction.RetailerName,
			&transaction.RetailerBusinessName,
			&transaction.MobileNumber,
			&transaction.BankName,
			&transaction.BeneficiaryName,
			&transaction.AccountNumber,
			&transaction.IFSCCode,
			&transaction.Amount,
			&transaction.TransferType,
			&transaction.RetailerCommision,
			&transaction.BeforeBalance,
			&transaction.AfterBalance,
			&transaction.TransactionStatus,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, res.Err()
}
