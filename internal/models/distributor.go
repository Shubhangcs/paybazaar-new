package models

import "time"

type CreateDistributorRequestModel struct {
	MasterDistributorID string `json:"master_distributor_id" validate:"required"`

	Name     string `json:"name" validate:"required,min=3,max=100"`
	Phone    string `json:"phone" validate:"required,phone"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,strpwd"`

	AadharNumber string    `json:"aadhar_number" validate:"required,aadhar"`
	PanNumber    string    `json:"pan_number" validate:"required,pan"`
	DateOfBirth  time.Time `json:"date_of_birth" validate:"required"`
	Gender       string    `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`

	City    string `json:"city" validate:"required"`
	State   string `json:"state" validate:"required"`
	Address string `json:"address" validate:"required"`
	Pincode string `json:"pincode" validate:"required,len=6,numeric"`

	BusinessName string `json:"business_name" validate:"required"`
	BusinessType string `json:"business_type" validate:"required"`

	GSTNumber string `json:"gst_number" validate:"omitempty,len=15"`
	MPIN      int    `json:"mpin" validate:"omitempty,min=1000,max=9999"`
}

type UpdateDistributorRequestModel struct {
	Name     *string `json:"name" validate:"omitempty,min=3,max=100"`
	Phone    *string `json:"phone" validate:"omitempty,phone"`
	Password *string `json:"password" validate:"omitempty,strpwd"`

	City    *string `json:"city"`
	State   *string `json:"state"`
	Address *string `json:"address"`
	Pincode *string `json:"pincode" validate:"omitempty,len=6,numeric"`

	BusinessName *string `json:"business_name"`
	BusinessType *string `json:"business_type"`

	GSTNumber *string `json:"gst_number" validate:"omitempty,len=15"`
	MPIN      *int    `json:"mpin" validate:"omitempty,min=1000,max=9999"`

	KYCStatus    *bool   `json:"kyc_status"`
	DocumentsURL *string `json:"documents_url"`

	IsBlocked     *bool    `json:"is_blocked"`
	WalletBalance *float64 `json:"wallet_balance" validate:"omitempty,gte=0"`
}

type GetDistributorResponseModel struct {
	DistributorID       string `json:"distributor_id"`
	MasterDistributorID string `json:"master_distributor_id"`

	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`

	AadharNumber string    `json:"aadhar_number"`
	PanNumber    string    `json:"pan_number"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Gender       string    `json:"gender"`

	City    string `json:"city"`
	State   string `json:"state"`
	Address string `json:"address"`
	Pincode string `json:"pincode"`

	BusinessName string `json:"business_name"`
	BusinessType string `json:"business_type"`

	GSTNumber    string `json:"gst_number"`
	KYCStatus    bool   `json:"kyc_status"`
	DocumentsURL string `json:"documents_url"`

	WalletBalance float64 `json:"wallet_balance"`
	IsBlocked     bool    `json:"is_blocked"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DistributorModel struct {
	AdminID             string
	DistributorID       string
	MasterDistributorID string

	Name     string
	Phone    string
	Email    string
	Password string

	AadharNumber string
	PanNumber    string
	DateOfBirth  time.Time
	Gender       string

	City    string
	State   string
	Address string
	Pincode string

	BusinessName string
	BusinessType string

	GSTNumber    string
	MPIN         int
	KYCStatus    bool
	DocumentsURL string

	WalletBalance float64
	IsBlocked     bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginDistributorModel struct {
	DistributorID       string `json:"distributor_id" validate:"required"`
	DistributorPassword string `json:"distributor_password" validate:"required"`
}
