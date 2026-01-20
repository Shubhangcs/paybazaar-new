package database

import (
	"context"
	"fmt"

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
			master_distributor_gst_number
		) VALUES (
			@admin_id,
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

	if _, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":      req.AdminID,
		"name":          req.MasterDistributorName,
		"phone":         req.MasterDistributorPhone,
		"email":         req.MasterDistributorEmail,
		"password":      req.MasterDistributorPassword,
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
	}); err != nil {
		return fmt.Errorf("failed to create master distributor")
	}

	return nil
}

func (db *Database) GetMasterDistributorDetailsByMasterDistributorIDQuery(
	ctx context.Context,
	masterDistributorID string,
) (*models.GetCompleteMasterDistributorDetailsResponseModel, error) {

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
		WHERE master_distributor_id = @md_id;
	`

	var res models.GetCompleteMasterDistributorDetailsResponseModel
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"md_id": masterDistributorID,
	}).Scan(
		&res.MasterDistributorID,
		&res.AdminID,
		&res.MasterDistributorName,
		&res.MasterDistributorPhone,
		&res.MasterDistributorEmail,
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
		&res.KYCStatus,
		&res.DocumentsURL,
		&res.GSTNumber,
		&res.WalletBalance,
		&res.IsBlocked,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch master distributor details")
	}

	return &res, nil
}

func (db *Database) GetMasterDistributorDetailsForLoginQuery(
	ctx context.Context,
	masterDistributorID string,
) (*models.GetMasterDistributorDetailsForLoginModel, error) {

	query := `
		SELECT
			master_distributor_id,
			admin_id,
			master_distributor_name,
			master_distributor_password,
			is_master_distributor_blocked
		FROM master_distributors
		WHERE master_distributor_id = @md_id;
	`

	var res models.GetMasterDistributorDetailsForLoginModel
	if err := db.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"md_id": masterDistributorID,
	}).Scan(
		&res.MasterDistributorID,
		&res.AdminID,
		&res.MasterDistributorName,
		&res.Password,
		&res.IsBlocked,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch master distributor login details")
	}

	return &res, nil
}

func (db *Database) UpdateMasterDistributorDetailsQuery(
	ctx context.Context,
	req models.UpdateMasterDistributorDetailsRequestModel,
) error {

	query := `
		UPDATE master_distributors
		SET
			master_distributor_name = COALESCE(@name, master_distributor_name),
			master_distributor_phone = COALESCE(@phone, master_distributor_phone),
			master_distributor_email = COALESCE(@email, master_distributor_email),

			master_distributor_aadhar_number = COALESCE(@aadhar, master_distributor_aadhar_number),
			master_distributor_pan_number = COALESCE(@pan, master_distributor_pan_number),
			master_distributor_date_of_birth = COALESCE(@dob, master_distributor_date_of_birth),
			master_distributor_gender = COALESCE(@gender, master_distributor_gender),

			master_distributor_city = COALESCE(@city, master_distributor_city),
			master_distributor_state = COALESCE(@state, master_distributor_state),
			master_distributor_address = COALESCE(@address, master_distributor_address),
			master_distributor_pincode = COALESCE(@pincode, master_distributor_pincode),

			master_distributor_business_name = COALESCE(@business_name, master_distributor_business_name),
			master_distributor_business_type = COALESCE(@business_type, master_distributor_business_type),

			master_distributor_gst_number = COALESCE(@gst, master_distributor_gst_number),
			master_distributor_documents_url = COALESCE(@documents_url, master_distributor_documents_url),
			updated_at = NOW()
		WHERE master_distributor_id = @md_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id":         req.MasterDistributorID,
		"name":          req.MasterDistributorName,
		"phone":         req.MasterDistributorPhone,
		"email":         req.MasterDistributorEmail,
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
	})

	if err != nil {
		return fmt.Errorf("failed to update master distributor details")
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid master distributor id or not found")
	}
	return nil
}

func (db *Database) UpdateMasterDistributorPasswordQuery(
	ctx context.Context,
	req models.UpdateMasterDistributorPasswordRequestModel,
) error {

	var oldPassword string
	getQuery := `
		SELECT master_distributor_password
		FROM master_distributors
		WHERE master_distributor_id = @md_id;
	`

	if err := db.pool.QueryRow(ctx, getQuery, pgx.NamedArgs{
		"md_id": req.MasterDistributorID,
	}).Scan(&oldPassword); err != nil {
		return fmt.Errorf("failed to fetch old password")
	}

	if oldPassword != req.OldPassword {
		return fmt.Errorf("incorrect old password")
	}

	updateQuery := `
		UPDATE master_distributors
		SET master_distributor_password = @new_password,
		    updated_at = NOW()
		WHERE master_distributor_id = @md_id;
	`

	if _, err := db.pool.Exec(ctx, updateQuery, pgx.NamedArgs{
		"md_id":        req.MasterDistributorID,
		"new_password": req.NewPassword,
	}); err != nil {
		return fmt.Errorf("failed to update master distributor password")
	}

	return nil
}

