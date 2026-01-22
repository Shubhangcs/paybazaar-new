package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateAdminQuery(
	ctx context.Context,
	req models.CreateAdminRequestModel,
) error {
	query := `
		INSERT INTO admins (
			admin_name,
			admin_email,
			admin_phone,
			admin_password
		) VALUES (
			@admin_name,
			@admin_email,
			@admin_phone,
			@admin_password
		)
	`
	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_name":     req.AdminName,
		"admin_email":    req.AdminEmail,
		"admin_phone":    req.AdminPhone,
		"admin_password": req.AdminPassword,
	}); err != nil {
		return fmt.Errorf("failed to create admin")
	}
	return nil
}

func (db *Database) GetAdminDetailsByAdminID(
	ctx context.Context,
	adminID string,
) (*models.GetCompleteAdminDetailsResponseModel, error) {
	query := `
		SELECT
			admin_id,
			admin_name,
			admin_email,
			admin_phone,
			admin_wallet_balance,
			is_admin_blocked,
			created_at,
			updated_at
		FROM admins
		WHERE admin_id = @admin_id;
	`
	var res models.GetCompleteAdminDetailsResponseModel
	if err := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"admin_id": adminID,
		},
	).Scan(
		&res.AdminID,
		&res.AdminName,
		&res.AdminEmail,
		&res.AdminPhone,
		&res.AdminWalletBalance,
		&res.IsAdminBlocked,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch admin details")
	}
	return &res, nil
}

func (db *Database) GetAdminDetailsForLogin(
	ctx context.Context,
	adminID string,
) (*models.GetAdminDetailsForLoginModel, error) {
	query := `
		SELECT
			admin_id,
			admin_name,
			admin_password,
			is_admin_blocked
		FROM admins
		WHERE admin_id = @admin_id;
	`
	var res models.GetAdminDetailsForLoginModel
	if err := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"admin_id": adminID,
		},
	).Scan(
		&res.AdminID,
		&res.AdminName,
		&res.AdminPassword,
		&res.IsAdminBlocked,
	); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to fetch admin details")
	}
	return &res, nil
}

