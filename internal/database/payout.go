package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) ValidateRequestQuery(
	ctx context.Context,
	req models.CreatePayoutRequestModel,
	retailerCommision float64,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM retailers
			WHERE retailer_id=@retailer_id AND
			retailer_mpin=@retailer_mpin AND
			retailer_kyc_status=TRUE AND
			retailer_wallet_balance >= @amount
		);
	`
	var isValid bool
	if err := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"retailer_id":   req.RetailerID,
			"retailer_mpin": req.MPIN,
			"amount":        req.Amount + retailerCommision,
		},
	).Scan(&isValid); err != nil {
		return false, err
	}
	return isValid, nil
}

func (db *Database) GetPayoutCommisionQuery(
	ctx context.Context,
	retailerID string,
	service string,
) (*models.GetCommisionResponseModel, error) {

	var (
		distributorID string
		mdID          string
		comm          models.GetCommisionResponseModel
	)

	err := db.pool.QueryRow(ctx, `
		SELECT
			r.distributor_id,
			d.master_distributor_id
		FROM retailers r
		JOIN distributors d
			ON d.distributor_id = r.distributor_id
		WHERE r.retailer_id = @retailer_id;
	`, pgx.NamedArgs{
		"retailer_id": retailerID,
	}).Scan(&distributorID, &mdID)

	if err != nil {
		return nil, fmt.Errorf("invalid retailer hierarchy")
	}

	userIDs := []string{
		retailerID,
		distributorID,
		mdID,
	}

	query := `
		SELECT
			commision_id,
			user_id,
			service,
			total_commision,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			created_at,
			updated_at
		FROM commisions
		WHERE user_id = @user_id
		AND service = @service;
	`

	for _, userID := range userIDs {
		err = db.pool.QueryRow(ctx, query, pgx.NamedArgs{
			"user_id": userID,
			"service": service,
		}).Scan(
			&comm.CommisionID,
			&comm.UserID,
			&comm.Service,
			&comm.TotalCommision,
			&comm.AdminCommision,
			&comm.MasterDistributorCommision,
			&comm.DistributorCommision,
			&comm.RetailerCommision,
			&comm.CreatedAt,
			&comm.UpdatedAt,
		)

		if err == nil {
			return &comm, nil
		}
	}

	return &models.GetCommisionResponseModel{
		UserID:                     retailerID,
		Service:                    service,
		TotalCommision:             1.20,
		AdminCommision:             0.30,
		MasterDistributorCommision: 0.10,
		DistributorCommision:       0.20,
		RetailerCommision:          0.60,
	}, nil
}

func (db *Database) PayoutPendingOrSuccessQuery(
	ctx context.Context,
	payoutReq models.CreatePayoutRequestModel,
	commission models.GetCommisionResponseModel,
	res models.PayoutAPIResponseModel,
	adminID string,
	status string, // PENDING or SUCCESS
) error {

	if status != "PENDING" && status != "SUCCESS" {
		return fmt.Errorf("invalid payout status")
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// --------------------------------------------------
	// 1️⃣ Fetch Distributor & Master Distributor IDs
	// --------------------------------------------------
	var distributorID, mdID string
	err = tx.QueryRow(ctx, `
		SELECT r.distributor_id, d.master_distributor_id
		FROM retailers r
		JOIN distributors d ON d.distributor_id = r.distributor_id
		WHERE r.retailer_id = @rid
	`, pgx.NamedArgs{
		"rid": payoutReq.RetailerID,
	}).Scan(&distributorID, &mdID)

	if err != nil {
		return fmt.Errorf("invalid retailer hierarchy")
	}

	// --------------------------------------------------
	// 2️⃣ COMMISSION CALCULATION (MATCHES OLD SYSTEM)
	// --------------------------------------------------

	// Total commission amount (₹)
	// Example: 1000 * (1.2 / 100) = 12
	payoutAmount := payoutReq.Amount
	totalCommission := payoutAmount * (commission.TotalCommision / 100)

	// Split commission using FRACTIONS (NOT percentages)
	adminComm := totalCommission * commission.AdminCommision
	mdComm := totalCommission * commission.MasterDistributorCommision
	disComm := totalCommission * commission.DistributorCommision
	retailerComm := totalCommission * commission.RetailerCommision

	// TDS 2% (except admin)
	tdsRate := 0.02

	adminTDS := 0.0
	mdTDS := mdComm * tdsRate
	disTDS := disComm * tdsRate
	retailerTDS := retailerComm * tdsRate

	adminNet := adminComm
	mdNet := mdComm - mdTDS
	disNet := disComm - disTDS
	retailerNet := retailerComm - retailerTDS

	// --------------------------------------------------
	// 3️⃣ Insert payout transaction
	// --------------------------------------------------
	var payoutTxnID string
	err = tx.QueryRow(ctx, `
		INSERT INTO payout_transactions (
			partner_request_id,
			operator_transaction_id,
			payout_transaction_status,
			retailer_id,
			order_id,
			mobile_number,
			beneficiary_bank_name,
			beneficiary_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			amount,
			transfer_type,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision
		) VALUES (
			@partner,
			@operator,
			@status,
			@retailer,
			@order_id,
			@mobile,
			@bank,
			@name,
			@account,
			@ifsc,
			@amount,
			@type,
			@admin_comm,
			@md_comm,
			@dis_comm,
			@ret_comm
		)
		RETURNING payout_transaction_id
	`, pgx.NamedArgs{
		"partner":    payoutReq.PartnerRequestID,
		"operator":   res.OperatorTransactionID,
		"status":     status,
		"retailer":   payoutReq.RetailerID,
		"order_id":   res.OrderID,
		"mobile":     payoutReq.MobileNumber,
		"bank":       payoutReq.BeneficiaryBankName,
		"name":       payoutReq.BeneficiaryName,
		"account":    payoutReq.BeneficiaryAccountNumber,
		"ifsc":       payoutReq.BeneficiaryIFSCCode,
		"amount":     payoutAmount,
		"type":       payoutReq.TransferType,
		"admin_comm": adminComm,
		"md_comm":    mdComm,
		"dis_comm":   disComm,
		"ret_comm":   retailerComm,
	}).Scan(&payoutTxnID)

	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 4️⃣ Deduct payout + total commission from retailer
	// --------------------------------------------------
	totalDebit := payoutAmount + totalCommission

	var retailerBefore float64
	err = tx.QueryRow(ctx, `
		SELECT retailer_wallet_balance
		FROM retailers
		WHERE retailer_id = @id
		FOR UPDATE
	`, pgx.NamedArgs{
		"id": payoutReq.RetailerID,
	}).Scan(&retailerBefore)

	if err != nil {
		return err
	}

	retailerAfter := retailerBefore - totalDebit
	if retailerAfter < 0 {
		return fmt.Errorf("insufficient retailer wallet balance")
	}

	_, err = tx.Exec(ctx, `
		UPDATE retailers
		SET retailer_wallet_balance = @bal
		WHERE retailer_id = @id
	`, pgx.NamedArgs{
		"id":  payoutReq.RetailerID,
		"bal": retailerAfter,
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			debit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@uid,
			@ref,
			@amt,
			@before,
			@after,
			'PAYOUT',
			'Payout + commission debit'
		)
	`, pgx.NamedArgs{
		"uid":    payoutReq.RetailerID,
		"ref":    payoutTxnID,
		"amt":    totalDebit,
		"before": retailerBefore,
		"after":  retailerAfter,
	})
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 5️⃣ Credit commission wallets + TDS
	// --------------------------------------------------
	type credit struct {
		UserID string
		Role   string
		Gross  float64
		Net    float64
		TDS    float64
	}

	credits := []credit{
		{adminID, "ADMIN", adminComm, adminNet, adminTDS},
		{mdID, "MASTER_DISTRIBUTOR", mdComm, mdNet, mdTDS},
		{distributorID, "DISTRIBUTOR", disComm, disNet, disTDS},
		{payoutReq.RetailerID, "RETAILER", retailerComm, retailerNet, retailerTDS},
	}

	for _, c := range credits {
		if c.UserID == "" || c.Gross <= 0 {
			continue
		}

		var before float64
		var pan string

		switch c.Role {

		case "ADMIN":
			err = tx.QueryRow(ctx,
				`SELECT admin_wallet_balance FROM admins WHERE admin_id=@id FOR UPDATE`,
				pgx.NamedArgs{"id": c.UserID},
			).Scan(&before)

		case "MASTER_DISTRIBUTOR":
			err = tx.QueryRow(ctx,
				`SELECT master_distributor_pan_number, master_distributor_wallet_balance
				 FROM master_distributors WHERE master_distributor_id=@id FOR UPDATE`,
				pgx.NamedArgs{"id": c.UserID},
			).Scan(&pan, &before)

		case "DISTRIBUTOR":
			err = tx.QueryRow(ctx,
				`SELECT distributor_pan_number, distributor_wallet_balance
				 FROM distributors WHERE distributor_id=@id FOR UPDATE`,
				pgx.NamedArgs{"id": c.UserID},
			).Scan(&pan, &before)

		case "RETAILER":
			err = tx.QueryRow(ctx,
				`SELECT retailer_pan_number, retailer_wallet_balance
				 FROM retailers WHERE retailer_id=@id FOR UPDATE`,
				pgx.NamedArgs{"id": c.UserID},
			).Scan(&pan, &before)
		}

		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return err
		}

		after := before + c.Net

		var update string
		switch c.Role {
		case "ADMIN":
			update = `UPDATE admins SET admin_wallet_balance=@bal WHERE admin_id=@id`
		case "MASTER_DISTRIBUTOR":
			update = `UPDATE master_distributors SET master_distributor_wallet_balance=@bal WHERE master_distributor_id=@id`
		case "DISTRIBUTOR":
			update = `UPDATE distributors SET distributor_wallet_balance=@bal WHERE distributor_id=@id`
		case "RETAILER":
			update = `UPDATE retailers SET retailer_wallet_balance=@bal WHERE retailer_id=@id`
		}

		_, err = tx.Exec(ctx, update, pgx.NamedArgs{
			"id":  c.UserID,
			"bal": after,
		})
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO wallet_transactions (
				user_id,
				reference_id,
				credit_amount,
				before_balance,
				after_balance,
				transaction_reason,
				remarks
			) VALUES (
				@uid,
				@ref,
				@amt,
				@before,
				@after,
				'PAYOUT',
				'Commission credited'
			)
		`, pgx.NamedArgs{
			"uid":    c.UserID,
			"ref":    payoutTxnID,
			"amt":    c.Net,
			"before": before,
			"after":  after,
		})
		if err != nil {
			return err
		}

		if c.Role != "ADMIN" {
			_, err = tx.Exec(ctx, `
				INSERT INTO tds_commision (
					transaction_id,
					user_id,
					user_name,
					commision,
					tds,
					paid_commision,
					pan_number,
					status
				) VALUES (
					@txn,
					@uid,
					@role,
					@gross,
					@tds,
					@net,
					@pan,
					'DEDUCTED'
				)
			`, pgx.NamedArgs{
				"txn":   payoutTxnID,
				"uid":   c.UserID,
				"role":  c.Role,
				"gross": c.Gross,
				"tds":   c.TDS,
				"net":   c.Net,
				"pan":   pan,
			})
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}


func (db *Database) PayoutFailedQuery(
	ctx context.Context,
	req models.CreatePayoutRequestModel,
) error {

	_, err := db.pool.Exec(ctx, `
		INSERT INTO payout_transactions (
			partner_request_id,
			operator_transaction_id,
			payout_transaction_status,
			retailer_id,
			order_id,
			mobile_number,
			beneficiary_bank_name,
			beneficiary_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			amount,
			transfer_type,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision
		) VALUES (
			@partner_req_id,
			@operator_txn_id,
			'FAILED',
			@retailer_id,
			@order_id,
			@mobile_number,
			@bank_name,
			@beneficiary_name,
			@account_number,
			@ifsc,
			@amount,
			@transfer_type,
			0, 0, 0, 0
		);
	`, pgx.NamedArgs{
		"partner_req_id":   req.PartnerRequestID,
		"operator_txn_id":  "", // empty or failure code from operator
		"retailer_id":      req.RetailerID,
		"order_id":         req.PartnerRequestID,
		"mobile_number":    req.MobileNumber,
		"bank_name":        req.BeneficiaryBankName,
		"beneficiary_name": req.BeneficiaryName,
		"account_number":   req.BeneficiaryAccountNumber,
		"ifsc":             req.BeneficiaryIFSCCode,
		"amount":           req.Amount,
		"transfer_type":    req.TransferType,
	})

	return err
}
