package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateRetailerQuery(
	ctx context.Context,
	req models.CreateRetailerRequestModel,
) error {

	query := `
		INSERT INTO retailers (
			distributor_id,
			retailer_name,
			retailer_phone,
			retailer_email,
			retailer_password,
			retailer_aadhar_number,
			retailer_pan_number,
			retailer_date_of_birth,
			retailer_gender,
			retailer_city,
			retailer_state,
			retailer_address,
			retailer_pincode,
			retailer_business_name,
			retailer_business_type,
			retailer_gst_number,
			retailer_mpin
		) VALUES (
			@distributor_id,
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
			@gst_number,
			@mpin
		)
	`

	_, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": req.DistributorID,
		"name":           req.Name,
		"phone":          req.Phone,
		"email":          req.Email,
		"password":       req.Password,
		"aadhar_number":  req.AadharNumber,
		"pan_number":     req.PanNumber,
		"date_of_birth":  req.DateOfBirth,
		"gender":         req.Gender,
		"city":           req.City,
		"state":          req.State,
		"address":        req.Address,
		"pincode":        req.Pincode,
		"business_name":  req.BusinessName,
		"business_type":  req.BusinessType,
		"gst_number":     req.GSTNumber,
		"mpin":           req.MPIN,
	})

	return err
}

