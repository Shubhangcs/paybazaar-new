package database

import (
	"context"
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
			remarks
		) VALUES (
			@requester_id,
			@request_to_id,
			@amount,
			@bank_name,
			@request_date,
			@utr_number,
			'PENDING',
			@remarks
		)
		RETURNING fund_request_id
	`

	var fundRequestID int64
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"requester_id":  req.RequesterID,
		"request_to_id": req.RequestToID,
		"amount":        req.Amount,
		"bank_name":     req.BankName,
		"request_date":  req.RequestDate,
		"utr_number":    req.UTRNumber,
		"remarks":       req.Remarks,
	}).Scan(&fundRequestID)

	if err != nil {
		return 0, err
	}

	return fundRequestID, nil
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
			created_at,
			updated_at
		FROM fund_requests
		WHERE fund_request_id = @fund_request_id
	`

	var fundRequest models.GetFundRequestResponseModel
	err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"fund_request_id": fundRequestID,
	}).Scan(
		&fundRequest.FundRequestID,
		&fundRequest.RequesterID,
		&fundRequest.RequestToID,
		&fundRequest.Amount,
		&fundRequest.BankName,
		&fundRequest.RequestDate,
		&fundRequest.UTRNumber,
		&fundRequest.RequestStatus,
		&fundRequest.Remarks,
		&fundRequest.CreatedAt,
		&fundRequest.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &fundRequest, nil
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
			created_at,
			updated_at
		FROM fund_requests
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fundRequests []models.GetFundRequestResponseModel

	for rows.Next() {
		var fundRequest models.GetFundRequestResponseModel
		err := rows.Scan(
			&fundRequest.FundRequestID,
			&fundRequest.RequesterID,
			&fundRequest.RequestToID,
			&fundRequest.Amount,
			&fundRequest.BankName,
			&fundRequest.RequestDate,
			&fundRequest.UTRNumber,
			&fundRequest.RequestStatus,
			&fundRequest.Remarks,
			&fundRequest.CreatedAt,
			&fundRequest.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fundRequest)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fundRequests, nil
}

func (db *Database) GetFundRequestsByRequesterIDQuery(
	ctx context.Context,
	req models.GetFundRequestFilterRequestModel,
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
			created_at,
			updated_at
		FROM fund_requests
		WHERE requester_id = @requester_id
	`

	args := pgx.NamedArgs{
		"requester_id": req.ID,
		"limit":        limit,
		"offset":       offset,
	}

	// 🔹 Optional created_at filters
	if req.StartDate != nil {
		query += ` AND created_at >= @start_date`
		args["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		query += ` AND created_at <= @end_date`
		args["end_date"] = *req.EndDate
	}

	// 🔹 Optional status filter
	if req.Status != nil {
		query += ` AND request_status = @status`
		args["status"] = *req.Status
	}

	query += `
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund requests by requester id: %w", err)
	}
	defer rows.Close()

	var fundRequests []models.GetFundRequestResponseModel

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
			&fr.CreatedAt,
			&fr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fundRequests, nil
}

func (db *Database) GetFundRequestsByRequestToIDQuery(
	ctx context.Context,
	req models.GetFundRequestFilterRequestModel,
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
			created_at,
			updated_at
		FROM fund_requests
		WHERE request_to_id = @request_to_id
	`

	args := pgx.NamedArgs{
		"request_to_id": req.ID,
		"limit":         limit,
		"offset":        offset,
	}

	// 🔹 Optional created_at filters
	if req.StartDate != nil {
		query += ` AND created_at >= @start_date`
		args["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		query += ` AND created_at <= @end_date`
		args["end_date"] = *req.EndDate
	}

	// 🔹 Optional status filter
	if req.Status != nil {
		query += ` AND request_status = @status`
		args["status"] = *req.Status
	}

	query += `
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund requests by request_to id: %w", err)
	}
	defer rows.Close()

	var fundRequests []models.GetFundRequestResponseModel

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
			&fr.CreatedAt,
			&fr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fundRequests, nil
}

