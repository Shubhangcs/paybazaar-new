package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateFundRequestQuery(
	ctx context.Context,
	req models.CreateFundRequestModel,
) (int64, error) {

	query := `
		INSERT INTO fund_requests (
			requester_id,
			request_to_id,
			amount,
			bank_name,
			request_date,
			utr_number,
			request_status,
			remarks,
			reject_remarks
		) VALUES (
			@requester_id,
			@request_to_id,
			@amount,
			@bank_name,
			@request_date,
			@utr_number,
			'PENDING',
			@remarks,
			''
		)
		RETURNING fund_request_id;
	`

	var id int64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"requester_id":  req.RequesterID,
		"request_to_id": req.RequestToID,
		"amount":        req.Amount,
		"bank_name":     req.BankName,
		"request_date":  req.RequestDate,
		"utr_number":    req.UTRNumber,
		"remarks":       req.Remarks,
	}).Scan(&id)

	return id, err
}

func (db *Database) GetFundRequestQuery(
	ctx context.Context,
	fundRequestID int64,
) (*models.GetFundRequestResponseModel, error) {

	query := `
		SELECT
			fund_request_id,
			requester_id,
			request_to_id,
			amount,
			bank_name,
			request_date,
			utr_number,
			request_status,
			remarks,
			reject_remarks,
			created_at,
			updated_at
		FROM fund_requests
		WHERE fund_request_id = @id;
	`

	var fr models.GetFundRequestResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id": fundRequestID,
	}).Scan(
		&fr.FundRequestID,
		&fr.RequesterID,
		&fr.RequestToID,
		&fr.Amount,
		&fr.BankName,
		&fr.RequestDate,
		&fr.UTRNumber,
		&fr.RequestStatus,
		&fr.Remarks,
		&fr.RejectRemarks,
		&fr.CreatedAt,
		&fr.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &fr, nil
}

