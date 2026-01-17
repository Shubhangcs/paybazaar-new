package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateMasterDistributorQuery(
	ctx context.Context,
	req models.CreateMasterDistributorRequestModel,
) error {

	query := `
		INSERT INTO master_distributors (
			admin_id,
			master_distributor_name,
			master_distributor_phone,
			master_distributor_email,
			master_distributor_password,
			master_distributor_aadhar_number,
			master_distributor_pan_number,
			master_distributor_date_of_birth,
			master_distributor_gender,
			master_distributor_city,
			master_distributor_state,
			master_distributor_address,
			master_distributor_pincode,
			master_distributor_business_name,
			master_distributor_business_type,
			master_distributor_mpin,
			master_distributor_gst_number
		) VALUES (
			@admin_id,
			@name,
			@phone,
			@email,
			@password,
			@aadhar_number,
			@pan_number,
			@date_of_birth,
			@gender,
			@city,
			@state,
			@address,
			@pincode,
			@business_name,
			@business_type,
			@mpin,
			@gst_number
		)
	`

	_, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":      req.AdminID,
		"name":          req.Name,
		"phone":         req.Phone,
		"email":         req.Email,
		"password":      req.Password,
		"aadhar_number": req.AadharNumber,
		"pan_number":    req.PanNumber,
		"date_of_birth": req.DateOfBirth,
		"gender":        req.Gender,
		"city":          req.City,
		"state":         req.State,
		"address":       req.Address,
		"pincode":       req.Pincode,
		"business_name": req.BusinessName,
		"business_type": req.BusinessType,
		"mpin":          req.MPIN,
		"gst_number":    req.GSTNumber,
	})

	return err
}

func (db *Database) GetMasterDistributorByIDQuery(
	ctx context.Context,
	mdID string,
) (*models.MasterDistributorModel, error) {

	query := `
		SELECT
			m.master_distributor_id,
			m.admin_id,
			m.master_distributor_name,
			m.master_distributor_phone,
			m.master_distributor_email,
			m.master_distributor_password,
			m.master_distributor_aadhar_number,
			m.master_distributor_pan_number,
			m.master_distributor_date_of_birth,
			m.master_distributor_gender,
			m.master_distributor_city,
			m.master_distributor_state,
			m.master_distributor_address,
			m.master_distributor_pincode,
			m.master_distributor_business_name,
			m.master_distributor_business_type,
			m.master_distributor_mpin,
			m.master_distributor_kyc_status,
			m.master_distributor_documents_url,
			m.master_distributor_gst_number,
			m.master_distributor_wallet_balance,
			m.is_master_distributor_blocked,
			m.created_at,
			m.updated_at
		FROM master_distributors m
		WHERE m.master_distributor_id = @md_id;
	`

	row := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"md_id": mdID,
		},
	)

	var md models.MasterDistributorModel
	err := row.Scan(
		&md.MasterDistributorID,
		&md.AdminID,
		&md.Name,
		&md.Phone,
		&md.Email,
		&md.Password,
		&md.AadharNumber,
		&md.PanNumber,
		&md.DateOfBirth,
		&md.Gender,
		&md.City,
		&md.State,
		&md.Address,
		&md.Pincode,
		&md.BusinessName,
		&md.BusinessType,
		&md.MPIN,
		&md.KYCStatus,
		&md.DocumentsURL,
		&md.GSTNumber,
		&md.WalletBalance,
		&md.IsBlocked,
		&md.CreatedAt,
		&md.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &md, nil
}

func (db *Database) GetMasterDistributorByEmailQuery(
	ctx context.Context,
	email string,
) (*models.MasterDistributorModel, error) {

	query := `
		SELECT
			master_distributor_id,
			admin_id,
			master_distributor_name,
			master_distributor_password,
			is_master_distributor_blocked
		FROM master_distributors
		WHERE master_distributor_email = @email
	`

	row := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"email": email,
	})

	var md models.MasterDistributorModel
	err := row.Scan(
		&md.MasterDistributorID,
		&md.AdminID,
		&md.Name,
		&md.Password,
		&md.IsBlocked,
	)

	if err != nil {
		return nil, err
	}

	return &md, nil
}

