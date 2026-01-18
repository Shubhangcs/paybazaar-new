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
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	/* -------------------------------------------------------
	   1. Decide payout status
	------------------------------------------------------- */
	status := "FAILED"
	if res.Status == 1 || res.Status == 2 {
		status = "SUCCESS"
	}

	/* -------------------------------------------------------
	   2. Insert payout_service (always)
	------------------------------------------------------- */
	_, err = tx.Exec(ctx, `
		INSERT INTO payout_service (
			partner_request_id,
			operator_transaction_id,
			retailer_id,
			order_id,
			mobile_number,
			beneficiary_bank_name,
			beneficiary_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			amount,
			transfer_type
		) VALUES (
			@partner_request_id,
			@operator_txn_id,
			@retailer_id,
			@order_id,
			@mobile,
			@bank,
			@name,
			@account,
			@ifsc,
			@amount,
			@type
		);
	`, pgx.NamedArgs{
		"partner_request_id": res.PartnerRequestID,
		"operator_txn_id":    res.OperatorTransactionID,
		"retailer_id":        req.RetailerID,
		"order_id":           res.OrderID,
		"mobile":             req.MobileNumber,
		"bank":               req.BeneficiaryBankName,
		"name":               req.BeneficiaryName,
		"account":            req.BeneficiaryAccountNumber,
		"ifsc":               req.BeneficiaryIFSCCode,
		"amount":             req.Amount,
		"type":               req.TransferType,
	})
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   3. If FAILED → stop here
	------------------------------------------------------- */
	if status == "FAILED" {
		return tx.Commit(ctx)
	}

	/* -------------------------------------------------------
	   4. Resolve hierarchy
	------------------------------------------------------- */
	var distributorID, mdID, adminID string
	err = tx.QueryRow(ctx, `
		SELECT
			d.distributor_id,
			md.master_distributor_id,
			md.admin_id
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
		WHERE r.retailer_id = @rid;
	`, pgx.NamedArgs{
		"rid": req.RetailerID,
	}).Scan(&distributorID, &mdID, &adminID)
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   5. Deduct retailer (amount + commission)
	------------------------------------------------------- */
	var retailerBefore float64
	err = tx.QueryRow(ctx, `
		SELECT retailer_wallet_balance
		FROM retailers
		WHERE retailer_id = @id
		FOR UPDATE;
	`, pgx.NamedArgs{"id": req.RetailerID}).Scan(&retailerBefore)
	if err != nil {
		return err
	}

	totalDebit := req.Amount + commision.TotalCommision
	retailerAfter := retailerBefore - totalDebit

	_, err = tx.Exec(ctx, `
		UPDATE retailers
		SET retailer_wallet_balance = @bal, updated_at = NOW()
		WHERE retailer_id = @id;
	`, pgx.NamedArgs{
		"id":  req.RetailerID,
		"bal": retailerAfter,
	})
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   6. Wallet transaction (retailer debit)
	------------------------------------------------------- */
	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions (
			user_id, reference_id, debit_amount,
			before_balance, after_balance,
			transaction_reason, remarks
		) VALUES (
			@uid, @ref, @amt,
			@before, @after,
			'PAYOUT', 'Payout transaction'
		);
	`, pgx.NamedArgs{
		"uid":    req.RetailerID,
		"ref":    res.OperatorTransactionID,
		"amt":    totalDebit,
		"before": retailerBefore,
		"after":  retailerAfter,
	})
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   7. Credit commissions (ADMIN / MD / DIS / RET)
	------------------------------------------------------- */
	type credit struct {
		id     string
		amount float64
		table  string
	}

	credits := []credit{
		{adminID, commision.AdminCommision, "admins"},
		{mdID, commision.MasterDistributorCommision, "master_distributors"},
		{distributorID, commision.DistributorCommision, "distributors"},
		{req.RetailerID, commision.RetailerCommision, "retailers"},
	}

	for _, c := range credits {
		if c.amount <= 0 {
			continue
		}

		var before float64
		err := tx.QueryRow(ctx, `
			SELECT `+c.table+`_wallet_balance
			FROM `+c.table+`
			WHERE `+c.table+`_id = @id
			FOR UPDATE;
		`, pgx.NamedArgs{"id": c.id}).Scan(&before)
		if err != nil {
			return err
		}

		after := before + c.amount

		_, err = tx.Exec(ctx, `
			UPDATE `+c.table+`
			SET `+c.table+`_wallet_balance = @bal, updated_at = NOW()
			WHERE `+c.table+`_id = @id;
		`, pgx.NamedArgs{
			"id":  c.id,
			"bal": after,
		})
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO wallet_transactions (
				user_id, reference_id, credit_amount,
				before_balance, after_balance,
				transaction_reason, remarks
			) VALUES (
				@uid, @ref, @amt,
				@before, @after,
				'COMMISSION', 'Payout commission'
			);
		`, pgx.NamedArgs{
			"uid":    c.id,
			"ref":    res.OperatorTransactionID,
			"amt":    c.amount,
			"before": before,
			"after":  after,
		})
		if err != nil {
			return err
		}
	}

	/* -------------------------------------------------------
	   8. Commit
	------------------------------------------------------- */
	return tx.Commit(ctx)
}
