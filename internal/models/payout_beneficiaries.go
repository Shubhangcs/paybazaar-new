package models

type BeneficiaryModel struct {
	BeneficiaryID       string `json:"beneficiary_id"`
	MobileNumber        string `json:"mobile_number"`
	BankName            string `json:"bank_name"`
	IFSCCode            string `json:"ifsc_code"`
	AccountNumber       string `json:"account_number"`
	BeneficiaryName     string `json:"beneficiary_name"`
	BeneficiaryPhone    string `json:"beneficiary_phone"`
	BeneficiaryVerified bool   `json:"beneficiary_verified"`
}

type VerifyBeneficiaryRequestModel struct {
	RetailerId    string `json:"retailer_id" validate:"required"`
	MobileNumber  string `json:"mobile_number" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	IFSCCode      string `json:"ifsc_code" validate:"required"`
}

type VerifyBeneficiaryResponseModel struct {
	StatusCode int                           `json:"status_code"`
	Status     bool                          `json:"status"`
	Message    string                        `json:"message"`
	Data       VerifyBeneficiaryDetailsModel `json:"data"`
}

type VerifyBeneficiaryDetailsModel struct {
	UserName   string `json:"c_name"`
	BankName   string `json:"bank_name"`
	BranchName string `json:"branch_name"`
}