func (db *Database) UpdateMasterDistributorQuery(
	ctx context.Context,
	mdID string,
	req models.UpdateMasterDistributorRequestModel,
) error {

	query := `
		UPDATE master_distributors
		SET
			master_distributor_name = COALESCE(@name, master_distributor_name),
			master_distributor_phone = COALESCE(@phone, master_distributor_phone),
			master_distributor_password = COALESCE(@password, master_distributor_password),
			master_distributor_city = COALESCE(@city, master_distributor_city),
			master_distributor_state = COALESCE(@state, master_distributor_state),
			master_distributor_address = COALESCE(@address, master_distributor_address),
			master_distributor_pincode = COALESCE(@pincode, master_distributor_pincode),
			master_distributor_business_name = COALESCE(@business_name, master_distributor_business_name),
			master_distributor_business_type = COALESCE(@business_type, master_distributor_business_type),
			master_distributor_mpin = COALESCE(@mpin, master_distributor_mpin),
			master_distributor_kyc_status = COALESCE(@kyc_status, master_distributor_kyc_status),
			master_distributor_documents_url = COALESCE(@documents_url, master_distributor_documents_url),
			master_distributor_gst_number = COALESCE(@gst_number, master_distributor_gst_number),
			master_distributor_wallet_balance = COALESCE(@wallet_balance, master_distributor_wallet_balance),
			is_master_distributor_blocked = COALESCE(@is_blocked, is_master_distributor_blocked),
			updated_at = NOW()
		WHERE master_distributor_id = @md_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id":          mdID,
		"name":           req.Name,
		"phone":          req.Phone,
		"password":       req.Password,
		"city":           req.City,
		"state":          req.State,
		"address":        req.Address,
		"pincode":        req.Pincode,
		"business_name":  req.BusinessName,
		"business_type":  req.BusinessType,
		"mpin":           req.MPIN,
		"kyc_status":     req.KYCStatus,
		"documents_url":  req.DocumentsURL,
		"gst_number":     req.GSTNumber,
		"wallet_balance": req.WalletBalance,
		"is_blocked":     req.IsBlocked,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) DeleteMasterDistributorQuery(
	ctx context.Context,
	mdID string,
) error {

	query := `
		DELETE FROM master_distributors
		WHERE master_distributor_id = @md_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id": mdID,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) ListMasterDistributorsQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetMasterDistributorResponseModel, error) {

	query := `
		SELECT
			master_distributor_id,
			admin_id,
			master_distributor_name,
			master_distributor_phone,
			master_distributor_email,
			master_distributor_aadhar_number,
			master_distributor_pan_number,
			master_distributor_date_of_birth,
			master_distributor_gender,
			master_distributor_city,
			master_distributor_state,
			master_distributor_address,
			master_distributor_pincode,
			master_distributor_business_name,
			master_distributor_business_type,
			master_distributor_kyc_status,
			master_distributor_documents_url,
			master_distributor_gst_number,
			master_distributor_wallet_balance,
			is_master_distributor_blocked,
			created_at,
			updated_at
		FROM master_distributors
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

	var list []models.GetMasterDistributorResponseModel

	for rows.Next() {
		var md models.GetMasterDistributorResponseModel
		err := rows.Scan(
			&md.MasterDistributorID,
			&md.AdminID,
			&md.Name,
			&md.Phone,
			&md.Email,
			&md.AadharNumber,
			&md.PanNumber,
			&md.DateOfBirth,
			&md.Gender,
			&md.City,
			&md.State,
			&md.Address,
			&md.Pincode,
			&md.BusinessName,
			&md.BusinessType,
			&md.KYCStatus,
			&md.DocumentsURL,
			&md.GSTNumber,
			&md.WalletBalance,
			&md.IsBlocked,
			&md.CreatedAt,
			&md.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, md)
	}

	return list, nil
}

func (db *Database) ListMasterDistributorsByAdminIDQuery(
	ctx context.Context,
	adminID string,
	limit, offset int,
) ([]models.GetMasterDistributorResponseModel, error) {

	query := `
		SELECT
			master_distributor_id,
			admin_id,
			master_distributor_name,
			master_distributor_phone,
			master_distributor_email,
			master_distributor_aadhar_number,
			master_distributor_pan_number,
			master_distributor_date_of_birth,
			master_distributor_gender,
			master_distributor_city,
			master_distributor_state,
			master_distributor_address,
			master_distributor_pincode,
			master_distributor_business_name,
			master_distributor_business_type,
			master_distributor_kyc_status,
			master_distributor_documents_url,
			master_distributor_gst_number,
			master_distributor_wallet_balance,
			is_master_distributor_blocked,
			created_at,
			updated_at
		FROM master_distributors
		WHERE admin_id = @admin_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
		"limit":    limit,
		"offset":   offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetMasterDistributorResponseModel

	for rows.Next() {
		var md models.GetMasterDistributorResponseModel
		if err := rows.Scan(
			&md.MasterDistributorID,
			&md.AdminID,
			&md.Name,
			&md.Phone,
			&md.Email,
			&md.AadharNumber,
			&md.PanNumber,
			&md.DateOfBirth,
			&md.Gender,
			&md.City,
			&md.State,
			&md.Address,
			&md.Pincode,
			&md.BusinessName,
			&md.BusinessType,
			&md.KYCStatus,
			&md.DocumentsURL,
			&md.GSTNumber,
			&md.WalletBalance,
			&md.IsBlocked,
			&md.CreatedAt,
			&md.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, md)
	}

	return list, nil
}

func (db *Database) GetMasterDistributorsByAdminIDForDropdownQuery(
	ctx context.Context,
	adminID string,
) ([]models.DropdownModel, error) {

	query := `
		SELECT
			master_distributor_id,
			master_distributor_name
		FROM master_distributors
		WHERE admin_id = @admin_id
		ORDER BY master_distributor_name ASC
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
	})
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

func (db *Database) UpdateMasterDistributorBlockStatusQuery(
	ctx context.Context,
	mdID string,
	isBlocked bool,
) error {

	query := `
		UPDATE master_distributors
		SET
			is_master_distributor_blocked = @is_blocked,
			updated_at = NOW()
		WHERE master_distributor_id = @md_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id":      mdID,
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

func (db *Database) UpdateMasterDistributorKYCStatusQuery(
	ctx context.Context,
	mdID string,
	kycStatus bool,
) error {

	query := `
		UPDATE master_distributors
		SET
			master_distributor_kyc_status = @kyc_status,
			updated_at = NOW()
		WHERE master_distributor_id = @md_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id":      mdID,
		"kyc_status": kycStatus,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
