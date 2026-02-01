package database

import (
	"context"

	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) GetUserDetailsForDMTWalletCreation(
	ctx context.Context,
	req *models.CreateDMTWalletVerifyRequestModel,
) error {
	query := `
		SELECT retailer_address, city, state, pincode
		FROM retailers 
		WHERE retailer_id=@retailer_id;
	`
	if err := db.pool.QueryRow(ctx, query, req.RetailerID).Scan(
		&req.AddressLine1,
		&req.City,
		&req.State,
		&req.PinCode,
	); err != nil {
		return err
	}
	if req.AddressLine2 == "" {
		req.AddressLine1 = req.AddressLine2
	}
	return nil
}
