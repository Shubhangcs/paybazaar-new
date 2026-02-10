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
		"recharge_type":           1,
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
		"recharge_type":           1,
		"recharge_status":         "FAILED",
		"commision":               0,
	}); err != nil {
		return err
	}
	return nil
}


