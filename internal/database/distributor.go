package database

import (
	"context"
	"fmt"

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
			distributor_gst_number
		) VALUES (
			@master_distributor_id,
			@distributor_name,
			@distributor_phone,
			@distributor_email,
			@distributor_password,
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
			@gst_number
		)
	`

	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"master_distributor_id": req.MasterDistributorID,
		"distributor_name":      req.DistributorName,
		"distributor_phone":     req.DistributorPhone,
		"distributor_email":     req.DistributorEmail,
		"distributor_password":  req.DistributorPassword,
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
	}); err != nil {
		return fmt.Errorf("failed to create distributor")
	}

	return nil
}

func (db *Database) GetDistributorDetailsByDistributorIDQuery(
	ctx context.Context,
	distributorID string,
) (*models.GetCompleteDistributorDetailsResponseModel, error) {

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
		WHERE distributor_id = @distributor_id;
	`

	var res models.GetCompleteDistributorDetailsResponseModel
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
	}).Scan(
		&res.DistributorID,
		&res.MasterDistributorID,
		&res.DistributorName,
		&res.DistributorPhone,
		&res.DistributorEmail,
		&res.AadharNumber,
		&res.PanNumber,
		&res.DateOfBirth,
		&res.Gender,
		&res.City,
		&res.State,
		&res.Address,
		&res.Pincode,
		&res.BusinessName,
		&res.BusinessType,
		&res.GSTNumber,
		&res.KYCStatus,
		&res.DocumentsURL,
		&res.WalletBalance,
		&res.IsBlocked,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch distributor details")
	}

	return &res, nil
}

func (db *Database) GetDistributorDetailsForLoginQuery(
	ctx context.Context,
	distributorID string,
) (*models.GetDistributorDetailsForLoginModel, error) {

	query := `
		SELECT
			d.distributor_id,
			d.distributor_name,
			d.distributor_password,
			d.is_distributor_blocked,
			md.admin_id
		FROM distributors d
		JOIN master_distributors md
			ON d.master_distributor_id = md.master_distributor_id
		WHERE d.distributor_id = @distributor_id;
	`

	var res models.GetDistributorDetailsForLoginModel
	if err := db.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{
			"distributor_id": distributorID,
		},
	).Scan(
		&res.DistributorID,
		&res.DistributorName,
		&res.DistributorPassword,
		&res.IsDistributorBlocked,
		&res.AdminID,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch distributor login details")
	}

	return &res, nil
}

func (db *Database) UpdateDistributorPasswordQuery(
	ctx context.Context,
	req models.UpdateDistributorPasswordRequestModel,
) error {

	query := `
		SELECT distributor_password
		FROM distributors
		WHERE distributor_id = @distributor_id;
	`

	var oldPassword string
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"distributor_id": req.DistributorID,
	}).Scan(&oldPassword); err != nil {
		return fmt.Errorf("failed to fetch old password")
	}

	if oldPassword != req.OldPassword {
		return fmt.Errorf("incorrect old password")
	}

	updateQuery := `
		UPDATE distributors
		SET distributor_password = @new_password,
		    updated_at = NOW()
		WHERE distributor_id = @distributor_id;
	`

	if _, err := db.pool.Exec(ctx, updateQuery, pgx.NamedArgs{
		"distributor_id": req.DistributorID,
		"new_password":   req.NewPassword,
	}); err != nil {
		return fmt.Errorf("failed to update distributor password")
	}

	return nil
}

func (db *Database) UpdateDistributorBlockStatusQuery(
	ctx context.Context,
	req models.UpdateDistributorBlockStatusRequestModel,
) error {

	query := `
		UPDATE distributors
		SET is_distributor_blocked = @block_status,
		    updated_at = NOW()
		WHERE distributor_id = @distributor_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": req.DistributorID,
		"block_status":   req.BlockStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to update distributor block status")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid distributor id or distributor not found")
	}
	return nil
}

func (db *Database) UpdateDistributorKYCStatusQuery(
	ctx context.Context,
	req models.UpdateDistributorKYCStatusRequestModel,
) error {

	query := `
		UPDATE distributors
		SET distributor_kyc_status = @kyc_status,
		    updated_at = NOW()
		WHERE distributor_id = @distributor_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": req.DistributorID,
		"kyc_status":     req.KYCStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to update distributor kyc status")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid distributor id or distributor not found")
	}
	return nil
}

