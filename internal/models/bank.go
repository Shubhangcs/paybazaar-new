package models

type CreateBankRequestModel struct {
	BankNams string `json:"bank_name" validate:"required"`
	IFSCCode string `json:"ifsc_code" validate:"required"`
}

type UpdateBankDetailsRequestModel struct {
	BankID   int64   `json:"bank_id" validate:"required"`
	BankName *string `json:"bank_name" validate:"omitempty"`
	IFSCCode *string `json:"ifsc_code" validate:"omitempty"`
}

type GetBankDetailsResponseModel struct {
	BankID   int64  `json:"bank_id"`
	BankName string `json:"bank_name"`
	IFSCCode string `json:"ifsc_code"`
}

type CreateAdminBankRequestModel struct {
	AdminID       string `json:"admin_id" validate:"required"`
	BankNams      string `json:"bank_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	IFSCCode      string `json:"ifsc_code" validate:"required"`
}

type UpdateAdminBankDetailsRequestModel struct {
	AdminBankID   int64   `json:"admin_bank_id" validate:"required"`
	BankName      *string `json:"bank_name" validate:"omitempty"`
	AccountNumber *string `json:"account_number" validate:"omitempty"`
	IFSCCode      *string `json:"ifsc_code" validate:"omitempty"`
}

type GetAdminBankDetailsResponseModel struct {
	AdminBankID   int64 `json:"admin_bank_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	IFSCCode      string `json:"ifsc_code"`
}
