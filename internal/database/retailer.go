package database

import (
	"context"
	"fmt"

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
			retailer_gst_number
		) VALUES (
			@distributor_id,
			@name,
			@phone,
			@email,
			@password,
			@aadhar,
			@pan,
			@dob,
			@gender,
			@city,
			@state,
			@address,
			@pincode,
			@business_name,
			@business_type,
			@gst
		);
	`

	_, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": req.DistributorID,
		"name":           req.RetailerName,
		"phone":          req.RetailerPhone,
		"email":          req.RetailerEmail,
		"password":       req.RetailerPassword,
		"aadhar":         req.AadharNumber,
		"pan":            req.PanNumber,
		"dob":            req.DateOfBirth,
		"gender":         req.Gender,
		"city":           req.City,
		"state":          req.State,
		"address":        req.Address,
		"pincode":        req.Pincode,
		"business_name":  req.BusinessName,
		"business_type":  req.BusinessType,
		"gst":            req.GSTNumber,
	})
	if err != nil {
		return fmt.Errorf("failed to create retailer")
	}

	return nil
}

func (db *Database) GetRetailerDetailsByRetailerIDQuery(
	ctx context.Context,
	retailerID string,
) (*models.GetCompleteRetailerDetailsResponseModel, error) {

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
		WHERE retailer_id = @retailer_id;
	`

	var r models.GetCompleteRetailerDetailsResponseModel
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
	}).Scan(
		&r.RetailerID,
		&r.DistributorID,
		&r.RetailerName,
		&r.RetailerPhone,
		&r.RetailerEmail,
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
		return nil, fmt.Errorf("failed to fetch retailer details")
	}

	return &r, nil
}

