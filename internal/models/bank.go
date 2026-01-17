package models

type CreateBankModel struct {
	BankNams string `json:"bank_name" validate:"required"`
	IFSCCode string `json:"ifsc_code" validate:"required"`
}

type GetBankModel struct {
	BankID   int64  `json:"bank_id"`
	BankName string `json:"bank_name"`
	IFSCCode string `json:"ifsc_code"`
}

type UpdateBankModel struct {
	BankName *string `json:"bank_name,omitempty"`
	IFSCCode *string `json:"ifsc_code,omitempty"`
}

type CreateAdminBankModel struct {
	AdminID       string `json:"admin_id" validate:"required"`
	BankNams      string `json:"bank_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	IFSCCode      string `json:"ifsc_code" validate:"required"`
}

type GetAdminBankModel struct {
	AdminBankID   string `json:"admin_bank_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	IFSCCode      string `json:"ifsc_code"`
}

type UpdateAdminBankModel struct {
	BankName      *string `json:"bank_name,omitempty"`
	AccountNumber *string `json:"account_number,omitempty"`
	IFSCCode      *string `json:"ifsc_code,omitempty"`
}
