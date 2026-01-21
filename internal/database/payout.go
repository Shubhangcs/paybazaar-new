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
	status string,
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
	// 1️⃣ Fetch hierarchy + balances + PAN + Names (LOCKED)
	// --------------------------------------------------
	type users struct {
		AdminID   string
		AdminName string
		AdminBal  float64

		MDID   string
		MDName string
		MDPAN  string
		MDBal  float64

		DisID   string
		DisName string
		DisPAN  string
		DisBal  float64

		RetID   string
		RetName string
		RetPAN  string
		RetBal  float64
	}

	var u users

	err = tx.QueryRow(ctx, `
		SELECT
			a.admin_id,
			a.admin_name,
			a.admin_wallet_balance,

			md.master_distributor_id,
			md.master_distributor_name,
			md.master_distributor_pan_number,
			md.master_distributor_wallet_balance,

			d.distributor_id,
			d.distributor_name,
			d.distributor_pan_number,
			d.distributor_wallet_balance,

			r.retailer_id,
			r.retailer_name,
			r.retailer_pan_number,
			r.retailer_wallet_balance
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
		JOIN admins a ON md.admin_id = a.admin_id
		WHERE r.retailer_id = @rid
		FOR UPDATE
	`, pgx.NamedArgs{
		"rid": payoutReq.RetailerID,
	}).Scan(
		&u.AdminID, &u.AdminName, &u.AdminBal,
		&u.MDID, &u.MDName, &u.MDPAN, &u.MDBal,
		&u.DisID, &u.DisName, &u.DisPAN, &u.DisBal,
		&u.RetID, &u.RetName, &u.RetPAN, &u.RetBal,
	)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 2️⃣ Commission Calculation (CORRECTED)
	// --------------------------------------------------
	amount := payoutReq.Amount

	// Commission percentages are already in percentage form (e.g., 2.5 means 2.5%)
	// TotalCommision is also in percentage
	totalCommPercentage := commission.TotalCommision / 100
	totalComm := amount * totalCommPercentage

	// Each user gets their share of total commission (these are fractions that sum to 1.0)
	adminComm := totalComm * commission.AdminCommision
	mdComm := totalComm * commission.MasterDistributorCommision
	disComm := totalComm * commission.DistributorCommision
	retComm := totalComm * commission.RetailerCommision

	// TDS is 2% (0.02) of each person's commission
	tdsRate := 0.02

	mdTDS := mdComm * tdsRate
	disTDS := disComm * tdsRate
	retTDS := retComm * tdsRate

	// Net commission after TDS deduction
	mdNet := mdComm - mdTDS
	disNet := disComm - disTDS
	retNet := retComm - retTDS

	// Admin gets their commission + all collected TDS
	adminFinal := adminComm + mdTDS + disTDS + retTDS

	// --------------------------------------------------
	// 3️⃣ Insert payout transaction
	// --------------------------------------------------
	var payoutTxnID string
	err = tx.QueryRow(ctx, `
		INSERT INTO payout_transactions (
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
			transfer_type,
			admin_commision,
			master_distributor_commision,
			distributor_commision,
			retailer_commision,
			payout_transaction_status
		) VALUES (
			@partner, @operator, @rid, @order,
			@mobile, @bank, @name, @account, @ifsc,
			@amount, @type,
			@admin, @md, @dis, @ret,
			@status
		)
		RETURNING payout_transaction_id::TEXT
	`, pgx.NamedArgs{
		"partner":  payoutReq.PartnerRequestID,
		"operator": res.OperatorTransactionID,
		"rid":      payoutReq.RetailerID,
		"order":    res.OrderID,
		"mobile":   payoutReq.MobileNumber,
		"bank":     payoutReq.BeneficiaryBankName,
		"name":     payoutReq.BeneficiaryName,
		"account":  payoutReq.BeneficiaryAccountNumber,
		"ifsc":     payoutReq.BeneficiaryIFSCCode,
		"amount":   amount,
		"type":     payoutReq.TransferType,
		"admin":    adminComm, // Store admin's commission only (not including TDS)
		"md":       mdComm,
		"dis":      disComm,
		"ret":      retComm,
		"status":   status,
	}).Scan(&payoutTxnID)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 4️⃣ Retailer Debit (amount + total commission)
	// --------------------------------------------------
	totalDebit := amount + totalComm
	if u.RetBal < totalDebit {
		return fmt.Errorf("insufficient retailer balance")
	}

	_, err = tx.Exec(ctx, `
		UPDATE retailers
		SET retailer_wallet_balance = retailer_wallet_balance - @amt
		WHERE retailer_id = @id
	`, pgx.NamedArgs{"amt": totalDebit, "id": u.RetID})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions (
			user_id, reference_id,
			debit_amount, before_balance, after_balance,
			transaction_reason, remarks
		) VALUES (
			@uid, @ref,
			@amt, @before, @after,
			'PAYOUT', 'Payout + commission debit'
		)
	`, pgx.NamedArgs{
		"uid":    u.RetID,
		"ref":    payoutTxnID,
		"amt":    totalDebit,
		"before": u.RetBal,
		"after":  u.RetBal - totalDebit,
	})
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 5️⃣ Admin Credit (commission + ALL TDS)
	// --------------------------------------------------
	_, err = tx.Exec(ctx, `
		UPDATE admins
		SET admin_wallet_balance = admin_wallet_balance + @amt
		WHERE admin_id = @id
	`, pgx.NamedArgs{"amt": adminFinal, "id": u.AdminID})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions (
			user_id, reference_id,
			credit_amount, before_balance, after_balance,
			transaction_reason, remarks
		) VALUES (
			@uid, @ref,
			@amt, @before, @after,
			'PAYOUT', 'Admin commission + TDS collection'
		)
	`, pgx.NamedArgs{
		"uid":    u.AdminID,
		"ref":    payoutTxnID,
		"amt":    adminFinal,
		"before": u.AdminBal,
		"after":  u.AdminBal + adminFinal,
	})
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 6️⃣ MD Credits + TDS Entry
	// --------------------------------------------------
	newMDBal := u.MDBal + mdNet

	_, err = tx.Exec(ctx, `
		UPDATE master_distributors 
		SET master_distributor_wallet_balance = master_distributor_wallet_balance + @amt 
		WHERE master_distributor_id = @id
	`, pgx.NamedArgs{"amt": mdNet, "id": u.MDID})
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
			'PAYOUT', 'MD commission (after TDS)'
		)
	`, pgx.NamedArgs{
		"uid":    u.MDID,
		"ref":    payoutTxnID,
		"amt":    mdNet,
		"before": u.MDBal,
		"after":  newMDBal,
	})
	if err != nil {
		return err
	}

	// TDS entry for MD
	_, err = tx.Exec(ctx, `
		INSERT INTO tds_commision (
			transaction_id, user_id, user_name,
			commision, tds, paid_commision,
			pan_number, status
		) VALUES (
			@txn, @uid, @name,
			@comm, @tds, @paid,
			@pan, @status
		)
	`, pgx.NamedArgs{
		"txn":    payoutTxnID,
		"uid":    u.MDID,
		"name":   u.MDName,
		"comm":   mdComm,
		"tds":    mdTDS,
		"paid":   mdNet,
		"pan":    u.MDPAN,
		"status": status,
	})
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 7️⃣ Distributor Credits + TDS Entry
	// --------------------------------------------------
	newDisBal := u.DisBal + disNet

	_, err = tx.Exec(ctx, `
		UPDATE distributors 
		SET distributor_wallet_balance = distributor_wallet_balance + @amt 
		WHERE distributor_id = @id
	`, pgx.NamedArgs{"amt": disNet, "id": u.DisID})
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
			'PAYOUT', 'Distributor commission (after TDS)'
		)
	`, pgx.NamedArgs{
		"uid":    u.DisID,
		"ref":    payoutTxnID,
		"amt":    disNet,
		"before": u.DisBal,
		"after":  newDisBal,
	})
	if err != nil {
		return err
	}

	// TDS entry for Distributor
	_, err = tx.Exec(ctx, `
		INSERT INTO tds_commision (
			transaction_id, user_id, user_name,
			commision, tds, paid_commision,
			pan_number, status
		) VALUES (
			@txn, @uid, @name,
			@comm, @tds, @paid,
			@pan, @status
		)
	`, pgx.NamedArgs{
		"txn":    payoutTxnID,
		"uid":    u.DisID,
		"name":   u.DisName,
		"comm":   disComm,
		"tds":    disTDS,
		"paid":   disNet,
		"pan":    u.DisPAN,
		"status": status,
	})
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// 8️⃣ Retailer Commission Credit + TDS Entry
	// --------------------------------------------------
	// Note: Retailer balance was already debited by totalDebit earlier
	beforeRetComm := u.RetBal - totalDebit
	afterRetComm := beforeRetComm + retNet

	_, err = tx.Exec(ctx, `
		UPDATE retailers 
		SET retailer_wallet_balance = retailer_wallet_balance + @amt 
		WHERE retailer_id = @id
	`, pgx.NamedArgs{"amt": retNet, "id": u.RetID})
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
			'PAYOUT', 'Retailer commission (after TDS)'
		)
	`, pgx.NamedArgs{
		"uid":    u.RetID,
		"ref":    payoutTxnID,
		"amt":    retNet,
		"before": beforeRetComm,
		"after":  afterRetComm,
	})
	if err != nil {
		return err
	}

	// TDS entry for Retailer
	_, err = tx.Exec(ctx, `
		INSERT INTO tds_commision (
			transaction_id, user_id, user_name,
			commision, tds, paid_commision,
			pan_number, status
		) VALUES (
			@txn, @uid, @name,
			@comm, @tds, @paid,
			@pan, @status
		)
	`, pgx.NamedArgs{
		"txn":    payoutTxnID,
		"uid":    u.RetID,
		"name":   u.RetName,
		"comm":   retComm,
		"tds":    retTDS,
		"paid":   retNet,
		"pan":    u.RetPAN,
		"status": status,
	})
	if err != nil {
		return err
	}

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
