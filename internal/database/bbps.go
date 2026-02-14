package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreatePostpaidMobileRechargeSuccessOrPendingQuery(
	ctx context.Context,
	req models.CreatePostpaidMobileRechargeAPIRequestModel,
	txn models.GetPostpaidMobileRechargeAPIResponseModel,
	status string,
) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getRetailerWalletBalanceQuery := `
		SELECT retailer_wallet_balance 
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`
	var retailerBeforeBalance float64
	if err := tx.QueryRow(ctx, getRetailerWalletBalanceQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(&retailerBeforeBalance); err != nil {
		return err
	}

	deductFromRetailerWalletQuery := `
		UPDATE retailers 
		SET retailer_wallet_balance = retailer_wallet_balance - @amount
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance;
	`
	var retailerAfterBalance float64
	if err := tx.QueryRow(ctx, deductFromRetailerWalletQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
		"amount":      req.Amount,
	}).Scan(
		&retailerAfterBalance,
	); err != nil {
		return err
	}

	insertToPostpaidMobileRechargeTable := `
		INSERT INTO mobile_recharge_postpaid (
			retailer_id,
			partner_request_id,
			order_id,
			operator_transaction_id,
			mobile_number,
			operator_code,
			amount,
			circle_code,
			circle_name,
			operator_name,
			recharge_type,
			recharge_status,
			commision
		) VALUES (
			@retailer_id,
			@partner_request_id,
			@order_id,
			@operator_transaction_id,
			@mobile_number,
			@operator_code,
			@amount,
			@circle_code,
			@circle_name,
			@operator_name,
			@recharge_type,
			@recharge_status,
			@commision 
		)
		RETURNING postpaid_recharge_transaction_id;
	`
	var transactionId int
	if err := tx.QueryRow(ctx, insertToPostpaidMobileRechargeTable, pgx.NamedArgs{
		"retailer_id":             req.RetailerID,
		"partner_request_id":      req.PartnerRequestID,
		"order_id":                txn.OrderID,
		"operator_transaction_id": txn.OperatorTransactionID,
		"mobile_number":           req.MobileNumber,
		"operator_code":           req.OperatorCode,
		"amount":                  req.Amount,
		"circle_code":             req.OperatorCircle,
		"circle_name":             req.CircleName,
		"operator_name":           req.OperatorName,
		"recharge_type":           fmt.Sprintf("%d", 1),
		"recharge_status":         status,
		"commision":               0,
	}).Scan(&transactionId); err != nil {
		return err
	}

	insertToWalletTransactions := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			before_balance,
			after_balance,
			debit_amount,
			credit_amount,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@before_balance,
			@after_balance,
			@debit_amount,
			@credit_amount,
			@transaction_reason,
			@remarks
		);
	`
	if _, err := tx.Exec(ctx, insertToWalletTransactions, pgx.NamedArgs{
		"user_id":            req.RetailerID,
		"reference_id":       fmt.Sprintf("%d", transactionId),
		"before_balance":     retailerBeforeBalance,
		"after_balance":      retailerAfterBalance,
		"debit_amount":       req.Amount,
		"credit_amount":      0,
		"transaction_reason": "POSTPAID_MOBILE_RECHARGE",
		"remarks":            fmt.Sprintf("Postpaid mobile recharge to: %s", req.MobileNumber),
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (db *Database) CreatePostpaidMobileRechargeFailureQuery(
	ctx context.Context,
	req models.CreatePostpaidMobileRechargeAPIRequestModel,
	txn models.GetPostpaidMobileRechargeAPIResponseModel,
) error {
	insertToPostpaidMobileRechargeTable := `
		INSERT INTO mobile_recharge_postpaid (
			retailer_id,
			partner_request_id,
			order_id,
			operator_transaction_id,
			mobile_number,
			operator_code,
			amount,
			circle_code,
			circle_name,
			operator_name,
			recharge_type,
			recharge_status,
			commision
		) VALUES (
			@retailer_id,
			@partner_request_id,
			@order_id,
			@operator_transaction_id,
			@mobile_number,
			@operator_code,
			@amount,
			@circle_code,
			@circle_name,
			@operator_name,
			@recharge_type,
			@recharge_status,
			@commision 
		);
	`
	if _, err := db.pool.Exec(ctx, insertToPostpaidMobileRechargeTable, pgx.NamedArgs{
		"retailer_id":             req.RetailerID,
		"partner_request_id":      req.PartnerRequestID,
		"order_id":                txn.OrderID,
		"operator_transaction_id": txn.OperatorTransactionID,
		"mobile_number":           req.MobileNumber,
		"operator_code":           req.OperatorCode,
		"amount":                  req.Amount,
		"circle_code":             req.OperatorCircle,
		"circle_name":             req.CircleName,
		"operator_name":           req.OperatorName,
		"recharge_type":           fmt.Sprintf("%d", 1),
		"recharge_status":         "FAILED",
		"commision":               0,
	}); err != nil {
		return err
	}
	return nil
}

func (db *Database) GetAllPostpaidMobileRechargeQuery(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.GetPostpaidMobileRechargeHistoryResponseModel, error) {

	query := `
		SELECT
			m.postpaid_recharge_transaction_id,
			m.retailer_id,
			r.retailer_name,
			r.retailer_business_name,
			m.partner_request_id,
			m.operator_transaction_id,
			m.order_id,
			m.mobile_number,
			m.operator_code,
			m.amount,
			w.before_balance,
			w.after_balance,
			m.circle_code,
			m.circle_name,
			m.operator_name,
			m.recharge_type,
			m.recharge_status,
			m.commision,
			m.created_at
		FROM mobile_recharge_postpaid m
		JOIN retailers r
			ON r.retailer_id = m.retailer_id
		JOIN wallet_transactions w
			ON w.user_id = m.retailer_id
			AND w.reference_id = m.postpaid_recharge_transaction_id::TEXT
			AND w.transaction_reason = 'POSTPAID_MOBILE_RECHARGE'
		ORDER BY m.created_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := db.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]models.GetPostpaidMobileRechargeHistoryResponseModel, 0)

	for rows.Next() {
		var item models.GetPostpaidMobileRechargeHistoryResponseModel

		err := rows.Scan(
			&item.PostpaidRechargeTransactionID,
			&item.RetailerID,
			&item.RetailerName,
			&item.RetailerBusinessName,
			&item.PartnerRequestID,
			&item.OperatorTransactionID,
			&item.OrderID,
			&item.MobileNumber,
			&item.OperatorCode,
			&item.Amount,
			&item.BeforeBalance,
			&item.AfterBalance,
			&item.CircleCode,
			&item.CircleName,
			&item.OperatorName,
			&item.RechargeType,
			&item.RechargeStatus,
			&item.Commission,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if item.RechargeStatus == "PENDING" {
			newStatus, err := db.DTHRechargeStatusCheck(item.PartnerRequestID)
			if err != nil {
				return nil, err
			}
			fmt.Println(newStatus)
			if newStatus != "PENDING" {
				if err := db.UpdatePostpaidMobileRechargeStatus(ctx, item.PostpaidRechargeTransactionID, newStatus); err != nil {
					return nil, err
				}
				item.RechargeStatus = newStatus
			}
		}

		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func (db *Database) UpdatePostpaidMobileRechargeStatus(
	ctx context.Context,
	transactionId int,
	status string,
) error {
	query := `
		UPDATE mobile_recharge_postpaid
		SET recharge_status = @status
		WHERE postpaid_recharge_transaction_id = @transaction_id;
	`
	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"status":         status,
		"transaction_id": transactionId,
	}); err != nil {
		return err
	}
	return nil
}

