package database

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) VerifyMPINAndKycQuery(
	ctx context.Context,
	retailerID string,
	mpin int64,
) (bool, error) {

	var isValid bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM retailers
			WHERE
				retailer_id = @retailer_id
				AND retailer_mpin = @mpin
				AND retailer_kyc_status = TRUE
		);
	`

	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
		"mpin":        mpin,
	}).Scan(&isValid); err != nil {
		return false, err
	}

	return isValid, nil
}

func (db *Database) CheckRetailerWalletBalance(
	ctx context.Context,
	retailerID string,
	amount float64,
	commision float64,
) (bool, error) {

	var walletBalance float64

	query := `
		SELECT retailer_wallet_balance
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`

	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
	}).Scan(&walletBalance); err != nil {
		return false, err
	}

	requiredAmount := amount + commision

	if walletBalance >= requiredAmount {
		return true, nil
	}

	return false, nil
}

func defaultPayoutCommision() *models.GetCommisionModel {
	return &models.GetCommisionModel{
		TotalCommision:             1.20,
		AdminCommision:             0.35,
		MasterDistributorCommision: 0.05,
		DistributorCommision:       0.20,
		RetailerCommision:          0.60,
	}
}

func (db *Database) getCommisionQuery(
	ctx context.Context,
	retailerID string,
) (*models.GetCommisionModel, error) {

	var c models.GetCommisionModel
	const service = "PAYOUT"

	query := `
		SELECT
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision
		FROM commisions
		WHERE user_id = @user_id
		  AND service = @service
		LIMIT 1;
	`

	/* -------------------------------------------------------
	   1. Retailer commission (PAYOUT)
	------------------------------------------------------- */
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": retailerID,
		"service": service,
	}).Scan(
		&c.TotalCommision,
		&c.AdminCommision,
		&c.MasterDistributorCommision,
		&c.DistributorCommision,
		&c.RetailerCommision,
	)

	if err == nil {
		return &c, nil
	}

	/* -------------------------------------------------------
	   2. Resolve distributor & MD
	------------------------------------------------------- */
	var distributorID, mdID string
	hierarchyQuery := `
		SELECT
			d.distributor_id,
			d.master_distributor_id
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		WHERE r.retailer_id = @retailer_id;
	`

	err = db.pool.QueryRow(ctx, hierarchyQuery, pgx.NamedArgs{
		"retailer_id": retailerID,
	}).Scan(&distributorID, &mdID)

	if err != nil {
		return defaultPayoutCommision(), nil
	}

	/* -------------------------------------------------------
	   3. Distributor commission (PAYOUT)
	------------------------------------------------------- */
	err = db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": distributorID,
		"service": service,
	}).Scan(
		&c.TotalCommision,
		&c.AdminCommision,
		&c.MasterDistributorCommision,
		&c.DistributorCommision,
		&c.RetailerCommision,
	)

	if err == nil {
		return &c, nil
	}

	/* -------------------------------------------------------
	   4. Master Distributor commission (PAYOUT)
	------------------------------------------------------- */
	err = db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": mdID,
		"service": service,
	}).Scan(
		&c.TotalCommision,
		&c.AdminCommision,
		&c.MasterDistributorCommision,
		&c.DistributorCommision,
		&c.RetailerCommision,
	)

	if err == nil {
		return &c, nil
	}

	/* -------------------------------------------------------
	   5. Default PAYOUT commission
	------------------------------------------------------- */
	return defaultPayoutCommision(), nil
}

func (db *Database) GetPayoutCommisionSplit(
	ctx context.Context,
	retailerID string,
	amount float64,
) (*models.GetCommisionModel, error) {

	// 1. Get commission percentages (your existing logic)
	c, err := db.getCommisionQuery(ctx, retailerID)
	if err != nil {
		return nil, err
	}

	// 2. Helper: percent → amount
	calc := func(percent float64) float64 {
		return round2((amount * percent) / 100)
	}

	// 3. Overwrite percentage values with actual amounts
	c.TotalCommision = calc(c.TotalCommision)
	c.AdminCommision = calc(c.AdminCommision)
	c.MasterDistributorCommision = calc(c.MasterDistributorCommision)
	c.DistributorCommision = calc(c.DistributorCommision)
	c.RetailerCommision = calc(c.RetailerCommision)

	return c, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func (db *Database) CreatePayoutQuery(
	ctx context.Context,
	req models.CreatePayoutRequestModel,
	res models.PayoutAPIResponseModel,
	commision models.GetCommisionModel,
) {

}

// func (db *Database) CheckRetailerWalletBalance(
// 	ctx context.Context,
// 	retailerID string,
// 	amount float64,
// ) {
// 	query := `
// 		SELECT EXISTS (SELECT 1 FROM retailers WHERE retailer_id=@retailer_id AND retailer_wallet_balance=)
// 	`
// }