func (db *Database) GetRetailerByIDQuery(
	ctx context.Context,
	retailerID string,
) (*models.RetailerModel, error) {

	query := `
		SELECT
			r.retailer_id,
			r.distributor_id,
			r.retailer_name,
			r.retailer_phone,
			r.retailer_email,
			r.retailer_password,
			r.retailer_aadhar_number,
			r.retailer_pan_number,
			r.retailer_date_of_birth,
			r.retailer_gender,
			r.retailer_city,
			r.retailer_state,
			r.retailer_address,
			r.retailer_pincode,
			r.retailer_business_name,
			r.retailer_business_type,
			r.retailer_gst_number,
			r.retailer_mpin,
			r.retailer_kyc_status,
			r.retailer_documents_url,
			r.retailer_wallet_balance,
			r.is_retailer_blocked,
			r.created_at,
			r.updated_at,
			md.admin_id
		FROM retailers r
		JOIN distributors d
			ON r.distributor_id = d.distributor_id
		JOIN master_distributors md
			ON d.master_distributor_id = md.master_distributor_id
		WHERE r.retailer_id = @retailer_id;
	`

	row := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"retailer_id": retailerID,
		},
	)

	var r models.RetailerModel
	err := row.Scan(
		&r.RetailerID,
		&r.DistributorID,
		&r.Name,
		&r.Phone,
		&r.Email,
		&r.Password,
		&r.AadharNumber,
		&r.PanNumber,
		&r.DateOfBirth,
		&r.Gender,
		&r.City,
		&r.State,
		&r.Address,
		&r.Pincode,
		&r.BusinessName,
		&r.BusinessType,
		&r.GSTNumber,
		&r.MPIN,
		&r.KYCStatus,
		&r.DocumentsURL,
		&r.WalletBalance,
		&r.IsBlocked,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.AdminID,
	)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (db *Database) GetRetailerByEmailQuery(
	ctx context.Context,
	email string,
) (*models.RetailerModel, error) {

	query := `
		SELECT
			r.retailer_id,
			r.distributor_id,
			r.retailer_name,
			r.retailer_password,
			r.is_retailer_blocked,
			md.admin_id
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
		WHERE r.retailer_email = @email
	`

	row := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"email": email,
	})

	var r models.RetailerModel
	err := row.Scan(
		&r.RetailerID,
		&r.DistributorID,
		&r.Name,
		&r.Password,
		&r.IsBlocked,
		&r.AdminID,
	)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (db *Database) GetRetailerByPhoneQuery(
	ctx context.Context,
	phone string,
) (*models.RetailerModel, error) {

	query := `
		SELECT
			r.retailer_id,
			r.distributor_id,
			r.retailer_name,
			r.retailer_password,
			r.is_retailer_blocked,
			md.admin_id
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
		WHERE r.retailer_phone = @phone;
	`

	row := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"phone": phone,
		},
	)

	var r models.RetailerModel
	err := row.Scan(
		&r.RetailerID,
		&r.DistributorID,
		&r.Name,
		&r.Password,
		&r.IsBlocked,
		&r.AdminID,
	)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (db *Database) UpdateRetailerQuery(
	ctx context.Context,
	retailerID string,
	req models.UpdateRetailerRequestModel,
) error {

	query := `
		UPDATE retailers
		SET
			retailer_name = COALESCE(@name, retailer_name),
			retailer_phone = COALESCE(@phone, retailer_phone),
			retailer_password = COALESCE(@password, retailer_password),
			retailer_city = COALESCE(@city, retailer_city),
			retailer_state = COALESCE(@state, retailer_state),
			retailer_address = COALESCE(@address, retailer_address),
			retailer_pincode = COALESCE(@pincode, retailer_pincode),
			retailer_business_name = COALESCE(@business_name, retailer_business_name),
			retailer_business_type = COALESCE(@business_type, retailer_business_type),
			retailer_gst_number = COALESCE(@gst_number, retailer_gst_number),
			retailer_mpin = COALESCE(@mpin, retailer_mpin),
			retailer_kyc_status = COALESCE(@kyc_status, retailer_kyc_status),
			retailer_documents_url = COALESCE(@documents_url, retailer_documents_url),
			retailer_wallet_balance = COALESCE(@wallet_balance, retailer_wallet_balance),
			is_retailer_blocked = COALESCE(@is_blocked, is_retailer_blocked),
			updated_at = NOW()
		WHERE retailer_id = @retailer_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id":    retailerID,
		"name":           req.Name,
		"phone":          req.Phone,
		"password":       req.Password,
		"city":           req.City,
		"state":          req.State,
		"address":        req.Address,
		"pincode":        req.Pincode,
		"business_name":  req.BusinessName,
		"business_type":  req.BusinessType,
		"gst_number":     req.GSTNumber,
		"mpin":           req.MPIN,
		"kyc_status":     req.KYCStatus,
		"documents_url":  req.DocumentsURL,
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

func (db *Database) DeleteRetailerQuery(
	ctx context.Context,
	retailerID string,
) error {

	query := `
		DELETE FROM retailers
		WHERE retailer_id = @retailer_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) ListRetailersQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetRetailerResponseModel, error) {

	query := `
		SELECT
			retailer_id,
			distributor_id,
			retailer_name,
			retailer_phone,
			retailer_email,
			retailer_aadhar_number,
			retailer_pan_number,
			retailer_date_of_birth,
			retailer_gender,
			retailer_city,
			retailer_state,
			retailer_address,
			retailer_pincode,
			retailer_business_name,
			retailer_business_type,
			retailer_gst_number,
			retailer_kyc_status,
			retailer_documents_url,
			retailer_wallet_balance,
			is_retailer_blocked,
			created_at,
			updated_at
		FROM retailers
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

	var list []models.GetRetailerResponseModel

	for rows.Next() {
		var r models.GetRetailerResponseModel
		err := rows.Scan(
			&r.RetailerID,
			&r.DistributorID,
			&r.Name,
			&r.Phone,
			&r.Email,
			&r.AadharNumber,
			&r.PanNumber,
			&r.DateOfBirth,
			&r.Gender,
			&r.City,
			&r.State,
			&r.Address,
			&r.Pincode,
			&r.BusinessName,
			&r.BusinessType,
			&r.GSTNumber,
			&r.KYCStatus,
			&r.DocumentsURL,
			&r.WalletBalance,
			&r.IsBlocked,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

func (db *Database) ListRetailersByDistributorIDQuery(
	ctx context.Context,
	distributorID string,
	limit, offset int,
) ([]models.GetRetailerResponseModel, error) {

	query := `
		SELECT
			retailer_id,
			distributor_id,
			retailer_name,
			retailer_phone,
			retailer_email,
			retailer_aadhar_number,
			retailer_pan_number,
			retailer_date_of_birth,
			retailer_gender,
			retailer_city,
			retailer_state,
			retailer_address,
			retailer_pincode,
			retailer_business_name,
			retailer_business_type,
			retailer_gst_number,
			retailer_kyc_status,
			retailer_documents_url,
			retailer_wallet_balance,
			is_retailer_blocked,
			created_at,
			updated_at
		FROM retailers
		WHERE distributor_id = @distributor_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
		"limit":          limit,
		"offset":         offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetRetailerResponseModel

	for rows.Next() {
		var r models.GetRetailerResponseModel
		if err := rows.Scan(
			&r.RetailerID,
			&r.DistributorID,
			&r.Name,
			&r.Phone,
			&r.Email,
			&r.AadharNumber,
			&r.PanNumber,
			&r.DateOfBirth,
			&r.Gender,
			&r.City,
			&r.State,
			&r.Address,
			&r.Pincode,
			&r.BusinessName,
			&r.BusinessType,
			&r.GSTNumber,
			&r.KYCStatus,
			&r.DocumentsURL,
			&r.WalletBalance,
			&r.IsBlocked,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

func (db *Database) ListRetailersByMasterDistributorIDQuery(
	ctx context.Context,
	masterDistributorID string,
	limit, offset int,
) ([]models.GetRetailerResponseModel, error) {

	query := `
		SELECT
			r.retailer_id,
			r.distributor_id,
			r.retailer_name,
			r.retailer_phone,
			r.retailer_email,
			r.retailer_aadhar_number,
			r.retailer_pan_number,
			r.retailer_date_of_birth,
			r.retailer_gender,
			r.retailer_city,
			r.retailer_state,
			r.retailer_address,
			r.retailer_pincode,
			r.retailer_business_name,
			r.retailer_business_type,
			r.retailer_gst_number,
			r.retailer_kyc_status,
			r.retailer_documents_url,
			r.retailer_wallet_balance,
			r.is_retailer_blocked,
			r.created_at,
			r.updated_at
		FROM retailers r
		INNER JOIN distributors d
			ON r.distributor_id = d.distributor_id
		WHERE d.master_distributor_id = @md_id
		ORDER BY r.created_at DESC
		LIMIT @limit OFFSET @offset
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"md_id":  masterDistributorID,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.GetRetailerResponseModel

	for rows.Next() {
		var r models.GetRetailerResponseModel
		if err := rows.Scan(
			&r.RetailerID,
			&r.DistributorID,
			&r.Name,
			&r.Phone,
			&r.Email,
			&r.AadharNumber,
			&r.PanNumber,
			&r.DateOfBirth,
			&r.Gender,
			&r.City,
			&r.State,
			&r.Address,
			&r.Pincode,
			&r.BusinessName,
			&r.BusinessType,
			&r.GSTNumber,
			&r.KYCStatus,
			&r.DocumentsURL,
			&r.WalletBalance,
			&r.IsBlocked,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

func (db *Database) GetRetailersByDistributorIDForDropdownQuery(
	ctx context.Context,
	distributorID string,
) ([]models.DropdownModel, error) {

	query := `
		SELECT
			retailer_id,
			retailer_name
		FROM retailers
		WHERE distributor_id = @distributor_id
		ORDER BY retailer_name ASC
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
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

func (db *Database) UpdateRetailerBlockStatusQuery(
	ctx context.Context,
	retailerID string,
	isBlocked bool,
) error {

	query := `
		UPDATE retailers
		SET
			is_retailer_blocked = @is_blocked,
			updated_at = NOW()
		WHERE retailer_id = @retailer_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
		"is_blocked":  isBlocked,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) UpdateRetailerKYCStatusQuery(
	ctx context.Context,
	retailerID string,
	kycStatus bool,
) error {

	query := `
		UPDATE retailers
		SET
			retailer_kyc_status = @kyc_status,
			updated_at = NOW()
		WHERE retailer_id = @retailer_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
		"kyc_status":  kycStatus,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