func (db *Database) GetPostpaidMobileRechargeByRetailerIDQuery(
	ctx context.Context,
	retailerID string,
	limit int,
	offset int,
) ([]models.GetPostpaidMobileRechargeHistoryResponseModel, error) {

	query := `
		SELECT
			m.postpaid_recharge_transaction_id,
			m.retailer_id,
			r.retailer_name,
			r.retailer_business_name,
			m.partner_request_id,
			m.operator_transaction_id,
			m.order_id,
			m.mobile_number,
			m.operator_code,
			m.amount,
			w.before_balance,
			w.after_balance,
			m.circle_code,
			m.circle_name,
			m.operator_name,
			m.recharge_type,
			m.recharge_status,
			m.commision,
			m.created_at
		FROM mobile_recharge_postpaid m
		JOIN retailers r
			ON r.retailer_id = m.retailer_id
		JOIN wallet_transactions w
			ON w.user_id = m.retailer_id
			AND w.reference_id = m.postpaid_recharge_transaction_id::TEXT
			AND w.transaction_reason = 'POSTPAID_MOBILE_RECHARGE'
		WHERE m.retailer_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := db.pool.Query(ctx, query, retailerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]models.GetPostpaidMobileRechargeHistoryResponseModel, 0)

	for rows.Next() {
		var item models.GetPostpaidMobileRechargeHistoryResponseModel

		err := rows.Scan(
			&item.PostpaidRechargeTransactionID,
			&item.RetailerID,
			&item.RetailerName,
			&item.RetailerBusinessName,
			&item.PartnerRequestID,
			&item.OperatorTransactionID,
			&item.OrderID,
			&item.MobileNumber,
			&item.OperatorCode,
			&item.Amount,
			&item.BeforeBalance,
			&item.AfterBalance,
			&item.CircleCode,
			&item.CircleName,
			&item.OperatorName,
			&item.RechargeType,
			&item.RechargeStatus,
			&item.Commission,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if item.RechargeStatus == "PENDING" {
			newStatus, err := db.DTHRechargeStatusCheck(item.PartnerRequestID)
			if err != nil {
				return nil, err
			}
			if newStatus != "PENDING" {
				if err := db.UpdatePostpaidMobileRechargeStatus(ctx, item.PostpaidRechargeTransactionID, newStatus); err != nil {
					return nil, err
				}
				item.RechargeStatus = newStatus
			}
		}

		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func (db *Database) RefundPostpaidMobileRechargeQuery(
	ctx context.Context,
	transactionId int,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getTransactionDetails := `
		UPDATE mobile_recharge_postpaid
		SET recharge_status = @status
		WHERE postpaid_recharge_transaction_id = @transaction_id
		RETURNING retailer_id, amount;
	`
	var (
		retailerID string
		amount     float64
	)
	if err := tx.QueryRow(ctx, getTransactionDetails, pgx.NamedArgs{
		"status":         "REFUND",
		"transaction_id": transactionId,
	}).Scan(
		&retailerID,
		&amount,
	); err != nil {
		return err
	}

	getRetailerBeforeBalanceQuery := `
		SELECT retailer_wallet_balance 
		FROM retailers
		WHERE retailer_id = @retailer_id
	`
	var retailerBeforeBalance float64
	if err := tx.QueryRow(ctx, getRetailerBeforeBalanceQuery, pgx.NamedArgs{
		"retailer_id": retailerID,
	}).Scan(
		retailerBeforeBalance,
	); err != nil {
		return err
	}

	updateRetailersWalletAndGetAfterBalanceQuery := `
		UPDATE retailers
		SET retailer_wallet_balance = retailer_wallet_balance + @amount
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance;
	`
	var retailerAfterBalance float64
	if err := tx.QueryRow(ctx, updateRetailersWalletAndGetAfterBalanceQuery, pgx.NamedArgs{
		"amount":      amount,
		"retailer_id": retailerID,
	}).Scan(
		&retailerAfterBalance,
	); err != nil {
		return err
	}

	insertToWalletTransactionsQuery := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			debit_amount,
			credit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@debit_amount,
			@credit_amount,
			@before_balance,
			@after_balance,
			@transaction_reason,
			@remarks 
		);
	`
	if _, err := tx.Exec(ctx, insertToWalletTransactionsQuery, pgx.NamedArgs{
		"user_id":            retailerID,
		"reference_id":       fmt.Sprintf("%d", transactionId),
		"debit_amount":       0,
		"credit_amount":      amount,
		"before_balance":     retailerBeforeBalance,
		"after_balance":      retailerAfterBalance,
		"transaction_reason": "POSTPAID_MOBILE_RECHARGE_REFUND",
		"remarks":            fmt.Sprintf("Refunded %d transaction to %s", transactionId, retailerID),
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *Database) CreateElectricityBillPaymentSuccessOrPendingQuery(
	ctx context.Context,
	req models.CreateElectricityBillPaymentRequestModel,
	txn models.GetElectricityBillPaymentAPIResponseModel,
	status string,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getRetailerWalletBeforeBalanceQuery := `
		SELECT retailer_wallet_balance
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`
	var retailerBeforeBalance float64
	if err := tx.QueryRow(ctx, getRetailerWalletBeforeBalanceQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(&retailerBeforeBalance); err != nil {
		return err
	}

	updateRetailerWalletAndGetAfterBalanceQuery := `
		UPDATE retailers
		SET retailer_wallet_balance = retailer_wallet_balance - @amount
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance;
	`
	var retailerAfterBalance float64
	if err := tx.QueryRow(ctx, updateRetailerWalletAndGetAfterBalanceQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
		"amount":      req.Amount,
	}).Scan(&retailerAfterBalance); err != nil {
		return err
	}

	insertToElectricityBillTransactionsQuery := `
		INSERT INTO electricity_bill_payments (
			retailer_id,
			order_id,
			operator_transaction_id,
			partner_request_id,
			customer_id,
			amount,
			operator_code,
			operator_name,
			customer_email,
			commision,
			transaction_status
		) VALUES (
			@retailer_id,
			@order_id,
			@operator_transaction_id,
			@partner_request_id,
			@customer_id,
			@amount,
			@operator_code,
			@operator_name,
			@customer_email,
			@commision,
			@transaction_status
		)
		RETURNING electricity_bill_transaction_id;
	`
	var transactionId int
	if err := tx.QueryRow(ctx, insertToElectricityBillTransactionsQuery, pgx.NamedArgs{
		"retailer_id":             req.RetailerID,
		"order_id":                txn.OrderID,
		"operator_transaction_id": txn.OperatorTransactionID,
		"partner_request_id":      txn.PartnerRequestID,
		"customer_id":             req.CustomerID,
		"customer_email":          req.CustomerEmail,
		"amount":                  req.Amount,
		"operator_code":           req.OperatorCode,
		"operator_name":           req.OperatorName,
		"commision":               0,
		"transaction_status":      status,
	}).Scan(&transactionId); err != nil {
		return err
	}

	insertIntoWalletTransactionsQuery := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			before_balance,
			after_balance,
			credit_amount,
			debit_amount,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@before_balance,
			@after_balance,
			@credit_amount,
			@debit_amount,
			@transaction_reason,
			@remarks
		);
	`

	if _, err := tx.Exec(ctx, insertIntoWalletTransactionsQuery, pgx.NamedArgs{
		"user_id":            req.RetailerID,
		"reference_id":       fmt.Sprintf("%d", transactionId),
		"before_balance":     retailerBeforeBalance,
		"after_balance":      retailerAfterBalance,
		"credit_amount":      0,
		"debit_amount":       req.Amount,
		"transaction_reason": "ELECTRICITY_BILL",
		"remarks":            fmt.Sprintf("electricity bill paid to: %s", req.CustomerID),
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (db *Database) CreateElectricityBillPaymentFailureQuery(
	ctx context.Context,
	req models.CreateElectricityBillPaymentRequestModel,
	txn models.GetElectricityBillPaymentAPIResponseModel,
) error {
	insertToElectricityBillTransactionsQuery := `
		INSERT INTO electricity_bill_payments (
			retailer_id,
			order_id,
			operator_transaction_id,
			partner_request_id,
			customer_id,
			amount,
			operator_code,
			operator_name,
			customer_email,
			commision,
			transaction_status
		) VALUES (
			@retailer_id,
			@order_id,
			@operator_transaction_id,
			@partner_request_id,
			@customer_id,
			@amount,
			@operator_code,
			@operator_name,
			@customer_email,
			@commision,
			@transaction_status
		)
		RETURNING electricity_bill_transaction_id;
	`
	if _, err := db.pool.Exec(ctx, insertToElectricityBillTransactionsQuery, pgx.NamedArgs{
		"retailer_id":             req.RetailerID,
		"order_id":                txn.OrderID,
		"operator_transaction_id": txn.OperatorTransactionID,
		"partner_request_id":      txn.PartnerRequestID,
		"customer_id":             req.CustomerID,
		"customer_email":          req.CustomerEmail,
		"amount":                  req.Amount,
		"operator_code":           req.OperatorCode,
		"operator_name":           req.OperatorName,
		"commision":               0,
		"transaction_status":      "FAILED",
	}); err != nil {
		return err
	}
	return nil
}

func (db *Database) GetElectricityOperatorsQuery(
	ctx context.Context,
) ([]models.GetElectricityOperatorResponseModel, error) {
	query := `
		SELECT
			operator_name,
			operator_code
		FROM electricity_operators;
	`

	res, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var operators []models.GetElectricityOperatorResponseModel
	for res.Next() {
		var operator models.GetElectricityOperatorResponseModel
		if err := res.Scan(
			&operator.OperatorName,
			&operator.OperatorCode,
		); err != nil {
			return nil, err
		}
		operators = append(operators, operator)
	}
	return operators, res.Err()
}

func (db *Database) GetAllElectricityBillPaymentTransactionsQuery(
	ctx context.Context,
	offset, limit int,
) ([]models.GetElectricityBillHistoryResponseModel, error) {

	query := `
		SELECT
			e.electricity_bill_transaction_id,
			e.retailer_id,
			r.retailer_name,
			r.retailer_business_name,
			e.order_id,
			e.operator_transaction_id,
			e.partner_request_id,
			e.customer_id,
			e.amount,
			w.before_balance,
			w.after_balance,
			e.operator_code,
			e.operator_name,
			e.customer_email,
			e.commision,
			e.transaction_status,
			e.created_at
		FROM electricity_bill_payments e
		JOIN retailers r
			ON r.retailer_id = e.retailer_id
		JOIN wallet_transactions w
			ON w.user_id = e.retailer_id
			AND w.reference_id = e.electricity_bill_transaction_id::TEXT
			AND w.transaction_reason = 'ELECTRICITY_BILL'
		ORDER BY e.created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.GetElectricityBillHistoryResponseModel

	for rows.Next() {
		var tx models.GetElectricityBillHistoryResponseModel

		err := rows.Scan(
			&tx.ElectricityBillTransactionID,
			&tx.RetailerID,
			&tx.RetailerName,
			&tx.RetailerBusinessName,
			&tx.OrderID,
			&tx.OperatorTransactionID,
			&tx.PartnerRequestID,
			&tx.CustomerID,
			&tx.Amount,
			&tx.BeforeBalance,
			&tx.AfterBalance,
			&tx.OperatorID, // operator_code
			&tx.OperatorName,
			&tx.CustomerEmail,
			&tx.Commision,
			&tx.TransactionStatus,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if tx.TransactionStatus == "PENDING" {
			newStatus, err := db.DTHRechargeStatusCheck(tx.PartnerRequestID)
			if err != nil {
				return nil, err
			}
			if newStatus != "PENDING" {
				if err := db.UpdateElectricityBillStatusByTransactionID(ctx, tx.ElectricityBillTransactionID, newStatus); err != nil {
					return nil, err
				}
				tx.TransactionStatus = newStatus
			}
		}

		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (db *Database) GetElectricityBillPaymentTransactionsByRetailerIDQuery(
	ctx context.Context,
	retailerID string,
	offset, limit int,
) ([]models.GetElectricityBillHistoryResponseModel, error) {

	query := `
		SELECT
			e.electricity_bill_transaction_id,
			e.retailer_id,
			r.retailer_name,
			r.retailer_business_name,
			e.order_id,
			e.operator_transaction_id,
			e.partner_request_id,
			e.customer_id,
			e.amount,
			w.before_balance,
			w.after_balance,
			e.operator_code,
			e.operator_name,
			e.customer_email,
			e.commision,
			e.transaction_status,
			e.created_at
		FROM electricity_bill_payments e
		JOIN retailers r
			ON r.retailer_id = e.retailer_id
		JOIN wallet_transactions w
			ON w.user_id = e.retailer_id
			AND w.reference_id = e.electricity_bill_transaction_id::TEXT
			AND w.transaction_reason = 'ELECTRICITY_BILL'
		WHERE e.retailer_id = $1
		ORDER BY e.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := db.pool.Query(ctx, query, retailerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.GetElectricityBillHistoryResponseModel

	for rows.Next() {
		var tx models.GetElectricityBillHistoryResponseModel

		err := rows.Scan(
			&tx.ElectricityBillTransactionID,
			&tx.RetailerID,
			&tx.RetailerName,
			&tx.RetailerBusinessName,
			&tx.OrderID,
			&tx.OperatorTransactionID,
			&tx.PartnerRequestID,
			&tx.CustomerID,
			&tx.Amount,
			&tx.BeforeBalance,
			&tx.AfterBalance,
			&tx.OperatorID,
			&tx.OperatorName,
			&tx.CustomerEmail,
			&tx.Commision,
			&tx.TransactionStatus,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if tx.TransactionStatus == "PENDING" {
			newStatus, err := db.DTHRechargeStatusCheck(tx.PartnerRequestID)
			if err != nil {
				return nil, err
			}
			if newStatus != "PENDING" {
				if err := db.UpdateElectricityBillStatusByTransactionID(ctx, tx.ElectricityBillTransactionID, newStatus); err != nil {
					return nil, err
				}
				tx.TransactionStatus = newStatus
			}
		}

		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (db *Database) UpdateElectricityBillStatusByTransactionID(
	ctx context.Context,
	electricityBillTransactionID int,
	status string,
) error {

	query := `
		UPDATE electricity_bill_payments
		SET
			transaction_status = $1
		WHERE electricity_bill_transaction_id = $2;
	`

	_, err := db.pool.Exec(
		ctx,
		query,
		status,
		electricityBillTransactionID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RefundElectricityBillPaymentQuery(
	ctx context.Context,
	transactionId int,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getTransactionDetailsQuery := `
		UPDATE electricity_bill_payments
		SET transaction_status = @status
		WHERE electricity_bill_transaction_id = @transaction_id
		RETURNING retailer_id, amount;
	`
	var (
		retailerId string
		amount     float64
	)
	if err := tx.QueryRow(ctx, getTransactionDetailsQuery, pgx.NamedArgs{
		"transaction_id": transactionId,
		"status":         "REFUND",
	}).Scan(
		&retailerId,
		&amount,
	); err != nil {
		return err
	}

	getRetailerWalletBalanceQuery := `
		SELECT retailer_wallet_balance 
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`
	var retailerBeforeBalance float64
	if err := tx.QueryRow(ctx, getRetailerWalletBalanceQuery, pgx.NamedArgs{
		"retailer_id": retailerId,
	}).Scan(&retailerBeforeBalance); err != nil {
		return err
	}

	updateRetailerWalletAndGetAfterBalanceQuery := `
		UPDATE retailers
		SET retailer_wallet_balance = retailer_wallet_balance + @amount
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance;
	`
	var retailerAfterBalance float64
	if err := tx.QueryRow(ctx, updateRetailerWalletAndGetAfterBalanceQuery, pgx.NamedArgs{
		"retailer_id": retailerId,
		"amount":      amount,
	}).Scan(
		&retailerAfterBalance,
	); err != nil {
		return err
	}

	insertToWalletTransactionsQuery := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			debit_amount,
			credit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@debit_amount,
			@credit_amount,
			@before_balance,
			@after_balance,
			@transaction_reason,
			@remarks 
		);
	`
	if _, err := tx.Exec(ctx, insertToWalletTransactionsQuery, pgx.NamedArgs{
		"user_id":            retailerId,
		"reference_id":       fmt.Sprintf("%d", transactionId),
		"debit_amount":       0,
		"credit_amount":      amount,
		"before_balance":     retailerBeforeBalance,
		"after_balance":      retailerAfterBalance,
		"transaction_reason": "ELECTRICITY_BILL_REFUND",
		"remarks":            fmt.Sprintf("transaction %d refunded to %s", transactionId, retailerId),
	}); err != nil {
		return err
	}
	return nil
}
