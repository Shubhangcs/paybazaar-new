package models

import "time"

type CreateRetailerRequestModel struct {
	DistributorID    string    `json:"distributor_id" validate:"required"`
	RetailerName     string    `json:"retailer_name" validate:"required,min=3,max=100"`
	RetailerPhone    string    `json:"retailer_phone" validate:"required,phone"`
	RetailerEmail    string    `json:"retailer_email" validate:"required,email"`
	RetailerPassword string    `json:"retailer_password" validate:"required,strpwd"`
	AadharNumber     string    `json:"aadhar_number" validate:"required,aadhar"`
	PanNumber        string    `json:"pan_number" validate:"required,pan"`
	DateOfBirth      time.Time `json:"date_of_birth" validate:"required"`
	Gender           string    `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	City             string    `json:"city" validate:"required"`
	State            string    `json:"state" validate:"required"`
	Address          string    `json:"address" validate:"required"`
	Pincode          string    `json:"pincode" validate:"required,len=6,numeric"`
	BusinessName     string    `json:"business_name" validate:"required"`
	BusinessType     string    `json:"business_type" validate:"required"`
	GSTNumber        string    `json:"gst_number" validate:"omitempty,len=15"`
}

type UpdateRetailerDetailsRequestModel struct {
	RetailerID    string     `json:"retailer_id"`
	RetailerName  *string    `json:"retailer_name"`
	RetailerPhone *string    `json:"retailer_phone"`
	RetailerEmail *string    `json:"retailer_email"`
	AadharNumber  *string    `json:"aadhar_number"`
	PanNumber     *string    `json:"pan_number"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Gender        *string    `json:"gender"`
	City          *string    `json:"city"`
	State         *string    `json:"state"`
	Address       *string    `json:"address"`
	Pincode       *string    `json:"pincode"`
	BusinessName  *string    `json:"business_name"`
	BusinessType  *string    `json:"business_type"`
	GSTNumber     *string    `json:"gst_number"`
	DocumentsURL  *string    `json:"documents_url"`
	WalletBalance *float64   `json:"wallet_balance"`
}

type UpdateRetailerPasswordRequestModel struct {
	RetailerID  string `json:"retailer_id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required,strpwd"`
	NewPassword string `json:"new_password" validate:"required,strpwd"`
}

type UpdateRetailerBlockStatusRequestModel struct {
	RetailerID  string `json:"retailer_id" validate:"required"`
	BlockStatus bool   `json:"block_status"`
}

type UpdateRetailerKYCStatusRequestModel struct {
	RetailerID string `json:"retailer_id" validate:"required"`
	KYCStatus  bool   `json:"kyc_status"`
}

type UpdateRetailerMPINRequestModel struct {
	RetailerID string `json:"retailer_id" validate:"required"`
	OldMPIN    int64  `json:"old_mpin" validate:"required,min=1000,max=9999"`
	NewMPIN    int64  `json:"new_mpin" validate:"required,min=1000,max=9999"`
}

type UpdateRetailerDistributorRequestModel struct {
	RetailerID    string `json:"retailer_id" validate:"required"`
	DistributorID string `json:"distributor_id" validate:"required"`
}

type RetailerLoginRequestModel struct {
	RetailerID       string `json:"retailer_id" validate:"required"`
	RetailerPassword string `json:"retailer_password" validate:"required,strpwd"`
}

type GetRetailerDetailsForLoginModel struct {
	AdminID      string `json:"admin_id"`
	RetailerID   string `json:"retailer_id"`
	RetailerName string `json:"retailer_name"`
	Password     string `json:"password"`
	IsBlocked    bool   `json:"is_blocked"`
}

type GetCompleteRetailerDetailsResponseModel struct {
	RetailerID       string    `json:"retailer_id"`
	DistributorID    string    `json:"distributor_id"`
	RetailerName     string    `json:"retailer_name"`
	RetailerPhone    string    `json:"retailer_phone"`
	RetailerEmail    string    `json:"retailer_email"`
	RetailerPassword string    `json:"retailer_password"`
	AadharNumber     string    `json:"aadhar_number"`
	PanNumber        string    `json:"pan_number"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	Gender           string    `json:"gender"`
	City             string    `json:"city"`
	State            string    `json:"state"`
	Address          string    `json:"address"`
	Pincode          string    `json:"pincode"`
	BusinessName     string    `json:"business_name"`
	BusinessType     string    `json:"business_type"`
	GSTNumber        string    `json:"gst_number"`
	KYCStatus        bool      `json:"kyc_status"`
	DocumentsURL     *string   `json:"documents_url"`
	WalletBalance    float64   `json:"wallet_balance"`
	IsBlocked        bool      `json:"is_blocked"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type GetRetailerForDropdownModel struct {
	RetailerID   string `json:"retailer_id"`
	RetailerName string `json:"retailer_name"`
}