func (db *Database) DeleteAdminQuery(
	ctx context.Context,
	adminID string,
) error {

	query := `
		DELETE FROM admins
		WHERE admin_id = @admin_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete admin")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid admin id or admin not found")
	}

	return nil
}

func (db *Database) UpdateAdminDetailsQuery(
	ctx context.Context,
	req models.UpdateAdminDetailsRequestModel,
) error {
	query := `
		UPDATE admins
		SET
			admin_name  = COALESCE(@admin_name, admin_name),
			admin_phone = COALESCE(@admin_phone, admin_phone),
			admin_email = COALESCE(@admin_email, admin_email),
			updated_at  = NOW()
		WHERE admin_id = @admin_id
	`
	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":    req.AdminID,
		"admin_name":  req.AdminName,
		"admin_phone": req.AdminPhone,
		"admin_email": req.AdminEmail,
	})

	if err != nil {
		return fmt.Errorf("failed to update admin")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid admin id or admin not found")
	}
	return nil
}

func (db *Database) UpdateAdminPasswordQuery(
	ctx context.Context,
	req models.UpdateAdminPasswordRequestModel,
) error {

	getAdminOldPasswordQuery := `
		SELECT admin_password FROM admins WHERE admin_id = @admin_id;
	`

	var oldPassword string
	if err := db.pool.QueryRow(
		ctx,
		getAdminOldPasswordQuery,
		pgx.NamedArgs{
			"admin_id": req.AdminID,
		},
	).Scan(&oldPassword); err != nil {
		return fmt.Errorf("failed to fetch old password")
	}

	if oldPassword != req.OldPassword {
		return fmt.Errorf("incorrect old password")
	}

	updateAdminPasswordQuery := `
		UPDATE admins
		SET admin_password = @new_admin_password,
		updated_at = NOW()
		WHERE admin_id = @admin_id;
	`

	if _, err := db.pool.Exec(
		ctx,
		updateAdminPasswordQuery,
		pgx.NamedArgs{
			"admin_id":           req.AdminID,
			"new_admin_password": req.NewPassword,
		},
	); err != nil {
		return fmt.Errorf("failed to update admin password")
	}
	return nil
}

func (db *Database) UpdateAdminWalletQuery(
	ctx context.Context,
	req models.UpdateAdminWalletRequestModel,
) error {

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction")
	}
	defer tx.Rollback(ctx)

	// 1️⃣ Get current wallet balance
	var beforeBalance float64
	getBalanceQuery := `
		SELECT admin_wallet_balance
		FROM admins
		WHERE admin_id = @admin_id;
	`

	if err := tx.QueryRow(
		ctx,
		getBalanceQuery,
		pgx.NamedArgs{
			"admin_id": req.AdminID,
		},
	).Scan(&beforeBalance); err != nil {
		return fmt.Errorf("failed to fetch admin wallet balance")
	}

	afterBalance := beforeBalance + req.Amount

	// 2️⃣ Update admin wallet
	updateWalletQuery := `
		UPDATE admins
		SET admin_wallet_balance = @after_balance,
		    updated_at = NOW()
		WHERE admin_id = @admin_id;
	`

	tag, err := tx.Exec(
		ctx,
		updateWalletQuery,
		pgx.NamedArgs{
			"admin_id":      req.AdminID,
			"after_balance": afterBalance,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update admin wallet")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid admin id or admin not found")
	}

	// 3️⃣ Insert wallet transaction record
	insertTransactionQuery := `
		INSERT INTO wallet_transactions (
			user_id,
			reference_id,
			credit_amount,
			before_balance,
			after_balance,
			transaction_reason,
			remarks
		) VALUES (
			@user_id,
			@reference_id,
			@credit_amount,
			@before_balance,
			@after_balance,
			@transaction_reason,
			@remarks
		);
	`

	if _, err := tx.Exec(
		ctx,
		insertTransactionQuery,
		pgx.NamedArgs{
			"user_id":            req.AdminID,
			"reference_id":       req.AdminID, // ✅ admin id as reference
			"credit_amount":      req.Amount,
			"before_balance":     beforeBalance,
			"after_balance":      afterBalance,
			"transaction_reason": "TOPUP",
			"remarks":            "Admin wallet topup",
		},
	); err != nil {
		return fmt.Errorf("failed to insert wallet transaction")
	}

	// 4️⃣ Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit wallet transaction")
	}

	return nil
}

func (db *Database) UpdateAdminBlockStatusQuery(
	ctx context.Context,
	req models.UpdateAdminBlockStatusRequestModel,
) error {
	query := `
		UPDATE admins
		SET is_admin_blocked = @admin_block_status,
		updated_at = NOW()
		WHERE admin_id = @admin_id;
	`

	tag, err := db.pool.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"admin_id":           req.AdminID,
			"admin_block_status": req.BlockStatus,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update admin block status")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid admin id or admin not found")
	}
	return nil
}

func (db *Database) GetAllAdminsQuery(
	ctx context.Context,
	offset, limit int,
) ([]models.GetCompleteAdminDetailsResponseModel, error) {
	query := `
		SELECT
			admin_id,
			admin_name,
			admin_email,
			admin_phone,
			admin_wallet_balance,
			is_admin_blocked,
			created_at,
			updated_at
		FROM admins
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch admins")
	}
	defer rows.Close()

	var list []models.GetCompleteAdminDetailsResponseModel

	for rows.Next() {
		var admin models.GetCompleteAdminDetailsResponseModel
		err := rows.Scan(
			&admin.AdminID,
			&admin.AdminName,
			&admin.AdminEmail,
			&admin.AdminPhone,
			&admin.AdminWalletBalance,
			&admin.IsAdminBlocked,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch admin details")
		}
		list = append(list, admin)
	}

	return list, rows.Err()
}

func (db *Database) GetAllAdminsForDropdownQuery(
	ctx context.Context,
) ([]models.GetAdminDetailsForDropdownModel, error) {

	query := `
		SELECT
			admin_id,
			admin_name
		FROM admins
		ORDER BY admin_name ASC
	`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch admin details")
	}
	defer rows.Close()

	var list []models.GetAdminDetailsForDropdownModel

	for rows.Next() {
		var d models.GetAdminDetailsForDropdownModel
		if err := rows.Scan(
			&d.AdminID,
			&d.AdminName,
		); err != nil {
			return nil, fmt.Errorf("failed to fetch admin details")
		}
		list = append(list, d)
	}

	return list, nil
}
