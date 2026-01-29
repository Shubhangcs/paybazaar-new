package database

import (
	"context"

	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) GetAllDTHOperatorsQuery(
	ctx context.Context,
) ([]models.GetDTHOperatorsResponseModel, error) {
	query := `
		SELECT operator_code, operator_name
		FROM dth_recharge;
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

// func (db *Database) CreateDTHRechargeRequest(
// 	ctx context.Context,
// 	req models.CreateDTHRechargeRequestModel,
// ) error {
// 	if req.Amount <= 99 {
// 		return db.dthRechargeWithoutCommision(ctx, req)
// 	}
// 	return db.dthRechargeWithCommision(ctx, req)
// }

// func (db *Database) dthRechargeWithoutCommision(
// 	ctx context.Context,
// 	req models.CreateDTHRechargeRequestModel,
// ) error {
// 	getRetailerBeforeBalanceQuery := `
// 		SELECT retailer_wallet_balance AS retailer_before_balance
// 		FROM retailers
// 		WHERE retailer_id = @retailer_id; 
// 	`

// 	deductRetailerAmountAndGetAfterBalance := `
// 		UPDATE retailers
// 		SET retailer_wallet_balance = retailer_wallet_balance - @amount
// 		WHERE retailer_id = @retailer_id
// 		RETURNING retailer_wallet_balance AS retailer_after_balance;
// 	`

// 	insertToDthRechargeTable := `
// 		INSERT INTO dth_recharge (
// 			retailer_id,
// 			customer_id,
// 			operator_name,
// 			operator_code,
// 			amount,
// 			partner_request_id,
// 			status,
// 			commision
// 		) VALUES (
// 			@retailer_id,
// 			@customer_id,
// 			@operator_name,
// 			@operator_code,
// 			@amount,
// 			@partner_request_id,
// 			@status,
// 			@commision
// 		)
// 		RETURNING dth_transaction_id AS transaction_id;
// 	`

// 	insertToRetailerWalletTransactionsTable := `
// 		INSERT INTO wallet_transaction (
// 			user_id,
// 			reference_id,
// 			debit_amount,
// 			before_balance,
// 			after_balance,
// 			transaction_reason,
// 			remarks
// 		) VALUES (
// 			@user_id,
// 			@reference_id,
// 			@debit_amount,
// 			@before_balance,
// 			@after_balance,
// 			@transaction_reason,
// 			@remarks 
// 		);
// 	`
	
// }

// func (db *Database) dthRechargeWithCommision(
// 	ctx context.Context,
// 	req models.CreateDTHRechargeRequestModel,
// ) error {

// }
