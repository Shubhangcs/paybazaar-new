package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateWalletTransactionQuery(
	ctx context.Context,
	req models.CreateWalletTransactionRequestModel,
) error {

	query := `
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
		)
	`

	_, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"user_id":            req.UserID,
		"reference_id":       req.ReferenceID,
		"credit_amount":      req.CreditAmount,
		"debit_amount":       req.DebitAmount,
		"before_balance":     req.BeforeBalance,
		"after_balance":      req.AfterBalance,
		"transaction_reason": req.TransactionReason,
		"remarks":            req.Remarks,
	})

	return err
}

func (db *Database) getWalletTransactionsByUserID(
	ctx context.Context,
	userID string,
	limit, offset int,
) ([]models.GetWalletTransactionResponseModel, error) {

	query := `
		SELECT
			wallet_transaction_id,
			user_id,
			reference_id,
			credit_amount,
			debit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks,
			created_at
		FROM wallet_transactions
		WHERE user_id = @user_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetWalletTransactionResponseModel

	for rows.Next() {
		var wt models.GetWalletTransactionResponseModel
		if err := rows.Scan(
			&wt.WalletTransactionID,
			&wt.UserID,
			&wt.ReferenceID,
			&wt.CreditAmount,
			&wt.DebitAmount,
			&wt.BeforeBalance,
			&wt.AfterBalance,
			&wt.TransactionReason,
			&wt.Remarks,
			&wt.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, wt)
	}

	return list, nil
}

func (db *Database) GetAdminWalletTransactionsQuery(
	ctx context.Context,
	adminID string,
	limit, offset int,
) ([]models.GetWalletTransactionResponseModel, error) {
	return db.getWalletTransactionsByUserID(ctx, adminID, limit, offset)
}

func (db *Database) GetMasterDistributorWalletTransactionsQuery(
	ctx context.Context,
	masterDistributorID string,
	limit, offset int,
) ([]models.GetWalletTransactionResponseModel, error) {
	return db.getWalletTransactionsByUserID(ctx, masterDistributorID, limit, offset)
}

func (db *Database) GetDistributorWalletTransactionsQuery(
	ctx context.Context,
	distributorID string,
	limit, offset int,
) ([]models.GetWalletTransactionResponseModel, error) {
	return db.getWalletTransactionsByUserID(ctx, distributorID, limit, offset)
}

func (db *Database) GetRetailerWalletTransactionsQuery(
	ctx context.Context,
	retailerID string,
	limit, offset int,
) ([]models.GetWalletTransactionResponseModel, error) {
	return db.getWalletTransactionsByUserID(ctx, retailerID, limit, offset)
}

func (db *Database) GetAdminWalletBalanceQuery(
	ctx context.Context,
	adminID string,
) (float64, error) {
	query := `
		SELECT admin_wallet_balance
		FROM admins
		WHERE admin_id = @admin_id
	`

	var balance float64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
	}).Scan(&balance)

	return balance, err
}

func (db *Database) GetMasterDistributorWalletBalanceQuery(
	ctx context.Context,
	masterDistributorID string,
) (float64, error) {
	query := `
		SELECT master_distributor_wallet_balance
		FROM master_distributors
		WHERE master_distributor_id = @md_id
	`

	var balance float64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"md_id": masterDistributorID,
	}).Scan(&balance)

	return balance, err
}

func (db *Database) GetDistributorWalletBalanceQuery(
	ctx context.Context,
	distributorID string,
) (float64, error) {
	query := `
		SELECT distributor_wallet_balance
		FROM distributors
		WHERE distributor_id = @distributor_id
	`

	var balance float64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
	}).Scan(&balance)

	return balance, err
}

func (db *Database) GetRetailerWalletBalanceQuery(
	ctx context.Context,
	retailerID string,
) (float64, error) {
	query := `
		SELECT retailer_wallet_balance
		FROM retailers
		WHERE retailer_id = @retailer_id
	`

	var balance float64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
	}).Scan(&balance)

	return balance, err
}