func (db *Database) UpdateDistributorDetailsQuery(
	ctx context.Context,
	req models.UpdateDistributorDetailsRequestModel,
) error {

	query := `
		UPDATE distributors
		SET
			distributor_name = COALESCE(@distributor_name, distributor_name),
			distributor_email = COALESCE(@distributor_email, distributor_email),
			distributor_phone = COALESCE(@distributor_phone, distributor_phone),

			distributor_aadhar_number = COALESCE(@aadhar_number, distributor_aadhar_number),
			distributor_pan_number = COALESCE(@pan_number, distributor_pan_number),
			distributor_date_of_birth = COALESCE(@date_of_birth, distributor_date_of_birth),
			distributor_gender = COALESCE(@gender, distributor_gender),

			distributor_city = COALESCE(@city, distributor_city),
			distributor_state = COALESCE(@state, distributor_state),
			distributor_address = COALESCE(@address, distributor_address),
			distributor_pincode = COALESCE(@pincode, distributor_pincode),

			distributor_business_name = COALESCE(@business_name, distributor_business_name),
			distributor_business_type = COALESCE(@business_type, distributor_business_type),
			distributor_gst_number = COALESCE(@gst_number, distributor_gst_number),

			distributor_documents_url = COALESCE(@documents_url, distributor_documents_url),
			distributor_wallet_balance = COALESCE(@distributor_wallet_balance , distributor_wallet_balance),
			updated_at = NOW()
		WHERE distributor_id = @distributor_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id":    req.DistributorID,
		"distributor_name":  req.DistributorName,
		"distributor_email": req.DistributorEmail,
		"distributor_phone": req.DistributorPhone,

		"aadhar_number": req.AadharNumber,
		"pan_number":    req.PanNumber,
		"date_of_birth": req.DateOfBirth,
		"gender":        req.Gender,

		"city":    req.City,
		"state":   req.State,
		"address": req.Address,
		"pincode": req.Pincode,

		"business_name": req.BusinessName,
		"business_type": req.BusinessType,
		"gst_number":    req.GSTNumber,

		"documents_url": req.DocumentsURL,
		"distributor_wallet_balance": req.WalletBalance,
	})

	if err != nil {
		return fmt.Errorf("failed to update distributor details")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid distributor id or distributor not found")
	}
	return nil
}

func (db *Database) GetDistributorsByAdminIDQuery(
	ctx context.Context,
	adminID string,
	limit, offset int,
) ([]models.GetCompleteDistributorDetailsResponseModel, error) {

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
			d.distributor_kyc_status,
			d.distributor_documents_url,
			d.distributor_wallet_balance,
			d.is_distributor_blocked,
			d.created_at,
			d.updated_at
		FROM distributors d
		JOIN master_distributors md
			ON d.master_distributor_id = md.master_distributor_id
		WHERE md.admin_id = @admin_id
		ORDER BY d.created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
		"limit":    limit,
		"offset":   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch distributors by admin id")
	}
	defer rows.Close()

	var list []models.GetCompleteDistributorDetailsResponseModel

	for rows.Next() {
		var d models.GetCompleteDistributorDetailsResponseModel
		if err := rows.Scan(
			&d.DistributorID,
			&d.MasterDistributorID,
			&d.DistributorName,
			&d.DistributorPhone,
			&d.DistributorEmail,
			&d.DistributorPassword,
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
			return nil, fmt.Errorf("failed to scan distributor details")
		}
		list = append(list, d)
	}

	return list, nil
}

func (db *Database) GetDistributorsByMasterDistributorIDQuery(
	ctx context.Context,
	masterDistributorID string,
	limit, offset int,
) ([]models.GetCompleteDistributorDetailsResponseModel, error) {

	query := `
		SELECT
			distributor_id,
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
			distributor_kyc_status,
			distributor_documents_url,
			distributor_wallet_balance,
			is_distributor_blocked,
			created_at,
			updated_at
		FROM distributors
		WHERE master_distributor_id = @md_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"md_id":  masterDistributorID,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch distributors by master distributor id")
	}
	defer rows.Close()

	var list []models.GetCompleteDistributorDetailsResponseModel

	for rows.Next() {
		var d models.GetCompleteDistributorDetailsResponseModel
		if err := rows.Scan(
			&d.DistributorID,
			&d.MasterDistributorID,
			&d.DistributorName,
			&d.DistributorPhone,
			&d.DistributorEmail,
			&d.DistributorPassword,
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
			return nil, fmt.Errorf("failed to scan distributor details")
		}
		list = append(list, d)
	}

	return list, nil
}

