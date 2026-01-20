package models

import "time"

type CreatePayoutBeneficiaryModel struct {
	RetailerID      string `json:"retailer_id" validate:"required"`
	MobileNumber    string `json:"mobile_number" validate:"required"`
	BankName        string `json:"bank_name" validate:"required"`
	BeneficiaryName string `json:"beneficiary_name" validate:"required"`
	AccountNumber   string `json:"account_number" validate:"required"`
	IFSCCode        string `json:"ifsc_code" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
}

type UpdatePayoutBeneficiaryModel struct {
	BankName        *string `json:"bank_name,omitempty"`
	BeneficiaryName *string `json:"beneficiary_name,omitempty"`
	AccountNumber   *string `json:"account_number,omitempty"`
	IFSCCode        *string `json:"ifsc_code,omitempty"`
	Phone           *string `json:"phone,omitempty"`
}

type UpdatePayoutVerificationRequestModel struct {
	RetailerID    string `json:"retailer_id"`
	PhoneNumber   string `json:"retailer_phone_number"`
	AccountNumber string `json:"account_number"`
	IFSCCode      string `json:"ifsc_code"`
}

type GetPayoutBeneficiaryResponseModel struct {
	BeneficiaryID   int64     `json:"beneficiary_id"`
	RetailerID      string    `json:"retailer_id"`
	MobileNumber    string    `json:"mobile_number"`
	BankName        string    `json:"bank_name"`
	BeneficiaryName string    `json:"beneficiary_name"`
	AccountNumber   string    `json:"account_number"`
	IFSCCode        string    `json:"ifsc_code"`
	Phone           string    `json:"phone"`
	IsVerified      bool      `json:"is_beneficiary_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
