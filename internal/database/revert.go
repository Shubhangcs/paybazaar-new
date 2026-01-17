package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateRevertQuery(
	ctx context.Context,
	req models.CreateRevertRequest,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	/* -------------------------------------------------------
	   1. Resolve tables
	------------------------------------------------------- */

	fromTable, err := db.getUserTable(req.FromID)
	if err != nil {
		return err
	}

	onTable, err := db.getUserTable(req.OnID)
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   2. Lock & get balances
	------------------------------------------------------- */

	var fromBefore, onBefore float64

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

	fromBefore, err = getBalance(fromTable, req.FromID)
	if err != nil {
		return err
	}

	onBefore, err = getBalance(onTable, req.OnID)
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   3. Balance validation
	------------------------------------------------------- */

	if onBefore < req.Amount {
		return fmt.Errorf("insufficient wallet balance")
	}

	fromAfter := fromBefore + req.Amount
	onAfter := onBefore - req.Amount

	/* -------------------------------------------------------
	   4. Create revert transaction (PENDING)
	------------------------------------------------------- */

	var revertID int64
	insertRevertQuery := `
		INSERT INTO revert_transactions (
			revert_by_id,
			revert_on_id,
			amount,
			revert_status,
			remarks
		) VALUES (
			@by_id,
			@on_id,
			@amount,
			'PENDING',
			@remarks
		)
		RETURNING revert_transaction_id;
	`

	err = tx.QueryRow(ctx, insertRevertQuery, pgx.NamedArgs{
		"by_id":   req.FromID,
		"on_id":   req.OnID,
		"amount":  req.Amount,
		"remarks": req.Remarks,
	}).Scan(&revertID)

	if err != nil {
		return err
	}

	refID := fmt.Sprintf("%d", revertID)

	/* -------------------------------------------------------
	   5. Update wallets
	------------------------------------------------------- */

	updateWallet := func(table, id string, balance float64) error {
		query := fmt.Sprintf(`
			UPDATE %ss
			SET %s_wallet_balance = @bal,
			    updated_at = NOW()
			WHERE %s_id = @id
		`, table, table, table)

		_, err := tx.Exec(ctx, query, pgx.NamedArgs{
			"id":  id,
			"bal": balance,
		})
		return err
	}

	if err := updateWallet(fromTable, req.FromID, fromAfter); err != nil {
		return err
	}

	if err := updateWallet(onTable, req.OnID, onAfter); err != nil {
		return err
	}

	/* -------------------------------------------------------
	   6. Wallet transactions
	------------------------------------------------------- */

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
			'REVERT',
			@remarks
		);
	`

	// Credit → FROM user
	_, err = tx.Exec(ctx, insertTxn, pgx.NamedArgs{
		"uid":    req.FromID,
		"ref":    refID,
		"credit": req.Amount,
		"debit":  nil,
		"before": fromBefore,
		"after":  fromAfter,
		"remarks": fmt.Sprintf(
			"Revert received from %s", req.OnID,
		),
	})
	if err != nil {
		return err
	}

	// Debit → ON user
	_, err = tx.Exec(ctx, insertTxn, pgx.NamedArgs{
		"uid":    req.OnID,
		"ref":    refID,
		"credit": nil,
		"debit":  req.Amount,
		"before": onBefore,
		"after":  onAfter,
		"remarks": fmt.Sprintf(
			"Revert sent to %s", req.FromID,
		),
	})
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   7. Mark revert SUCCESS
	------------------------------------------------------- */

	_, err = tx.Exec(ctx, `
		UPDATE revert_transactions
		SET revert_status = 'SUCCESS'
		WHERE revert_transaction_id = @id
	`, pgx.NamedArgs{
		"id": revertID,
	})
	if err != nil {
		return err
	}

	/* -------------------------------------------------------
	   8. Commit
	------------------------------------------------------- */

	return tx.Commit(ctx)
}