func (db *Database) GetRetailersByDistributorIDQuery(
	ctx context.Context,
	distributorID string,
	limit, offset int,
) ([]models.GetCompleteRetailerDetailsResponseModel, error) {

	query := `
		SELECT
			retailer_id,
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
			retailer_kyc_status,
			retailer_documents_url,
			retailer_wallet_balance,
			is_retailer_blocked,
			created_at,
			updated_at
		FROM retailers
		WHERE distributor_id = @distributor_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
		"limit":          limit,
		"offset":         offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch retailers by distributor id")
	}
	defer rows.Close()

	var list []models.GetCompleteRetailerDetailsResponseModel
	for rows.Next() {
		var r models.GetCompleteRetailerDetailsResponseModel
		if err := rows.Scan(
			&r.RetailerID,
			&r.DistributorID,
			&r.RetailerName,
			&r.RetailerPhone,
			&r.RetailerEmail,
			&r.RetailerPassword,
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
			return nil, fmt.Errorf("failed to scan retailer")
		}
		list = append(list, r)
	}

	return list, nil
}

func (db *Database) GetRetailersByMasterDistributorIDQuery(
	ctx context.Context,
	masterDistributorID string,
	limit, offset int,
) ([]models.GetCompleteRetailerDetailsResponseModel, error) {

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
			r.retailer_kyc_status,
			r.retailer_documents_url,
			r.retailer_wallet_balance,
			r.is_retailer_blocked,
			r.created_at,
			r.updated_at
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		WHERE d.master_distributor_id = @md_id
		ORDER BY r.created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"md_id":  masterDistributorID,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch retailers by master distributor id")
	}
	defer rows.Close()

	var list []models.GetCompleteRetailerDetailsResponseModel
	for rows.Next() {
		var r models.GetCompleteRetailerDetailsResponseModel
		if err := rows.Scan(
			&r.RetailerID,
			&r.DistributorID,
			&r.RetailerName,
			&r.RetailerPhone,
			&r.RetailerEmail,
			&r.RetailerPassword,
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
			return nil, fmt.Errorf("failed to scan retailer")
		}
		list = append(list, r)
	}

	return list, nil
}

func (db *Database) GetRetailersByAdminIDQuery(
	ctx context.Context,
	adminID string,
	limit, offset int,
) ([]models.GetCompleteRetailerDetailsResponseModel, error) {

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
			r.retailer_kyc_status,
			r.retailer_documents_url,
			r.retailer_wallet_balance,
			r.is_retailer_blocked,
			r.created_at,
			r.updated_at
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
		WHERE md.admin_id = @admin_id
		ORDER BY r.created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
		"limit":    limit,
		"offset":   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch retailers by admin id")
	}
	defer rows.Close()

	var list []models.GetCompleteRetailerDetailsResponseModel
	for rows.Next() {
		var r models.GetCompleteRetailerDetailsResponseModel
		if err := rows.Scan(
			&r.RetailerID,
			&r.DistributorID,
			&r.RetailerName,
			&r.RetailerPhone,
			&r.RetailerEmail,
			&r.RetailerPassword,
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
			return nil, fmt.Errorf("failed to scan retailer")
		}
		list = append(list, r)
	}

	return list, nil
}

func (db *Database) GetRetailerDetailsForLoginQuery(
	ctx context.Context,
	retailerID string,
) (*models.GetRetailerDetailsForLoginModel, error) {

	query := `
		SELECT
			r.retailer_id,
			r.retailer_name,
			r.retailer_password,
			r.is_retailer_blocked,
			md.admin_id
		FROM retailers r
		JOIN distributors d ON r.distributor_id = d.distributor_id
		JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
		WHERE r.retailer_id = @retailer_id;
	`

	var res models.GetRetailerDetailsForLoginModel
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
	}).Scan(
		&res.RetailerID,
		&res.RetailerName,
		&res.Password,
		&res.IsBlocked,
		&res.AdminID,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch retailer login details")
	}

	return &res, nil
}

func (db *Database) UpdateRetailerDetailsQuery(
	ctx context.Context,
	req models.UpdateRetailerDetailsRequestModel,
) error {

	query := `
		UPDATE retailers
		SET
			retailer_name = COALESCE(@name, retailer_name),
			retailer_phone = COALESCE(@phone, retailer_phone),
			retailer_email = COALESCE(@email, retailer_email),

			retailer_aadhar_number = COALESCE(@aadhar, retailer_aadhar_number),
			retailer_pan_number = COALESCE(@pan, retailer_pan_number),
			retailer_date_of_birth = COALESCE(@dob, retailer_date_of_birth),
			retailer_gender = COALESCE(@gender, retailer_gender),

			retailer_city = COALESCE(@city, retailer_city),
			retailer_state = COALESCE(@state, retailer_state),
			retailer_address = COALESCE(@address, retailer_address),
			retailer_pincode = COALESCE(@pincode, retailer_pincode),

			retailer_business_name = COALESCE(@business_name, retailer_business_name),
			retailer_business_type = COALESCE(@business_type, retailer_business_type),

			retailer_gst_number = COALESCE(@gst, retailer_gst_number),
			retailer_documents_url = COALESCE(@documents_url, retailer_documents_url),
			retailer_wallet_balance = COALESCE(@wallet_balance , retailer_wallet_balance),
			updated_at = NOW()
		WHERE retailer_id = @retailer_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id":   req.RetailerID,
		"name":          req.RetailerName,
		"phone":         req.RetailerPhone,
		"email":         req.RetailerEmail,
		"aadhar":        req.AadharNumber,
		"pan":           req.PanNumber,
		"dob":           req.DateOfBirth,
		"gender":        req.Gender,
		"city":          req.City,
		"state":         req.State,
		"address":       req.Address,
		"pincode":       req.Pincode,
		"business_name": req.BusinessName,
		"business_type": req.BusinessType,
		"gst":           req.GSTNumber,
		"documents_url": req.DocumentsURL,
		"wallet_balance": req.WalletBalance,
	})

	if err != nil {
		return fmt.Errorf("failed to update retailer details")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid retailer id or retailer not found")
	}

	return nil
}

