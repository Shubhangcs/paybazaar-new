package models

import "time"

type CreateMasterDistributorRequestModel struct {
	AdminID                   string    `json:"admin_id" validate:"required"`
	MasterDistributorName     string    `json:"master_distributor_name" validate:"required,min=3,max=100"`
	MasterDistributorPhone    string    `json:"master_distributor_phone" validate:"required,phone"`
	MasterDistributorEmail    string    `json:"master_distributor_email" validate:"required,email"`
	MasterDistributorPassword string    `json:"master_distributor_password" validate:"required,strpwd"`
	AadharNumber              string    `json:"aadhar_number" validate:"required,aadhar"`
	PanNumber                 string    `json:"pan_number" validate:"required,pan"`
	DateOfBirth               time.Time `json:"date_of_birth" validate:"required"`
	Gender                    string    `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	City                      string    `json:"city" validate:"required"`
	State                     string    `json:"state" validate:"required"`
	Address                   string    `json:"address" validate:"required"`
	Pincode                   string    `json:"pincode" validate:"required,len=6,numeric"`
	BusinessName              string    `json:"business_name" validate:"required"`
	BusinessType              string    `json:"business_type" validate:"required"`
	GSTNumber                 string    `json:"gst_number" validate:"omitempty,len=15"`
}

type UpdateMasterDistributorDetailsRequestModel struct {
	MasterDistributorID    string     `json:"master_distributor_id"`
	MasterDistributorName  *string    `json:"master_distributor_name"`
	MasterDistributorPhone *string    `json:"master_distributor_phone"`
	MasterDistributorEmail *string    `json:"master_distributor_email"`
	AadharNumber           *string    `json:"aadhar_number"`
	PanNumber              *string    `json:"pan_number"`
	DateOfBirth            *time.Time `json:"date_of_birth"`
	Gender                 *string    `json:"gender"`
	City                   *string    `json:"city"`
	State                  *string    `json:"state"`
	Address                *string    `json:"address"`
	Pincode                *string    `json:"pincode"`
	BusinessName           *string    `json:"business_name"`
	BusinessType           *string    `json:"business_type"`
	GSTNumber              *string    `json:"gst_number"`
	DocumentsURL           *string    `json:"documents_url"`
	WalletBalance          *float64   `json:"wallet_balance"`
}

type UpdateMasterDistributorPasswordRequestModel struct {
	MasterDistributorID string `json:"master_distributor_id" validate:"required"`
	OldPassword         string `json:"old_password" validate:"required,strpwd"`
	NewPassword         string `json:"new_password" validate:"required,strpwd"`
}

type UpdateMasterDistributorBlockStatusRequestModel struct {
	MasterDistributorID string `json:"master_distributor_id" validate:"required"`
	BlockStatus         bool   `json:"block_status"`
}

type UpdateMasterDistributorKYCStatusRequestModel struct {
	MasterDistributorID string `json:"master_distributor_id" validate:"required"`
	KYCStatus           bool   `json:"kyc_status"`
}

type UpdateMasterDistributorMPINRequestModel struct {
	MasterDistributorID string `json:"master_distributor_id" validate:"required"`
	OldMPIN             int64  `json:"old_mpin" validate:"required,min=1000,max=9999"`
	NewMPIN             int64  `json:"new_mpin" validate:"required,min=1000,max=9999"`
}

type LoginMasterDistributorRequestModel struct {
	MasterDistributorID       string `json:"master_distributor_id" validate:"required"`
	MasterDistributorPassword string `json:"master_distributor_password" validate:"required,strpwd"`
}

type GetMasterDistributorDetailsForLoginModel struct {
	AdminID               string `json:"admin_id"`
	MasterDistributorID   string `json:"master_distributor_id"`
	MasterDistributorName string `json:"master_distributor_name"`
	Password              string `json:"password"`
	IsBlocked             bool   `json:"is_blocked"`
}

type GetMasterDistributorForDropdownModel struct {
	MasterDistributorID   string `json:"master_distributor_id"`
	MasterDistributorName string `json:"master_distributor_name"`
}

type GetCompleteMasterDistributorDetailsResponseModel struct {
	MasterDistributorID       string    `json:"master_distributor_id"`
	AdminID                   string    `json:"admin_id"`
	MasterDistributorName     string    `json:"master_distributor_name"`
	MasterDistributorPhone    string    `json:"master_distributor_phone"`
	MasterDistributorEmail    string    `json:"master_distributor_email"`
	MasterDistributorPassword string    `json:"master_distributor_password"`
	AadharNumber              string    `json:"aadhar_number"`
	PanNumber                 string    `json:"pan_number"`
	DateOfBirth               time.Time `json:"date_of_birth"`
	Gender                    string    `json:"gender"`
	City                      string    `json:"city"`
	State                     string    `json:"state"`
	Address                   string    `json:"address"`
	Pincode                   string    `json:"pincode"`
	BusinessName              string    `json:"business_name"`
	BusinessType              string    `json:"business_type"`
	KYCStatus                 bool      `json:"kyc_status"`
	DocumentsURL              *string   `json:"documents_url"`
	GSTNumber                 *string   `json:"gst_number"`
	WalletBalance             float64   `json:"wallet_balance"`
	IsBlocked                 bool      `json:"is_blocked"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}
