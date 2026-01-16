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

	_, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_name":     req.AdminName,
		"admin_email":    req.AdminEmail,
		"admin_phone":    req.AdminPhone,
		"admin_password": req.AdminPassword,
	})

	return err
}

func (db *Database) GetAdminByIDQuery(
	ctx context.Context,
	adminID string,
) (*models.AdminModel, error) {

	query := `
		SELECT
			admin_id,
			admin_name,
			admin_email,
			admin_phone,
			admin_password,
			is_admin_blocked
		FROM admins
		WHERE admin_id = @admin_id
	`

	row := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
	})

	var admin models.AdminModel
	err := row.Scan(
		&admin.AdminID,
		&admin.AdminName,
		&admin.AdminEmail,
		&admin.AdminPhone,
		&admin.AdminPassword,
		&admin.IsAdminBlocked,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (db *Database) GetAdminByEmailQuery(
	ctx context.Context,
	email string,
) (*models.AdminModel, error) {

	query := `
		SELECT
			admin_id,
			admin_name,
			admin_email,
			admin_phone,
			admin_password,
			is_admin_blocked
		FROM admins
		WHERE admin_email = @admin_email
	`

	row := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"admin_email": email,
	})

	var admin models.AdminModel
	err := row.Scan(
		&admin.AdminID,
		&admin.AdminName,
		&admin.AdminEmail,
		&admin.AdminPhone,
		&admin.AdminPassword,
		&admin.IsAdminBlocked,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
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
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) UpdateAdminQuery(
	ctx context.Context,
	adminID string,
	req models.UpdateAdminRequestModel,
) error {

	query := `
		UPDATE admins
		SET
			admin_name = COALESCE(@admin_name, admin_name),
			admin_phone = COALESCE(@admin_phone, admin_phone),
			admin_password = COALESCE(@admin_password, admin_password),
			admin_wallet_balance = COALESCE(@wallet_balance, admin_wallet_balance),
			is_admin_blocked = COALESCE(@is_admin_blocked, is_admin_blocked),
			updated_at = NOW()
		WHERE admin_id = @admin_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":         adminID,
		"admin_name":       req.AdminName,
		"admin_phone":      req.AdminPhone,
		"admin_password":   req.AdminPassword,
		"wallet_balance":   req.WalletBalance,
		"is_admin_blocked": req.IsAdminBlocked,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) ListAdminsQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetAdminResponseModel, error) {

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
		return nil, err
	}
	defer rows.Close()

	var list []models.GetAdminResponseModel

	for rows.Next() {
		var admin models.GetAdminResponseModel
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
			return nil, err
		}
		fmt.Println(admin)
		list = append(list, admin)
	}
	fmt.Println("hello")

	return list, nil
}

func (db *Database) GetAllAdminsForDropdownQuery(
	ctx context.Context,
) ([]models.DropdownModel, error) {

	query := `
		SELECT
			admin_id,
			admin_name
		FROM admins
		ORDER BY admin_name ASC
	`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DropdownModel

	for rows.Next() {
		var d models.DropdownModel
		if err := rows.Scan(
			&d.ID,
			&d.Name,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}

func (db *Database) UpdateAdminBlockStatusQuery(
	ctx context.Context,
	adminID string,
	isBlocked bool,
) error {

	query := `
		UPDATE admins
		SET
			is_admin_blocked = @is_blocked,
			updated_at = NOW()
		WHERE admin_id = @admin_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":   adminID,
		"is_blocked": isBlocked,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
