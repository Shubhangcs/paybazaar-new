package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

func (db *Database) CreateDistributorQuery(
	ctx context.Context,
	req models.CreateDistributorRequestModel,
) error {

	query := `
		INSERT INTO distributors (
			master_distributor_id,
			distributor_name,
			distributor_phone,
			distributor_email,
			distributor_password,
			distributor_aadhar_number,
			distributor_pan_number,
			distributor_date_of_birth,
			distributor_gender,
			distributor_city,
			distributor_state,
			distributor_address,
			distributor_pincode,
			distributor_business_name,
			distributor_business_type,
			distributor_gst_number,
			distributor_mpin
		) VALUES (
			@master_distributor_id,
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
		"master_distributor_id": req.MasterDistributorID,
		"name":                  req.Name,
		"phone":                 req.Phone,
		"email":                 req.Email,
		"password":              req.Password,
		"aadhar_number":         req.AadharNumber,
		"pan_number":            req.PanNumber,
		"date_of_birth":         req.DateOfBirth,
		"gender":                req.Gender,
		"city":                  req.City,
		"state":                 req.State,
		"address":               req.Address,
		"pincode":               req.Pincode,
		"business_name":         req.BusinessName,
		"business_type":         req.BusinessType,
		"gst_number":            req.GSTNumber,
		"mpin":                  req.MPIN,
	})

	return err
}

func (db *Database) GetDistributorByIDQuery(
	ctx context.Context,
	distributorID string,
) (*models.DistributorModel, error) {

	query := `
		SELECT
			d.distributor_id,
			d.master_distributor_id,
			d.distributor_name,
			d.distributor_phone,
			d.distributor_email,
			d.distributor_password,
			d.distributor_aadhar_number,
			d.distributor_pan_number,
			d.distributor_date_of_birth,
			d.distributor_gender,
			d.distributor_city,
			d.distributor_state,
			d.distributor_address,
			d.distributor_pincode,
			d.distributor_business_name,
			d.distributor_business_type,
			d.distributor_gst_number,
			d.distributor_mpin,
			d.distributor_kyc_status,
			d.distributor_documents_url,
			d.distributor_wallet_balance,
			d.is_distributor_blocked,
			d.created_at,
			d.updated_at,
			md.admin_id
		FROM distributors d
		JOIN master_distributors md
			ON d.master_distributor_id = md.master_distributor_id
		WHERE d.distributor_id = @distributor_id;
	`

	row := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"distributor_id": distributorID,
		},
	)

	var d models.DistributorModel
	err := row.Scan(
		&d.DistributorID,
		&d.MasterDistributorID,
		&d.Name,
		&d.Phone,
		&d.Email,
		&d.Password,
		&d.AadharNumber,
		&d.PanNumber,
		&d.DateOfBirth,
		&d.Gender,
		&d.City,
		&d.State,
		&d.Address,
		&d.Pincode,
		&d.BusinessName,
		&d.BusinessType,
		&d.GSTNumber,
		&d.MPIN,
		&d.KYCStatus,
		&d.DocumentsURL,
		&d.WalletBalance,
		&d.IsBlocked,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.AdminID,
	)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (db *Database) GetDistributorByEmailQuery(
	ctx context.Context,
	email string,
) (*models.DistributorModel, error) {

	query := `
		SELECT
			d.distributor_id,
			d.master_distributor_id,
			d.distributor_name,
			d.distributor_password,
			d.is_distributor_blocked,
			md.admin_id
		FROM distributors d
		JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
		WHERE d.distributor_email = @email
	`

	row := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"email": email,
	})

	var d models.DistributorModel
	err := row.Scan(
		&d.DistributorID,
		&d.MasterDistributorID,
		&d.Name,
		&d.Password,
		&d.IsBlocked,
		&d.AdminID,
	)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (db *Database) UpdateDistributorQuery(
	ctx context.Context,
	distributorID string,
	req models.UpdateDistributorRequestModel,
) error {

	query := `
		UPDATE distributors
		SET
			distributor_name = COALESCE(@name, distributor_name),
			distributor_phone = COALESCE(@phone, distributor_phone),
			distributor_password = COALESCE(@password, distributor_password),
			distributor_city = COALESCE(@city, distributor_city),
			distributor_state = COALESCE(@state, distributor_state),
			distributor_address = COALESCE(@address, distributor_address),
			distributor_pincode = COALESCE(@pincode, distributor_pincode),
			distributor_business_name = COALESCE(@business_name, distributor_business_name),
			distributor_business_type = COALESCE(@business_type, distributor_business_type),
			distributor_gst_number = COALESCE(@gst_number, distributor_gst_number),
			distributor_mpin = COALESCE(@mpin, distributor_mpin),
			distributor_kyc_status = COALESCE(@kyc_status, distributor_kyc_status),
			distributor_documents_url = COALESCE(@documents_url, distributor_documents_url),
			distributor_wallet_balance = COALESCE(@wallet_balance, distributor_wallet_balance),
			is_distributor_blocked = COALESCE(@is_blocked, is_distributor_blocked),
			updated_at = NOW()
		WHERE distributor_id = @distributor_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
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

func (db *Database) DeleteDistributorQuery(
	ctx context.Context,
	distributorID string,
) error {

	query := `
		DELETE FROM distributors
		WHERE distributor_id = @distributor_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
	})

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) ListDistributorsQuery(
	ctx context.Context,
	limit, offset int,
) ([]models.GetDistributorResponseModel, error) {

	query := `
		SELECT
			distributor_id,
			master_distributor_id,
			distributor_name,
			distributor_phone,
			distributor_email,
			distributor_aadhar_number,
			distributor_pan_number,
			distributor_date_of_birth,
			distributor_gender,
			distributor_city,
			distributor_state,
			distributor_address,
			distributor_pincode,
			distributor_business_name,
			distributor_business_type,
			distributor_gst_number,
			distributor_kyc_status,
			distributor_documents_url,
			distributor_wallet_balance,
			is_distributor_blocked,
			created_at,
			updated_at
		FROM distributors
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

	var list []models.GetDistributorResponseModel

	for rows.Next() {
		var d models.GetDistributorResponseModel
		err := rows.Scan(
			&d.DistributorID,
			&d.MasterDistributorID,
			&d.Name,
			&d.Phone,
			&d.Email,
			&d.AadharNumber,
			&d.PanNumber,
			&d.DateOfBirth,
			&d.Gender,
			&d.City,
			&d.State,
			&d.Address,
			&d.Pincode,
			&d.BusinessName,
			&d.BusinessType,
			&d.GSTNumber,
			&d.KYCStatus,
			&d.DocumentsURL,
			&d.WalletBalance,
			&d.IsBlocked,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}

func (db *Database) ListDistributorsByMasterDistributorIDQuery(
	ctx context.Context,
	masterDistributorID string,
	limit, offset int,
) ([]models.GetDistributorResponseModel, error) {

	query := `
		SELECT
			distributor_id,
			master_distributor_id,
			distributor_name,
			distributor_phone,
			distributor_email,
			distributor_aadhar_number,
			distributor_pan_number,
			distributor_date_of_birth,
			distributor_gender,
			distributor_city,
			distributor_state,
			distributor_address,
			distributor_pincode,
			distributor_business_name,
			distributor_business_type,
			distributor_gst_number,
			distributor_kyc_status,
			distributor_documents_url,
			distributor_wallet_balance,
			is_distributor_blocked,
			created_at,
			updated_at
		FROM distributors
		WHERE master_distributor_id = @md_id
		ORDER BY created_at DESC
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

	var list []models.GetDistributorResponseModel

	for rows.Next() {
		var d models.GetDistributorResponseModel
		if err := rows.Scan(
			&d.DistributorID,
			&d.MasterDistributorID,
			&d.Name,
			&d.Phone,
			&d.Email,
			&d.AadharNumber,
			&d.PanNumber,
			&d.DateOfBirth,
			&d.Gender,
			&d.City,
			&d.State,
			&d.Address,
			&d.Pincode,
			&d.BusinessName,
			&d.BusinessType,
			&d.GSTNumber,
			&d.KYCStatus,
			&d.DocumentsURL,
			&d.WalletBalance,
			&d.IsBlocked,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}

func (db *Database) GetDistributorsByMasterDistributorIDForDropdownQuery(
	ctx context.Context,
	masterDistributorID string,
) ([]models.DropdownModel, error) {

	query := `
		SELECT
			distributor_id,
			distributor_name
		FROM distributors
		WHERE master_distributor_id = @md_id
		ORDER BY distributor_name ASC
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"md_id": masterDistributorID,
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

func (db *Database) UpdateDistributorBlockStatusQuery(
	ctx context.Context,
	distributorID string,
	isBlocked bool,
) error {

	query := `
		UPDATE distributors
		SET
			is_distributor_blocked = @is_blocked,
			updated_at = NOW()
		WHERE distributor_id = @distributor_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
		"is_blocked":     isBlocked,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (db *Database) UpdateDistributorKYCStatusQuery(
	ctx context.Context,
	distributorID string,
	kycStatus bool,
) error {

	query := `
		UPDATE distributors
		SET
			distributor_kyc_status = @kyc_status,
			updated_at = NOW()
		WHERE distributor_id = @distributor_id
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
		"kyc_status":     kycStatus,
	})
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
