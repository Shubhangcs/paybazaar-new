package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateFundTransferQuery(
	ctx context.Context,
	req models.CreateFundTransferModel,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	fromTable, err := db.getUserTable(req.FromID)
	if err != nil {
		return err
	}

	toTable, err := db.getUserTable(req.ToID)
	if err != nil {
		return err
	}

	getBalance := func(table, id string) (float64, error) {
		query := fmt.Sprintf(`
			SELECT %s_wallet_balance
			FROM %ss
			WHERE %s_id = @id
			FOR UPDATE
		`, table, table, table)

		var bal float64
		err := tx.QueryRow(ctx, query, pgx.NamedArgs{
			"id": id,
		}).Scan(&bal)

		return bal, err
	}

	fromBefore, err := getBalance(fromTable, req.FromID)
	if err != nil {
		return err
	}

	toBefore, err := getBalance(toTable, req.ToID)
	if err != nil {
		return err
	}

	if fromBefore < req.Amount {
		return fmt.Errorf("insufficient wallet balance")
	}

	fromAfter := fromBefore - req.Amount
	toAfter := toBefore + req.Amount

	var transferID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO fund_transfers (
			fund_transfer_by_id,
			fund_treanfer_to_id,
			amount,
			fund_transfer_status,
			remarks
		) VALUES (
			@from_id,
			@to_id,
			@amount,
			'PENDING',
			@remarks
		)
		RETURNING fund_transfer_id
	`, pgx.NamedArgs{
		"from_id": req.FromID,
		"to_id":   req.ToID,
		"amount":  req.Amount,
		"remarks": req.Remarks,
	}).Scan(&transferID)

	if err != nil {
		return err
	}

	refID := fmt.Sprintf("%d", transferID)

	updateWallet := func(table, id string, bal float64) error {
		query := fmt.Sprintf(`
			UPDATE %ss
			SET %s_wallet_balance = @bal,
			    updated_at = NOW()
			WHERE %s_id = @id
		`, table, table, table)

		_, err := tx.Exec(ctx, query, pgx.NamedArgs{
			"id":  id,
			"bal": bal,
		})
		return err
	}

	if err := updateWallet(fromTable, req.FromID, fromAfter); err != nil {
		return err
	}

	if err := updateWallet(toTable, req.ToID, toAfter); err != nil {
		return err
	}

	insertTxn := `
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
			@uid,
			@ref,
			@credit,
			@debit,
			@before,
			@after,
			'FUND_TRANSFER',
			@remarks
		)
	`

	// Debit FROM user
	_, err = tx.Exec(ctx, insertTxn, pgx.NamedArgs{
		"uid":    req.FromID,
		"ref":    refID,
		"credit": nil,
		"debit":  req.Amount,
		"before": fromBefore,
		"after":  fromAfter,
		"remarks": fmt.Sprintf(
			"Fund transfer to %s", req.ToID,
		),
	})
	if err != nil {
		return err
	}

	// Credit TO user
	_, err = tx.Exec(ctx, insertTxn, pgx.NamedArgs{
		"uid":    req.ToID,
		"ref":    refID,
		"credit": req.Amount,
		"debit":  nil,
		"before": toBefore,
		"after":  toAfter,
		"remarks": fmt.Sprintf(
			"Fund transfer from %s", req.FromID,
		),
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE fund_transfers
		SET fund_transfer_status = 'SUCCESS'
		WHERE fund_transfer_id = @id
	`, pgx.NamedArgs{
		"id": transferID,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *Database) GetFundTransfersByFromIDQuery(
	ctx context.Context,
	req models.GetFundTransferFilterRequestModel,
	limit, offset int,
) ([]models.GetFundTransferResponseModel, error) {

	query := `
		SELECT
			ft.fund_transfer_id,
			ft.fund_transfer_by_id,
			ft.fund_treanfer_to_id,

			COALESCE(
				a.admin_name,
				md.master_distributor_name,
				d.distributor_name,
				r.retailer_name
			) AS from_name,

			COALESCE(
				a2.admin_name,
				md2.master_distributor_name,
				d2.distributor_name,
				r2.retailer_name
			) AS to_name,

			ft.amount,
			ft.remarks,
			ft.created_at
		FROM fund_transfers ft

		LEFT JOIN admins a ON ft.fund_transfer_by_id = a.admin_id
		LEFT JOIN master_distributors md ON ft.fund_transfer_by_id = md.master_distributor_id
		LEFT JOIN distributors d ON ft.fund_transfer_by_id = d.distributor_id
		LEFT JOIN retailers r ON ft.fund_transfer_by_id = r.retailer_id

		LEFT JOIN admins a2 ON ft.fund_treanfer_to_id = a2.admin_id
		LEFT JOIN master_distributors md2 ON ft.fund_treanfer_to_id = md2.master_distributor_id
		LEFT JOIN distributors d2 ON ft.fund_treanfer_to_id = d2.distributor_id
		LEFT JOIN retailers r2 ON ft.fund_treanfer_to_id = r2.retailer_id

		WHERE ft.fund_transfer_by_id = @id
	`

	args := pgx.NamedArgs{
		"id":     req.ID,
		"limit":  limit,
		"offset": offset,
	}

	if req.StartDate != nil {
		query += ` AND ft.created_at >= @start_date`
		args["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		query += ` AND ft.created_at <= @end_date`
		args["end_date"] = *req.EndDate
	}

	if req.Status != nil {
		query += ` AND ft.fund_transfer_status = @status`
		args["status"] = *req.Status
	}

	query += `
		ORDER BY ft.created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.GetFundTransferResponseModel
	for rows.Next() {
		var r models.GetFundTransferResponseModel
		if err := rows.Scan(
			&r.FundTransferID,
			&r.FundTransferFromID,
			&r.FundTransferToID,
			&r.FundTransferFromName,
			&r.FundTransferToName,
			&r.Amount,
			&r.Remarks,
			&r.CreatedAT,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}

func (db *Database) GetFundTransfersByToIDQuery(
	ctx context.Context,
	req models.GetFundTransferFilterRequestModel,
	limit, offset int,
) ([]models.GetFundTransferResponseModel, error) {

	query := `
		SELECT
			ft.fund_transfer_id,
			ft.fund_transfer_by_id,
			ft.fund_treanfer_to_id,

			COALESCE(
				a.admin_name,
				md.master_distributor_name,
				d.distributor_name,
				r.retailer_name
			) AS from_name,

			COALESCE(
				a2.admin_name,
				md2.master_distributor_name,
				d2.distributor_name,
				r2.retailer_name
			) AS to_name,

			ft.amount,
			ft.remarks,
			ft.created_at
		FROM fund_transfers ft

		LEFT JOIN admins a ON ft.fund_transfer_by_id = a.admin_id
		LEFT JOIN master_distributors md ON ft.fund_transfer_by_id = md.master_distributor_id
		LEFT JOIN distributors d ON ft.fund_transfer_by_id = d.distributor_id
		LEFT JOIN retailers r ON ft.fund_transfer_by_id = r.retailer_id

		LEFT JOIN admins a2 ON ft.fund_treanfer_to_id = a2.admin_id
		LEFT JOIN master_distributors md2 ON ft.fund_treanfer_to_id = md2.master_distributor_id
		LEFT JOIN distributors d2 ON ft.fund_treanfer_to_id = d2.distributor_id
		LEFT JOIN retailers r2 ON ft.fund_treanfer_to_id = r2.retailer_id

		WHERE ft.fund_treanfer_to_id = @id
	`

	args := pgx.NamedArgs{
		"id":     req.ID,
		"limit":  limit,
		"offset": offset,
	}

	if req.StartDate != nil {
		query += ` AND ft.created_at >= @start_date`
		args["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		query += ` AND ft.created_at <= @end_date`
		args["end_date"] = *req.EndDate
	}

	if req.Status != nil {
		query += ` AND ft.fund_transfer_status = @status`
		args["status"] = *req.Status
	}

	query += `
		ORDER BY ft.created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.GetFundTransferResponseModel
	for rows.Next() {
		var r models.GetFundTransferResponseModel
		if err := rows.Scan(
			&r.FundTransferID,
			&r.FundTransferFromID,
			&r.FundTransferToID,
			&r.FundTransferFromName,
			&r.FundTransferToName,
			&r.Amount,
			&r.Remarks,
			&r.CreatedAT,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, rows.Err()
}

