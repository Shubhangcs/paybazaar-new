package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) AddNewBeneficiary(req *models.BeneficiaryModel) error {
	query := `INSERT INTO beneficiaries (mobile_number, bank_name, ifsc_code, account_number, beneficiary_name, beneficiary_phone) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.pool.Exec(context.Background(), query, req.MobileNumber, req.BankName, req.IFSCCode, req.AccountNumber, req.BeneficiaryName, req.BeneficiaryPhone)
	return err
}

func (db *Database) GetBeneficiaries(mobileNumber string) (*[]models.BeneficiaryModel, error) {
	query := `SELECT beneficiary_id ,mobile_number ,bank_name, ifsc_code, account_number, beneficiary_name, beneficiary_phone, beneficiary_verified FROM beneficiaries WHERE mobile_number = $1`
	rows, err := db.pool.Query(context.Background(), query, mobileNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beneficiaries []models.BeneficiaryModel
	for rows.Next() {
		var beneficiary models.BeneficiaryModel
		err := rows.Scan(&beneficiary.BeneficiaryID, &beneficiary.MobileNumber, &beneficiary.BankName, &beneficiary.IFSCCode, &beneficiary.AccountNumber, &beneficiary.BeneficiaryName, &beneficiary.BeneficiaryPhone, &beneficiary.BeneficiaryVerified)
		if err != nil {
			return nil, err
		}
		beneficiaries = append(beneficiaries, beneficiary)
	}
	return &beneficiaries, nil
}

func (db *Database) VerifyBenificary(ctx context.Context, amount float64, retailerId, beneficiaryId string) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	deductAmountQuery := `
		UPDATE retailers 
		SET retailer_wallet_balance = retailer_wallet_balance - @amount
		WHERE retailer_id = @retailer_id;
	`
	if _, err := tx.Exec(ctx, deductAmountQuery, pgx.NamedArgs{
		"retailer_id": retailerId,
		"amount":      amount,
	}); err != nil {
		return err
	}
	query := `UPDATE beneficiaries SET beneficiary_verified = TRUE WHERE beneficiary_id = $1`
	if _, err := tx.Exec(ctx, query, beneficiaryId); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (db *Database) DeleteBeneficiary(beneficiaryId string) error {
	query := `DELETE FROM beneficiaries WHERE beneficiary_id=$1`
	_, err := db.pool.Exec(context.Background(), query, beneficiaryId)
	return err
}