func (db *Database) GetDistributorsForDropdownByMasterDistributorIDQuery(
	ctx context.Context,
	masterDistributorID string,
) ([]models.GetDistributorForDropdownModel, error) {

	query := `
		SELECT
			distributor_id,
			distributor_name
		FROM distributors
		WHERE master_distributor_id = @md_id
		ORDER BY distributor_name ASC;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"md_id": masterDistributorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch distributors for dropdown")
	}
	defer rows.Close()

	var list []models.GetDistributorForDropdownModel

	for rows.Next() {
		var d models.GetDistributorForDropdownModel
		if err := rows.Scan(
			&d.DistributorID,
			&d.DistributorName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan distributor dropdown data")
		}
		list = append(list, d)
	}

	return list, nil
}

func (db *Database) UpdateDistributorMasterDistributorQuery(
	ctx context.Context,
	req models.UpdateDistributorMasterDistributorRequestModel,
) error {

	query := `
		UPDATE distributors
		SET
			master_distributor_id = @master_distributor_id,
			updated_at = NOW()
		WHERE distributor_id = @distributor_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id":        req.DistributorID,
		"master_distributor_id": req.MasterDistributorID,
	})

	if err != nil {
		return fmt.Errorf("failed to update distributor master distributor")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid distributor id or distributor not found")
	}

	return nil
}

func (db *Database) DeleteDistributorQuery(
	ctx context.Context,
	distributorID string,
) error {

	query := `
		DELETE FROM distributors
		WHERE distributor_id = @distributor_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"distributor_id": distributorID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete distributor")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid distributor id or distributor not found")
	}

	return nil
}

func (db *Database) UpdateDistributorMPINQuery(
	ctx context.Context,
	req models.UpdateDistributorMPINRequestModel,
) error {

	var existingMPIN int64
	getMPINQuery := `
		SELECT distributor_mpin
		FROM distributors
		WHERE distributor_id = @distributor_id;
	`

	if err := db.pool.QueryRow(
		ctx,
		getMPINQuery,
		pgx.NamedArgs{
			"distributor_id": req.DistributorID,
		},
	).Scan(&existingMPIN); err != nil {
		return fmt.Errorf("failed to fetch distributor mpin")
	}

	if existingMPIN != req.OldMPIN {
		return fmt.Errorf("incorrect old mpin")
	}

	updateQuery := `
		UPDATE distributors
		SET
			distributor_mpin = @new_mpin,
			updated_at = NOW()
		WHERE distributor_id = @distributor_id;
	`

	tag, err := db.pool.Exec(
		ctx,
		updateQuery,
		pgx.NamedArgs{
			"distributor_id": req.DistributorID,
			"new_mpin":       req.NewMPIN,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update distributor mpin")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid distributor id or distributor not found")
	}

	return nil
}
