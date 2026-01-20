package models

import "time"

type CreateDistributorRequestModel struct {
	MasterDistributorID string    `json:"master_distributor_id" validate:"required"`
	DistributorName     string    `json:"distributor_name" validate:"required,min=3,max=100"`
	DistributorPhone    string    `json:"distributor_phone" validate:"required,phone"`
	DistributorEmail    string    `json:"distributor_email" validate:"required,email"`
	DistributorPassword string    `json:"distributor_password" validate:"required,strpwd"`
	AadharNumber        string    `json:"aadhar_number" validate:"required,aadhar"`
	PanNumber           string    `json:"pan_number" validate:"required,pan"`
	DateOfBirth         time.Time `json:"date_of_birth" validate:"required"`
	Gender              string    `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	City                string    `json:"city" validate:"required"`
	State               string    `json:"state" validate:"required"`
	Address             string    `json:"address" validate:"required"`
	Pincode             string    `json:"pincode" validate:"required,len=6,numeric"`
	BusinessName        string    `json:"business_name" validate:"required"`
	BusinessType        string    `json:"business_type" validate:"required"`
	GSTNumber           string    `json:"gst_number" validate:"omitempty,len=15"`
}

type UpdateDistributorDetailsRequestModel struct {
	DistributorID    string     `json:"distributor_id" validate:"required"`
	DistributorName  *string    `json:"distributor_name" validate:"omitempty,min=3,max=100"`
	DistributorEmail *string    `json:"distributor_email" validate:"omitempty,email"`
	DistributorPhone *string    `json:"distributor_phone" validate:"omitempty,phone"`
	AadharNumber     *string    `json:"aadhar_number" validate:"omitempty,aadhar"`
	PanNumber        *string    `json:"pan_number" validate:"omitempty,pan"`
	DateOfBirth      *time.Time `json:"date_of_birth" validate:"omitempty"`
	Gender           *string    `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	City             *string    `json:"city" validate:"omitempty"`
	State            *string    `json:"state" validate:"omitempty"`
	Address          *string    `json:"address" validate:"omitempty"`
	Pincode          *string    `json:"pincode" validate:"omitempty,len=6,numeric"`
	BusinessName     *string    `json:"business_name" validate:"omitempty"`
	BusinessType     *string    `json:"business_type" validate:"omitempty"`
	GSTNumber        *string    `json:"gst_number" validate:"omitempty,len=15"`
	DocumentsURL     *string    `json:"documents_url" validate:"omitempty"`
}

type UpdateDistributorPasswordRequestModel struct {
	DistributorID string `json:"distributor_id" validate:"required"`
	OldPassword   string `json:"old_password" validate:"required,strpwd"`
	NewPassword   string `json:"new_password" validate:"required,strpwd"`
}

type UpdateDistributorBlockStatusRequestModel struct {
	DistributorID string `json:"distributor_id" validate:"required"`
	BlockStatus   bool   `json:"block_status"`
}

type UpdateDistributorKYCStatusRequestModel struct {
	DistributorID string `json:"distributor_id" validate:"required"`
	KYCStatus     bool   `json:"kyc_status" validate:"required"`
}

type UpdateDistributorMasterDistributorRequestModel struct {
	DistributorID       string `json:"distributor_id" validate:"required"`
	MasterDistributorID string `json:"master_distributor_id" validate:"required"`
}

type UpdateDistributorMPINRequestModel struct {
	DistributorID string `json:"distributor_id" validate:"required"`
	OldMPIN       int64  `json:"old_mpin" validate:"required,min=1000,max=9999"`
	NewMPIN       int64  `json:"new_mpin" validate:"required,min=1000,max=9999"`
}

type DistributorLoginRequestModel struct {
	DistributorID       string `json:"distributor_id" validate:"required"`
	DistributorPassword string `json:"distributor_password" validate:"required,strpwd"`
}

type GetDistributorDetailsForLoginModel struct {
	AdminID              string `json:"admin_id"`
	DistributorID        string `json:"distributor_id"`
	DistributorName      string `json:"distributor_name"`
	DistributorPassword  string `json:"distributor_password"`
	IsDistributorBlocked bool   `json:"is_distributor_blocked"`
}

type GetCompleteDistributorDetailsResponseModel struct {
	DistributorID       string    `json:"distributor_id"`
	MasterDistributorID string    `json:"master_distributor_id"`
	DistributorName     string    `json:"distributor_name"`
	DistributorPhone    string    `json:"distributor_phone"`
	DistributorEmail    string    `json:"distributor_email"`
	AadharNumber        string    `json:"aadhar_number"`
	PanNumber           string    `json:"pan_number"`
	DateOfBirth         time.Time `json:"date_of_birth"`
	Gender              string    `json:"gender"`
	City                string    `json:"city"`
	State               string    `json:"state"`
	Address             string    `json:"address"`
	Pincode             string    `json:"pincode"`
	BusinessName        string    `json:"business_name"`
	BusinessType        string    `json:"business_type"`
	GSTNumber           string    `json:"gst_number"`
	KYCStatus           bool      `json:"kyc_status"`
	DocumentsURL        *string    `json:"documents_url"`
	WalletBalance       float64   `json:"wallet_balance"`
	IsBlocked           bool      `json:"is_blocked"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type GetDistributorForDropdownModel struct {
	DistributorID   string `json:"distributor_id"`
	DistributorName string `json:"distributor_name"`
}