func (db *Database) UpdateMasterDistributorBlockStatusQuery(
	ctx context.Context,
	req models.UpdateMasterDistributorBlockStatusRequestModel,
) error {

	query := `
		UPDATE master_distributors
		SET is_master_distributor_blocked = @block_status,
		    updated_at = NOW()
		WHERE master_distributor_id = @md_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id":        req.MasterDistributorID,
		"block_status": req.BlockStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to update master distributor block status")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid master distributor id or not found")
	}

	return nil
}

func (db *Database) UpdateMasterDistributorKYCStatusQuery(
	ctx context.Context,
	req models.UpdateMasterDistributorKYCStatusRequestModel,
) error {

	query := `
		UPDATE master_distributors
		SET master_distributor_kyc_status = @kyc_status,
		    updated_at = NOW()
		WHERE master_distributor_id = @md_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id":      req.MasterDistributorID,
		"kyc_status": req.KYCStatus,
	})
	if err != nil {
		return fmt.Errorf("failed to update master distributor kyc status")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid master distributor id or not found")
	}

	return nil
}

func (db *Database) UpdateMasterDistributorMPINQuery(
	ctx context.Context,
	req models.UpdateMasterDistributorMPINRequestModel,
) error {

	var existingMPIN int64
	getQuery := `
		SELECT master_distributor_mpin
		FROM master_distributors
		WHERE master_distributor_id = @md_id;
	`

	if err := db.pool.QueryRow(ctx, getQuery, pgx.NamedArgs{
		"md_id": req.MasterDistributorID,
	}).Scan(&existingMPIN); err != nil {
		return fmt.Errorf("failed to fetch master distributor mpin")
	}

	if existingMPIN != req.OldMPIN {
		return fmt.Errorf("incorrect old mpin")
	}

	updateQuery := `
		UPDATE master_distributors
		SET master_distributor_mpin = @new_mpin,
		    updated_at = NOW()
		WHERE master_distributor_id = @md_id;
	`

	tag, err := db.pool.Exec(ctx, updateQuery, pgx.NamedArgs{
		"md_id":    req.MasterDistributorID,
		"new_mpin": req.NewMPIN,
	})
	if err != nil {
		return fmt.Errorf("failed to update master distributor mpin")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid master distributor id or not found")
	}

	return nil
}

func (db *Database) DeleteMasterDistributorQuery(
	ctx context.Context,
	masterDistributorID string,
) error {

	query := `
		DELETE FROM master_distributors
		WHERE master_distributor_id = @md_id;
	`

	tag, err := db.pool.Exec(ctx, query, pgx.NamedArgs{
		"md_id": masterDistributorID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete master distributor")
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid master distributor id or not found")
	}

	return nil
}

func (db *Database) GetMasterDistributorsByAdminIDQuery(
	ctx context.Context,
	adminID string,
	limit, offset int,
) ([]models.GetCompleteMasterDistributorDetailsResponseModel, error) {

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
		LIMIT @limit OFFSET @offset;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
		"limit":    limit,
		"offset":   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch master distributors by admin id")
	}
	defer rows.Close()

	var list []models.GetCompleteMasterDistributorDetailsResponseModel

	for rows.Next() {
		var md models.GetCompleteMasterDistributorDetailsResponseModel
		if err := rows.Scan(
			&md.MasterDistributorID,
			&md.AdminID,
			&md.MasterDistributorName,
			&md.MasterDistributorPhone,
			&md.MasterDistributorEmail,
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
			return nil, fmt.Errorf("failed to scan master distributor details")
		}
		list = append(list, md)
	}

	return list, nil
}

func (db *Database) GetMasterDistributorsForDropdownByAdminIDQuery(
	ctx context.Context,
	adminID string,
) ([]models.GetMasterDistributorForDropdownModel, error) {

	query := `
		SELECT
			master_distributor_id,
			master_distributor_name
		FROM master_distributors
		WHERE admin_id = @admin_id
		ORDER BY master_distributor_name ASC;
	`

	rows, err := db.pool.Query(ctx, query, pgx.NamedArgs{
		"admin_id": adminID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch master distributors for dropdown")
	}
	defer rows.Close()

	var list []models.GetMasterDistributorForDropdownModel

	for rows.Next() {
		var md models.GetMasterDistributorForDropdownModel
		if err := rows.Scan(
			&md.MasterDistributorID,
			&md.MasterDistributorName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan master distributor dropdown data")
		}
		list = append(list, md)
	}

	return list, nil
}
