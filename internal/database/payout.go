package database

import (
	"context"
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
			retailer_wallet_balance >= (@amount + @retailer_commision);
		);
	`
	var isValid bool
	if err := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"retailer_id":        req.RetailerID,
			"retailer_mpin":      req.MPIN,
			"amount":             req.Amount,
			"retailer_commision": retailerCommision,
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
	comm models.GetCommisionResponseModel,
	adminID string,
	status string, // must be PENDING or SUCCESS
) error {

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
		WHERE r.retailer_id = @rid;
	`, pgx.NamedArgs{
		"rid": payoutReq.RetailerID,
	}).Scan(&distributorID, &mdID)

	if err != nil {
		return fmt.Errorf("invalid retailer hierarchy")
	}

	// --------------------------------------------------
	// 2️⃣ Insert payout transaction (WITH STATUS)
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
			@operator_txn,
			@status,
			@retailer,
			@order,
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
		RETURNING payout_transaction_id;
	`, pgx.NamedArgs{
		"partner":      payoutReq.PartnerRequestID,
		"operator_txn": "", // fill later when operator responds
		"status":       status, // PENDING or SUCCESS
		"retailer":     payoutReq.RetailerID,
		"order":        payoutReq.PartnerRequestID,
		"mobile":       payoutReq.MobileNumber,
		"bank":         payoutReq.BeneficiaryBankName,
		"name":         payoutReq.BeneficiaryName,
		"account":      payoutReq.BeneficiaryAccountNumber,
		"ifsc":         payoutReq.BeneficiaryIFSCCode,
		"amount":       payoutReq.Amount,
		"type":         payoutReq.TransferType,
		"admin_comm":   comm.AdminCommision,
		"md_comm":      comm.MasterDistributorCommision,
		"dis_comm":     comm.DistributorCommision,
		"ret_comm":     comm.RetailerCommision,
	}).Scan(&payoutTxnID)

	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 3️⃣ Deduct Amount + Total Commission from Retailer
	// --------------------------------------------------
	totalDebit := payoutReq.Amount + comm.TotalCommision

	var beforeRetailerBalance float64
	err = tx.QueryRow(ctx, `
		SELECT retailer_wallet_balance
		FROM retailers
		WHERE retailer_id = @id
		FOR UPDATE;
	`, pgx.NamedArgs{"id": payoutReq.RetailerID}).Scan(&beforeRetailerBalance)

	if err != nil {
		return err
	}

	afterRetailerBalance := beforeRetailerBalance - totalDebit

	_, err = tx.Exec(ctx, `
		UPDATE retailers
		SET retailer_wallet_balance = @bal
		WHERE retailer_id = @id;
	`, pgx.NamedArgs{
		"id":  payoutReq.RetailerID,
		"bal": afterRetailerBalance,
	})
	if err != nil {
		return err
	}

	// Retailer debit wallet entry
	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions (
			user_id, reference_id, debit_amount,
			before_balance, after_balance,
			transaction_reason, remarks
		) VALUES (
			@uid, @ref, @amt, @before, @after,
			'PAYOUT', 'Payout initiated'
		);
	`, pgx.NamedArgs{
		"uid":    payoutReq.RetailerID,
		"ref":    payoutTxnID,
		"amt":    totalDebit,
		"before": beforeRetailerBalance,
		"after":  afterRetailerBalance,
	})
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 4️⃣ Commission Distribution with PAN + TDS
	// --------------------------------------------------
	type split struct {
		UserID string
		Amount float64
		Role   string
	}

	splits := []split{
		{adminID, comm.AdminCommision, "ADMIN"},
		{mdID, comm.MasterDistributorCommision, "MASTER_DISTRIBUTOR"},
		{distributorID, comm.DistributorCommision, "DISTRIBUTOR"},
		{payoutReq.RetailerID, comm.RetailerCommision, "RETAILER"},
	}

	tdsRate := 0.02

	for _, s := range splits {
		if s.Amount <= 0 {
			continue
		}

		var before, after, tds, net float64
		var pan string

		switch s.Role {

		case "ADMIN":
			err = tx.QueryRow(ctx, `
				SELECT admin_wallet_balance
				FROM admins WHERE admin_id=@id FOR UPDATE;
			`, pgx.NamedArgs{"id": s.UserID}).Scan(&before)
			tds = 0
			net = s.Amount

		case "MASTER_DISTRIBUTOR":
			err = tx.QueryRow(ctx, `
				SELECT master_distributor_pan_number, master_distributor_wallet_balance
				FROM master_distributors WHERE master_distributor_id=@id FOR UPDATE;
			`, pgx.NamedArgs{"id": s.UserID}).Scan(&pan, &before)
			tds = s.Amount * tdsRate
			net = s.Amount - tds

		case "DISTRIBUTOR":
			err = tx.QueryRow(ctx, `
				SELECT distributor_pan_number, distributor_wallet_balance
				FROM distributors WHERE distributor_id=@id FOR UPDATE;
			`, pgx.NamedArgs{"id": s.UserID}).Scan(&pan, &before)
			tds = s.Amount * tdsRate
			net = s.Amount - tds

		case "RETAILER":
			err = tx.QueryRow(ctx, `
				SELECT retailer_pan_number, retailer_wallet_balance
				FROM retailers WHERE retailer_id=@id FOR UPDATE;
			`, pgx.NamedArgs{"id": s.UserID}).Scan(&pan, &before)
			tds = s.Amount * tdsRate
			net = s.Amount - tds
		}

		if err != nil {
			return err
		}

		after = before + net

		var update string
		switch s.Role {
		case "ADMIN":
			update = `UPDATE admins SET admin_wallet_balance=@bal WHERE admin_id=@id`
		case "MASTER_DISTRIBUTOR":
			update = `UPDATE master_distributors SET master_distributor_wallet_balance=@bal WHERE master_distributor_id=@id`
		case "DISTRIBUTOR":
			update = `UPDATE distributors SET distributor_wallet_balance=@bal WHERE distributor_id=@id`
		case "RETAILER":
			update = `UPDATE retailers SET retailer_wallet_balance=@bal WHERE retailer_id=@id`
		}

		_, err = tx.Exec(ctx, update, pgx.NamedArgs{"id": s.UserID, "bal": after})
		if err != nil {
			return err
		}

		// Wallet credit entry
		_, err = tx.Exec(ctx, `
			INSERT INTO wallet_transactions (
				user_id, reference_id, credit_amount,
				before_balance, after_balance,
				transaction_reason, remarks
			) VALUES (
				@uid, @ref, @amt, @before, @after,
				'PAYOUT', 'Commission credited'
			);
		`, pgx.NamedArgs{
			"uid":    s.UserID,
			"ref":    payoutTxnID,
			"amt":    net,
			"before": before,
			"after":  after,
		})
		if err != nil {
			return err
		}

		// TDS entry (except admin)
		if s.Role != "ADMIN" {
			_, err = tx.Exec(ctx, `
				INSERT INTO tds_commision (
					transaction_id, user_id, user_name,
					commision, tds, paid_commision,
					pan_number, status
				) VALUES (
					@txn, @uid, @role,
					@comm, @tds, @paid,
					@pan, 'DEDUCTED'
				);
			`, pgx.NamedArgs{
				"txn":  payoutTxnID,
				"uid":  s.UserID,
				"role": s.Role,
				"comm": s.Amount,
				"tds":  tds,
				"paid": net,
				"pan":  pan,
			})
			if err != nil {
				return err
			}
		}
	}

	// --------------------------------------------------
	// 5️⃣ Commit Transaction
	// --------------------------------------------------
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
