package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateBankQuery(
	ctx context.Context,
	req models.CreateBankRequestModel,
) error {
	query := `
		INSERT INTO banks (
			bank_name,
			ifsc_code
		) VALUES (
			@bank_name,
			@ifsc_code
		);
	`

	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"bank_name": req.BankNams,
		"ifsc_code": req.IFSCCode,
	}); err != nil {
		return fmt.Errorf("failed to create bank")
	}
	return nil
}

func (db *Database) GetBankDetailsByBankIDQuery(
	ctx context.Context,
	bankID int64,
) (*models.GetBankDetailsResponseModel, error) {

	query := `
		SELECT
			bank_id,
			bank_name,
			ifsc_code
		FROM banks
		WHERE bank_id = @bank_id;
	`

	var bank models.GetBankDetailsResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"bank_id": bankID,
	}).Scan(
		&bank.BankID,
		&bank.BankName,
		&bank.IFSCCode,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch bank details")
	}

	return &bank, nil
}

func (db *Database) GetAllBanksQuery(
	ctx context.Context,
) ([]models.GetBankDetailsResponseModel, error) {

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
		return nil, fmt.Errorf("failed to fetch bank details")
	}
	defer rows.Close()

	var banks []models.GetBankDetailsResponseModel

	for rows.Next() {
		var b models.GetBankDetailsResponseModel
		if err := rows.Scan(
			&b.BankID,
			&b.BankName,
			&b.IFSCCode,
		); err != nil {
			return nil, fmt.Errorf("failed to fetch bank details")
		}
		banks = append(banks, b)
	}

	return banks, rows.Err()
}

func (db *Database) UpdateBankQuery(
	ctx context.Context,
	req models.UpdateBankDetailsRequestModel,
) error {

	query := `
		UPDATE banks
		SET
			bank_name = COALESCE(@bank_name, bank_name),
			ifsc_code = COALESCE(@ifsc_code, ifsc_code)
		WHERE bank_id = @bank_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"bank_id":   req.BankID,
		"bank_name": req.BankName,
		"ifsc_code": req.IFSCCode,
	})

	if err != nil {
		return fmt.Errorf("failed to update bank details")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid bank or bank not found")
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
		return fmt.Errorf("failed to delete bank")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid bank id or bank not found")
	}

	return nil
}

func (db *Database) CreateAdminBankQuery(
	ctx context.Context,
	req models.CreateAdminBankRequestModel,
) error {

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
		);
	`

	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":       req.AdminID,
		"bank_name":      req.BankNams,
		"account_number": req.AccountNumber,
		"ifsc_code":      req.IFSCCode,
	}); err != nil {
		return fmt.Errorf("failed to create admin bank")
	}

	return nil
}

func (db *Database) GetAdminBankDetailsByAdminBankIDQuery(
	ctx context.Context,
	adminBankID int64,
) (*models.GetAdminBankDetailsResponseModel, error) {

	query := `
		SELECT
			admin_bank_id,
			admin_bank_name,
			admin_bank_account_number,
			admin_bank_ifsc_code
		FROM admin_banks
		WHERE admin_bank_id = @admin_bank_id;
	`

	var ab models.GetAdminBankDetailsResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"admin_bank_id": adminBankID,
	}).Scan(
		&ab.AdminBankID,
		&ab.BankName,
		&ab.AccountNumber,
		&ab.IFSCCode,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch admin bank details")
	}

	return &ab, nil
}

func (db *Database) GetAdminBanksByAdminIDQuery(
	ctx context.Context,
	adminID string,
) ([]models.GetAdminBankDetailsResponseModel, error) {

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
		return nil, fmt.Errorf("failed to fetch admin banks")
	}
	defer rows.Close()

	var banks []models.GetAdminBankDetailsResponseModel

	for rows.Next() {
		var ab models.GetAdminBankDetailsResponseModel
		if err := rows.Scan(
			&ab.AdminBankID,
			&ab.BankName,
			&ab.AccountNumber,
			&ab.IFSCCode,
		); err != nil {
			return nil, fmt.Errorf("failed to fetch admin banks")
		}
		banks = append(banks, ab)
	}

	return banks, rows.Err()
}

func (db *Database) UpdateAdminBankQuery(
	ctx context.Context,
	req models.UpdateAdminBankDetailsRequestModel,
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
		"admin_bank_id":  req.AdminBankID,
		"bank_name":      req.BankName,
		"account_number": req.AccountNumber,
		"ifsc_code":      req.IFSCCode,
	})

	if err != nil {
		return fmt.Errorf("failed to update admin bank details")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid admin bank id or admin bank not found")
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
		return fmt.Errorf("failed to delete admin bank")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid admin bank id or admin bank not found")
	}

	return nil
}