func (db *Database) UpdateRetailerPasswordQuery(
	ctx context.Context,
	req models.UpdateRetailerPasswordRequestModel,
) error {

	var oldPassword string
	getQuery := `
		SELECT retailer_password
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`

	if err := db.pool.QueryRow(ctx, getQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(&oldPassword); err != nil {
		return fmt.Errorf("failed to fetch old password")
	}

	if oldPassword != req.OldPassword {
		return fmt.Errorf("incorrect old password")
	}

	updateQuery := `
		UPDATE retailers
		SET retailer_password = @new_password,
		    updated_at = NOW()
		WHERE retailer_id = @retailer_id;
	`

	if _, err := db.pool.Exec(ctx, updateQuery, pgx.NamedArgs{
		"retailer_id":  req.RetailerID,
		"new_password": req.NewPassword,
	}); err != nil {
		return fmt.Errorf("failed to update retailer password")
	}

	return nil
}

func (db *Database) UpdateRetailerBlockStatusQuery(
	ctx context.Context,
	req models.UpdateRetailerBlockStatusRequestModel,
) error {

	query := `
		UPDATE retailers
		SET is_retailer_blocked = @block_status,
		    updated_at = NOW()
		WHERE retailer_id = @retailer_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id":  req.RetailerID,
		"block_status": req.BlockStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to update retailer block status")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid retailer id or retailer not found")
	}
	return nil
}

func (db *Database) UpdateRetailerKYCStatusQuery(
	ctx context.Context,
	req models.UpdateRetailerKYCStatusRequestModel,
) error {

	query := `
		UPDATE retailers
		SET retailer_kyc_status = @kyc_status,
		    updated_at = NOW()
		WHERE retailer_id = @retailer_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
		"kyc_status":  req.KYCStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to update retailer kyc status")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid retailer id or retailer not found")
	}
	return nil
}

func (db *Database) UpdateRetailerMPINQuery(
	ctx context.Context,
	req models.UpdateRetailerMPINRequestModel,
) error {

	var existingMPIN int64
	getQuery := `
		SELECT retailer_mpin
		FROM retailers
		WHERE retailer_id = @retailer_id;
	`

	if err := db.pool.QueryRow(ctx, getQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
	}).Scan(&existingMPIN); err != nil {
		return fmt.Errorf("failed to fetch retailer mpin")
	}

	if existingMPIN != req.OldMPIN {
		return fmt.Errorf("incorrect old mpin")
	}

	updateQuery := `
		UPDATE retailers
		SET retailer_mpin = @new_mpin,
		    updated_at = NOW()
		WHERE retailer_id = @retailer_id;
	`

	tag, err := db.pool.Exec(ctx, updateQuery, pgx.NamedArgs{
		"retailer_id": req.RetailerID,
		"new_mpin":    req.NewMPIN,
	})
	if err != nil {
		return fmt.Errorf("failed to update retailer mpin")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid retailer id or retailer not found")
	}

	return nil
}

func (db *Database) UpdateRetailerDistributorQuery(
	ctx context.Context,
	req models.UpdateRetailerDistributorRequestModel,
) error {

	query := `
		UPDATE retailers
		SET distributor_id = @distributor_id,
		    updated_at = NOW()
		WHERE retailer_id = @retailer_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id":    req.RetailerID,
		"distributor_id": req.DistributorID,
	})
	if err != nil {
		return fmt.Errorf("failed to update retailer distributor")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid retailer id or retailer not found")
	}

	return nil
}

func (db *Database) DeleteRetailerQuery(
	ctx context.Context,
	retailerID string,
) error {

	query := `
		DELETE FROM retailers
		WHERE retailer_id = @retailer_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"retailer_id": retailerID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete retailer")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid retailer id or retailer not found")
	}

	return nil
}

func (db *Database) GetRetailersForDropdownByDistributorIDQuery(
	ctx context.Context,
	distributorID string,
) ([]models.GetRetailerForDropdownModel, error) {

	query := `
		SELECT
			retailer_id,
			retailer_name
		FROM retailers
		WHERE distributor_id = @distributor_id
		ORDER BY retailer_name ASC;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch retailers for dropdown")
	}
	defer rows.Close()

	var list []models.GetRetailerForDropdownModel

	for rows.Next() {
		var r models.GetRetailerForDropdownModel
		if err := rows.Scan(
			&r.RetailerID,
			&r.RetailerName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan retailer dropdown data")
		}
		list = append(list, r)
	}

	return list, nil
}
