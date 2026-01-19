package database

import (
	"context"
	"fmt"

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
	query := `
		UPDATE admins 
		SET admin_wallet_balance = admin_wallet_balance + @amount,
		updated_at = NOW()
		WHERE admin_id = @admin_id;
	`
	tag, err := db.pool.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"admin_id": req.AdminID,
			"amount":   req.Amount,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update admin wallet")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid admin id or admin not found")
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
