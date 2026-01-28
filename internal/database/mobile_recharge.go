package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) GetAllMobileRechargeOperatorsQuery(
	ctx context.Context,
) ([]models.GetMobileRechargeOperatorsResponseModel, error) {
	query := `
		SELECT 
			operator_name,
			operator_code 
		FROM mobile_recharge_operators;
	`
	res, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var operators []models.GetMobileRechargeOperatorsResponseModel
	for res.Next() {
		var operator models.GetMobileRechargeOperatorsResponseModel
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

func (db *Database) GetAllMobileRechargeCirclesQuery(
	ctx context.Context,
) ([]models.GetMobileRechargeCircleResponseModel, error) {
	query := `
		SELECT 
			circle_name,
			circle_code
		FROM mobile_recharge_circles;
	`
	res, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var circles []models.GetMobileRechargeCircleResponseModel
	for res.Next() {
		var circle models.GetMobileRechargeCircleResponseModel
		if err := res.Scan(
			&circle.CircleName,
			&circle.CircleCode,
		); err != nil {
			return nil, err
		}
		circles = append(circles, circle)
	}
	return circles, res.Err()
}

func (db *Database) CreateMobileRechargeQuery(
	ctx context.Context,
	req models.CreateMobileRechargeRequestModel,
) error {
	if req.Amount <= 99 {
		return db.mobileRechargeWithoutCommision(ctx, req)
	}
	return db.mobileRechargeWithCommision(ctx, req)
}

func (db *Database) mobileRechargeWithoutCommision(
	ctx context.Context,
	req models.CreateMobileRechargeRequestModel,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	getRetailerBeforeBalance := `
		SELECT retailer_wallet_balance AS before_balance 
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`
	var beforeBalance float64
	if err := tx.QueryRow(ctx, getRetailerBeforeBalance, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(
		&beforeBalance,
	); err != nil {
		return err
	}

	deductAmountFromRetailerQuery := `
		UPDATE retailers 
		SET retailer_wallet_balance = retailer_wallet_balance - @amount
		WHERE retailer_id = @retailer_id
		RETURNING retailer_wallet_balance as after_balance;
	`
	var afterBalance float64
	if err := tx.QueryRow(ctx, deductAmountFromRetailerQuery, pgx.NamedArgs{
		"amount":      req.Amount,
		"retailer_id": req.RetailerID,
	}).Scan(
		&afterBalance,
	); err != nil {
		return err
	}

	insertToMobileRechargeTable := `
		INSERT INTO mobile_recharge (
    		retailer_id,
    		partner_request_id,
    		mobile_number,
    		operator_name,
    		circle_name,
    		operator_code,
    		circle_code,
    		amount,
    		commision,
    		recharge_type
		) VALUES (
    		@retailer_id,
    		@partner_request_id,
    		@mobile_number,
    		@operator_name,
    		@circle_name,
    		@operator_code,
    		@circle_code,
    		@amount,
    		@commision,
    		@recharge_type
		)
		RETURNING mobile_recharge_transaction_id AS transaction_id;
	`
	var transactionId string
	if err := tx.QueryRow(ctx, insertToMobileRechargeTable, pgx.NamedArgs{
		"retailer_id":        req.RetailerID,
		"partner_request_id": req.PartnerRequestID,
		"mobile_number":      req.MobileNumber,
		"operator_name":      req.OperatorName,
		"circle_name":        req.CircleName,
		"operator_code":      req.OperatorCode,
		"circle_code":        req.CircleCode,
		"amount":             req.Amount,
		"commision":          0,
		"recharge_type":      1,
	}).Scan(
		&transactionId,
	); err != nil {
		return err
	}

	insertToWalletTransactionsTable := `
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
	if _, err := tx.Exec(ctx, insertToWalletTransactionsTable, pgx.NamedArgs{
		"user_id":            req.RetailerID,
		"reference_id":       transactionId,
		"debit_amount":       req.Amount,
		"before_balance":     beforeBalance,
		"after_balance":      afterBalance,
		"transaction_reason": "MOBILE_RECHARGE",
		"remarks":            fmt.Sprintf("Mobile Recharge to: %s", req.MobileNumber),
	}); err != nil {
		return err
	}
	return nil
}

func (db *Database) mobileRechargeWithCommision(
	ctx context.Context,
	req models.CreateMobileRechargeRequestModel,
) error {
	getAdminBeforeBalanceQuery := `
		SELECT admin_wallet_balance AS admin_before_balance
		FROM admins AS ad
		JOIN master_distributors AS md
			ON md.admin_id = ad.admin_id
		JOIN distributors AS d
			ON d.master_distributor_id = md.master_distributor_id
		JOIN retailers AS r
			ON r.distributor_id = d.distributor_id
		WHERE r.retailer_id = @retailer_id;
	`

	deductAdminAmountAndGetAfterBalanceQuery := `
		UPDATE admins AS ad
		SET ad.admin_wallet_balance = ad.admin_wallet_balance - 1
		JOIN master_distributors AS md
			ON md.admin_id = ad.admin_id
		JOIN distributors AS d
			ON d.master_distributor_id = md.master_distributor_id
		JOIN retailers AS r
			ON r.distributor_id = d.distributor_id
		WHERE r.retailer_id = @retailer_id
		RETURNING ad.admin_wallet_balance as admin_after_balance;
	`
}
