package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateBankQuery(
	ctx context.Context,
	req models.CreateBankModel,
) (int64, error) {

	query := `
		INSERT INTO banks (
			bank_name,
			ifsc_code
		) VALUES (
			@bank_name,
			@ifsc_code
		)
		RETURNING bank_id;
	`

	var bankID int64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"bank_name": req.BankNams,
		"ifsc_code": req.IFSCCode,
	}).Scan(&bankID)

	return bankID, err
}

func (db *Database) GetBankByIDQuery(
	ctx context.Context,
	bankID int64,
) (*models.GetBankModel, error) {

	query := `
		SELECT
			bank_id,
			bank_name,
			ifsc_code
		FROM banks
		WHERE bank_id = @bank_id;
	`

	var bank models.GetBankModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"bank_id": bankID,
	}).Scan(
		&bank.BankID,
		&bank.BankName,
		&bank.IFSCCode,
	)

	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (db *Database) GetAllBanksQuery(
	ctx context.Context,
) ([]models.GetBankModel, error) {

	query := `
		SELECT
			bank_id,
			bank_name,
			ifsc_code
		FROM banks
		ORDER BY bank_name ASC;
	`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banks []models.GetBankModel

	for rows.Next() {
		var b models.GetBankModel
		if err := rows.Scan(
			&b.BankID,
			&b.BankName,
			&b.IFSCCode,
		); err != nil {
			return nil, err
		}
		banks = append(banks, b)
	}

	return banks, rows.Err()
}

func (db *Database) UpdateBankQuery(
	ctx context.Context,
	bankID int64,
	req models.UpdateBankModel,
) error {

	query := `
		UPDATE banks
		SET
			bank_name = COALESCE(@bank_name, bank_name),
			ifsc_code = COALESCE(@ifsc_code, ifsc_code)
		WHERE bank_id = @bank_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"bank_id":   bankID,
		"bank_name": req.BankName,
		"ifsc_code": req.IFSCCode,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) DeleteBankQuery(
	ctx context.Context,
	bankID int64,
) error {

	query := `
		DELETE FROM banks
		WHERE bank_id = @bank_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"bank_id": bankID,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) CreateAdminBankQuery(
	ctx context.Context,
	req models.CreateAdminBankModel,
) (int64, error) {

	query := `
		INSERT INTO admin_banks (
			admin_id,
			admin_bank_name,
			admin_bank_account_number,
			admin_bank_ifsc_code
		) VALUES (
			@admin_id,
			@bank_name,
			@account_number,
			@ifsc_code
		)
		RETURNING admin_bank_id;
	`

	var adminBankID int64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"admin_id":       req.AdminID,
		"bank_name":      req.BankNams,
		"account_number": req.AccountNumber,
		"ifsc_code":      req.IFSCCode,
	}).Scan(&adminBankID)

	return adminBankID, err
}

func (db *Database) GetAdminBankByIDQuery(
	ctx context.Context,
	adminBankID int64,
) (*models.GetAdminBankModel, error) {

	query := `
		SELECT
			admin_bank_id,
			admin_bank_name,
			admin_bank_account_number,
			admin_bank_ifsc_code
		FROM admin_banks
		WHERE admin_bank_id = @admin_bank_id;
	`

	var ab models.GetAdminBankModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"admin_bank_id": adminBankID,
	}).Scan(
		&ab.AdminBankID,
		&ab.BankName,
		&ab.AccountNumber,
		&ab.IFSCCode,
	)

	if err != nil {
		return nil, err
	}

	return &ab, nil
}

func (db *Database) GetAdminBanksByAdminIDQuery(
	ctx context.Context,
	adminID string,
) ([]models.GetAdminBankModel, error) {

	query := `
		SELECT
			admin_bank_id,
			admin_bank_name,
			admin_bank_account_number,
			admin_bank_ifsc_code
		FROM admin_banks
		WHERE admin_id = @admin_id
		ORDER BY admin_bank_id DESC;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banks []models.GetAdminBankModel

	for rows.Next() {
		var ab models.GetAdminBankModel
		if err := rows.Scan(
			&ab.AdminBankID,
			&ab.BankName,
			&ab.AccountNumber,
			&ab.IFSCCode,
		); err != nil {
			return nil, err
		}
		banks = append(banks, ab)
	}

	return banks, rows.Err()
}

func (db *Database) UpdateAdminBankQuery(
	ctx context.Context,
	adminBankID int64,
	req models.UpdateAdminBankModel,
) error {

	query := `
		UPDATE admin_banks
		SET
			admin_bank_name = COALESCE(@bank_name, admin_bank_name),
			admin_bank_account_number = COALESCE(@account_number, admin_bank_account_number),
			admin_bank_ifsc_code = COALESCE(@ifsc_code, admin_bank_ifsc_code)
		WHERE admin_bank_id = @admin_bank_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_bank_id":  adminBankID,
		"bank_name":      req.BankName,
		"account_number": req.AccountNumber,
		"ifsc_code":      req.IFSCCode,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) DeleteAdminBankQuery(
	ctx context.Context,
	adminBankID int64,
) error {

	query := `
		DELETE FROM admin_banks
		WHERE admin_bank_id = @admin_bank_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_bank_id": adminBankID,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