func (db *Database) RejectFundRequestQuery(
	ctx context.Context,
	fundRequestID int64,
) error {

	query := `
		UPDATE fund_requests
		SET
			request_status = 'REJECTED',
			updated_at = NOW()
		WHERE fund_request_id = @fund_request_id
		AND request_status = 'PENDING'
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"fund_request_id": fundRequestID,
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

	// Start transaction
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get fund request details
	var fundRequest struct {
		RequesterID   string
		RequestToID   string
		Amount        string
		RequestStatus string
	}

	getFundRequestQuery := `
		SELECT requester_id, request_to_id, amount, request_status
		FROM fund_requests
		WHERE fund_request_id = @fund_request_id
	`

	err = tx.QueryRow(ctx, getFundRequestQuery, pgx.NamedArgs{
		"fund_request_id": fundRequestID,
	}).Scan(
		&fundRequest.RequesterID,
		&fundRequest.RequestToID,
		&fundRequest.Amount,
		&fundRequest.RequestStatus,
	)

	if err != nil {
		return err
	}

	// Check if already processed
	if fundRequest.RequestStatus != "PENDING" {
		return pgx.ErrNoRows // or custom error
	}

	// Convert amount to numeric for calculations
	var amountNumeric float64
	_, err = fmt.Sscanf(fundRequest.Amount, "%f", &amountNumeric)
	if err != nil {
		return err
	}

	// Get current balance of request_to user (sender)
	var senderBeforeBalance float64
	var senderTable string

	senderTable, err = db.getUserTable(fundRequest.RequestToID)
	if err != nil {
		return err
	}

	getSenderBalanceQuery := fmt.Sprintf(`
		SELECT %s_wallet_balance
		FROM %s
		WHERE %s_id = @user_id
	`, senderTable, senderTable+"s", senderTable)

	err = tx.QueryRow(ctx, getSenderBalanceQuery, pgx.NamedArgs{
		"user_id": fundRequest.RequestToID,
	}).Scan(&senderBeforeBalance)

	if err != nil {
		return err
	}

	// Check sufficient balance
	if senderBeforeBalance < amountNumeric {
		return fmt.Errorf("insufficient balance")
	}

	senderAfterBalance := senderBeforeBalance - amountNumeric

	// Get current balance of requester (receiver)
	var receiverBeforeBalance float64
	var receiverTable string

	receiverTable, err = db.getUserTable(fundRequest.RequesterID)
	if err != nil {
		return err
	}

	getReceiverBalanceQuery := fmt.Sprintf(`
		SELECT %s_wallet_balance
		FROM %s
		WHERE %s_id = @user_id
	`, receiverTable, receiverTable+"s", receiverTable)

	err = tx.QueryRow(ctx, getReceiverBalanceQuery, pgx.NamedArgs{
		"user_id": fundRequest.RequesterID,
	}).Scan(&receiverBeforeBalance)

	if err != nil {
		return err
	}

	receiverAfterBalance := receiverBeforeBalance + amountNumeric

	// Update sender wallet (deduct amount)
	updateSenderQuery := fmt.Sprintf(`
		UPDATE %s
		SET %s_wallet_balance = @new_balance,
		    updated_at = NOW()
		WHERE %s_id = @user_id
	`, senderTable+"s", senderTable, senderTable)

	_, err = tx.Exec(ctx, updateSenderQuery, pgx.NamedArgs{
		"user_id":     fundRequest.RequestToID,
		"new_balance": senderAfterBalance,
	})

	if err != nil {
		return err
	}

	// Update receiver wallet (add amount)
	updateReceiverQuery := fmt.Sprintf(`
		UPDATE %s
		SET %s_wallet_balance = @new_balance,
		    updated_at = NOW()
		WHERE %s_id = @user_id
	`, receiverTable+"s", receiverTable, receiverTable)

	_, err = tx.Exec(ctx, updateReceiverQuery, pgx.NamedArgs{
		"user_id":     fundRequest.RequesterID,
		"new_balance": receiverAfterBalance,
	})

	if err != nil {
		return err
	}

	// Insert debit transaction for sender
	insertDebitTransactionQuery := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			credit_amount,
			debit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			NULL,
			@debit_amount,
			@before_balance,
			@after_balance,
			'FUND_REQUEST',
			@remarks
		)
	`

	_, err = tx.Exec(ctx, insertDebitTransactionQuery, pgx.NamedArgs{
		"user_id":        fundRequest.RequestToID,
		"reference_id":   fmt.Sprintf("FR%d", fundRequestID),
		"debit_amount":   fundRequest.Amount,
		"before_balance": fmt.Sprintf("%.2f", senderBeforeBalance),
		"after_balance":  fmt.Sprintf("%.2f", senderAfterBalance),
		"remarks":        fmt.Sprintf("Fund request to %s", fundRequest.RequesterID),
	})

	if err != nil {
		return err
	}

	// Insert credit transaction for receiver
	insertCreditTransactionQuery := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			credit_amount,
			debit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@credit_amount,
			NULL,
			@before_balance,
			@after_balance,
			'FUND_REQUEST',
			@remarks
		)
	`

	_, err = tx.Exec(ctx, insertCreditTransactionQuery, pgx.NamedArgs{
		"user_id":        fundRequest.RequesterID,
		"reference_id":   fmt.Sprintf("FR%d", fundRequestID),
		"credit_amount":  fundRequest.Amount,
		"before_balance": fmt.Sprintf("%.2f", receiverBeforeBalance),
		"after_balance":  fmt.Sprintf("%.2f", receiverAfterBalance),
		"remarks":        fmt.Sprintf("Fund request from %s", fundRequest.RequestToID),
	})

	if err != nil {
		return err
	}

	// Update fund request status to ACCEPTED
	updateFundRequestQuery := `
		UPDATE fund_requests
		SET request_status = 'ACCEPTED',
		    updated_at = NOW()
		WHERE fund_request_id = @fund_request_id
	`

	_, err = tx.Exec(ctx, updateFundRequestQuery, pgx.NamedArgs{
		"fund_request_id": fundRequestID,
	})

	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(ctx)
}

func (db *Database) getUserTable(userID string) (string, error) {
	if len(userID) == 0 {
		return "", fmt.Errorf("invalid user ID")
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
		return "", fmt.Errorf("unknown user type")
	}
}
