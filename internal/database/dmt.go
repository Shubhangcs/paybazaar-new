package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (db *Database) GetRetailerAadharNumberForDMTQuery(
	ctx context.Context,
	retailerId string,
) (string, error) {
	query := `
		SELECT retailer_aadhar_number
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`
	var aadharNumber string
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerId,
	}).Scan(
		&aadharNumber,
	); err != nil {
		return "", err
	}
	return aadharNumber, nil
}