func (db *Database) GetAllFundRequestsQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetFundRequestResponseModel, error) {

	query := `
		SELECT
			fund_request_id,
			requester_id,
			request_to_id,
			amount,
			bank_name,
			request_date,
			utr_number,
			request_status,
			remarks,
			reject_remarks,
			created_at,
			updated_at
		FROM fund_requests
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetFundRequestResponseModel

	for rows.Next() {
		var fr models.GetFundRequestResponseModel
		if err := rows.Scan(
			&fr.FundRequestID,
			&fr.RequesterID,
			&fr.RequestToID,
			&fr.Amount,
			&fr.BankName,
			&fr.RequestDate,
			&fr.UTRNumber,
			&fr.RequestStatus,
			&fr.Remarks,
			&fr.RejectRemarks,
			&fr.CreatedAt,
			&fr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, fr)
	}

	return list, rows.Err()
}

func (db *Database) getFundRequestsByUser(
	ctx context.Context,
	column string,
	req models.GetFundRequestFilterRequestModel,
	limit, offset int,
) ([]models.GetFundRequestResponseModel, error) {

	query := fmt.Sprintf(`
		SELECT
			fund_request_id,
			requester_id,
			request_to_id,
			amount,
			bank_name,
			request_date,
			utr_number,
			request_status,
			remarks,
			reject_remarks,
			created_at,
			updated_at
		FROM fund_requests
		WHERE %s = @id
	`, column)

	args := pgx.NamedArgs{
		"id":     req.ID,
		"limit":  limit,
		"offset": offset,
	}

	if req.StartDate != nil {
		query += ` AND created_at >= @start_date`
		args["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		query += ` AND created_at <= @end_date`
		args["end_date"] = *req.EndDate
	}

	if req.Status != nil {
		query += ` AND request_status = @status`
		args["status"] = *req.Status
	}

	query += ` ORDER BY created_at DESC LIMIT @limit OFFSET @offset;`

	rows, err := db.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetFundRequestResponseModel

	for rows.Next() {
		var fr models.GetFundRequestResponseModel
		if err := rows.Scan(
			&fr.FundRequestID,
			&fr.RequesterID,
			&fr.RequestToID,
			&fr.Amount,
			&fr.BankName,
			&fr.RequestDate,
			&fr.UTRNumber,
			&fr.RequestStatus,
			&fr.Remarks,
			&fr.RejectRemarks,
			&fr.CreatedAt,
			&fr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, fr)
	}

	return list, rows.Err()
}

func (db *Database) GetFundRequestsByRequesterIDQuery(
	ctx context.Context,
	req models.GetFundRequestFilterRequestModel,
	limit, offset int,
) ([]models.GetFundRequestResponseModel, error) {
	return db.getFundRequestsByUser(ctx, "requester_id", req, limit, offset)
}

func (db *Database) GetFundRequestsByRequestToIDQuery(
	ctx context.Context,
	req models.GetFundRequestFilterRequestModel,
	limit, offset int,
) ([]models.GetFundRequestResponseModel, error) {
	return db.getFundRequestsByUser(ctx, "request_to_id", req, limit, offset)
}

func (db *Database) RejectFundRequestQuery(
	ctx context.Context,
	fundRequestID int64,
	rejectRemarks string,
) error {

	query := `
		UPDATE fund_requests
		SET
			request_status = 'REJECTED',
			reject_remarks = @remarks,
			updated_at = NOW()
		WHERE fund_request_id = @id
		  AND request_status = 'PENDING';
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"id":      fundRequestID,
		"remarks": rejectRemarks,
	})

	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) AcceptFundRequestQuery(
	ctx context.Context,
	fundRequestID int64,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var fr struct {
		RequesterID   string
		RequestToID   string
		Amount        float64
		RequestStatus string
	}

	err = tx.QueryRow(ctx, `
		SELECT requester_id, request_to_id, amount, request_status
		FROM fund_requests
		WHERE fund_request_id = @id
		FOR UPDATE;
	`, pgx.NamedArgs{"id": fundRequestID}).Scan(
		&fr.RequesterID,
		&fr.RequestToID,
		&fr.Amount,
		&fr.RequestStatus,
	)

	if err != nil {
		return err
	}

	if fr.RequestStatus != "PENDING" {
		return errors.New("fund request already processed")
	}

	senderTable, err := db.getUserTable(fr.RequestToID)
	if err != nil {
		return err
	}
	receiverTable, err := db.getUserTable(fr.RequesterID)
	if err != nil {
		return err
	}

	var senderBalance, receiverBalance float64

	err = tx.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s_wallet_balance FROM %ss WHERE %s_id=@id`,
			senderTable, senderTable, senderTable),
		pgx.NamedArgs{"id": fr.RequestToID},
	).Scan(&senderBalance)
	if err != nil {
		return err
	}

	if senderBalance < fr.Amount {
		return errors.New("insufficient balance")
	}

	err = tx.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s_wallet_balance FROM %ss WHERE %s_id=@id`,
			receiverTable, receiverTable, receiverTable),
		pgx.NamedArgs{"id": fr.RequesterID},
	).Scan(&receiverBalance)
	if err != nil {
		return err
	}

	senderAfter := senderBalance - fr.Amount
	receiverAfter := receiverBalance + fr.Amount

	_, err = tx.Exec(ctx,
		fmt.Sprintf(`UPDATE %ss SET %s_wallet_balance=@b WHERE %s_id=@id`,
			senderTable, senderTable, senderTable),
		pgx.NamedArgs{"b": senderAfter, "id": fr.RequestToID},
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		fmt.Sprintf(`UPDATE %ss SET %s_wallet_balance=@b WHERE %s_id=@id`,
			receiverTable, receiverTable, receiverTable),
		pgx.NamedArgs{"b": receiverAfter, "id": fr.RequesterID},
	)
	if err != nil {
		return err
	}

	ref := fmt.Sprintf("%d", fundRequestID)

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions
		(user_id, reference_id, debit_amount, before_balance, after_balance, transaction_reason, remarks)
		VALUES (@u,@r,@a,@bb,@ab,'FUND_REQUEST',@rm);
	`, pgx.NamedArgs{
		"u":  fr.RequestToID,
		"r":  ref,
		"a":  fr.Amount,
		"bb": senderBalance,
		"ab": senderAfter,
		"rm": fmt.Sprintf("Fund sent to %s", fr.RequesterID),
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions
		(user_id, reference_id, credit_amount, before_balance, after_balance, transaction_reason, remarks)
		VALUES (@u,@r,@a,@bb,@ab,'FUND_REQUEST',@rm);
	`, pgx.NamedArgs{
		"u":  fr.RequesterID,
		"r":  ref,
		"a":  fr.Amount,
		"bb": receiverBalance,
		"ab": receiverAfter,
		"rm": fmt.Sprintf("Fund received from %s", fr.RequestToID),
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE fund_requests
		SET request_status='ACCEPTED', updated_at=NOW()
		WHERE fund_request_id=@id;
	`, pgx.NamedArgs{"id": fundRequestID})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *Database) getUserTable(userID string) (string, error) {
	if len(userID) == 0 {
		return "", errors.New("invalid user id")
	}

	switch userID[0] {
	case 'A':
		return "admin", nil
	case 'M':
		return "master_distributor", nil
	case 'D':
		return "distributor", nil
	case 'R':
		return "retailer", nil
	default:
		return "", errors.New("unknown user type")
	}
}
