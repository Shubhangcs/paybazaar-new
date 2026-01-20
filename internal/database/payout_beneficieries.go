package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreatePayoutBeneficiaryQuery(
	ctx context.Context,
	req models.CreatePayoutBeneficiaryModel,
) (int64, error) {

	query := `
		INSERT INTO payout_beneficiaries (
			retailer_id,
			mobile_number,
			beneficiary_bank_name,
			beneficiary_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			beneficiary_phone
		) VALUES (
			@retailer_id,
			@mobile_number,
			@bank_name,
			@name,
			@account_number,
			@ifsc,
			@phone
		)
		RETURNING beneficiary_id;
	`

	var beneficiaryID int64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id":    req.RetailerID,
		"mobile_number":  req.MobileNumber,
		"bank_name":      req.BankName,
		"name":           req.BeneficiaryName,
		"account_number": req.AccountNumber,
		"ifsc":           req.IFSCCode,
		"phone":          req.Phone,
	}).Scan(&beneficiaryID)

	return beneficiaryID, err
}

func (db *Database) GetPayoutBeneficiaryByIDQuery(
	ctx context.Context,
	beneficiaryID int64,
) (*models.GetPayoutBeneficiaryResponseModel, error) {

	query := `
		SELECT
			beneficiary_id,
			retailer_id,
			mobile_number,
			beneficiary_bank_name,
			beneficiary_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			beneficiary_phone,
			is_beneficiary_verified,
			created_at,
			updated_at
		FROM payout_beneficiaries
		WHERE beneficiary_id = @id;
	`

	var b models.GetPayoutBeneficiaryResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id": beneficiaryID,
	}).Scan(
		&b.BeneficiaryID,
		&b.RetailerID,
		&b.MobileNumber,
		&b.BankName,
		&b.BeneficiaryName,
		&b.AccountNumber,
		&b.IFSCCode,
		&b.Phone,
		&b.IsVerified,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (db *Database) GetPayoutBeneficiariesByRetailerIDQuery(
	ctx context.Context,
	retailerID string,
	limit, offset int,
) ([]models.GetPayoutBeneficiaryResponseModel, error) {

	query := `
		SELECT
			beneficiary_id,
			retailer_id,
			mobile_number,
			beneficiary_bank_name,
			beneficiary_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			beneficiary_phone,
			is_beneficiary_verified,
			created_at,
			updated_at
		FROM payout_beneficiaries
		WHERE retailer_id = @retailer_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
		"limit":       limit,
		"offset":      offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetPayoutBeneficiaryResponseModel

	for rows.Next() {
		var b models.GetPayoutBeneficiaryResponseModel
		if err := rows.Scan(
			&b.BeneficiaryID,
			&b.RetailerID,
			&b.MobileNumber,
			&b.BankName,
			&b.BeneficiaryName,
			&b.AccountNumber,
			&b.IFSCCode,
			&b.Phone,
			&b.IsVerified,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	return list, rows.Err()
}

func (db *Database) GetPayoutBeneficiariesByMobileNumberQuery(
	ctx context.Context,
	mobileNumber string,
	limit, offset int,
) ([]models.GetPayoutBeneficiaryResponseModel, error) {

	query := `
		SELECT
			beneficiary_id,
			retailer_id,
			mobile_number,
			beneficiary_bank_name,
			beneficiary_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			beneficiary_phone,
			is_beneficiary_verified,
			created_at,
			updated_at
		FROM payout_beneficiaries
		WHERE mobile_number = @mobile_number
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"mobile_number": mobileNumber,
		"limit":         limit,
		"offset":        offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetPayoutBeneficiaryResponseModel

	for rows.Next() {
		var b models.GetPayoutBeneficiaryResponseModel
		if err := rows.Scan(
			&b.BeneficiaryID,
			&b.RetailerID,
			&b.MobileNumber,
			&b.BankName,
			&b.BeneficiaryName,
			&b.AccountNumber,
			&b.IFSCCode,
			&b.Phone,
			&b.IsVerified,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	return list, rows.Err()
}

func (db *Database) UpdatePayoutBeneficiaryQuery(
	ctx context.Context,
	beneficiaryID int64,
	req models.UpdatePayoutBeneficiaryModel,
) error {

	query := `
		UPDATE payout_beneficiaries
		SET
			beneficiary_bank_name = COALESCE(@bank_name, beneficiary_bank_name),
			beneficiary_name = COALESCE(@name, beneficiary_name),
			beneficiary_account_number = COALESCE(@account_number, beneficiary_account_number),
			beneficiary_ifsc_code = COALESCE(@ifsc, beneficiary_ifsc_code),
			beneficiary_phone = COALESCE(@phone, beneficiary_phone),
			updated_at = NOW()
		WHERE beneficiary_id = @id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"id":             beneficiaryID,
		"bank_name":      req.BankName,
		"name":           req.BeneficiaryName,
		"account_number": req.AccountNumber,
		"ifsc":           req.IFSCCode,
		"phone":          req.Phone,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) UpdatePayoutBeneficiaryVerificationQuery(
	ctx context.Context,
	beneficiaryID int64,
	isVerified bool,
) error {

	query := `
		UPDATE payout_beneficiaries
		SET
			is_beneficiary_verified = @verified,
			updated_at = NOW()
		WHERE beneficiary_id = @id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"id":       beneficiaryID,
		"verified": isVerified,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) DeletePayoutBeneficiaryQuery(
	ctx context.Context,
	beneficiaryID int64,
) error {

	tag, err := db.pool.Exec(ctx, `
		DELETE FROM payout_beneficiaries
		WHERE beneficiary_id = @id;
	`, pgx.NamedArgs{
		"id": beneficiaryID,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
