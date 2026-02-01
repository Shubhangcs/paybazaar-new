package models

type CheckDMTWalletExistsRequestModel struct {
	MobileNumber string `json:"mobile_no" validate:"required"`
}

type CheckDMTWalletExistsResponseModel struct {
	Error        int     `json:"error"`
	Message      string  `json:"msg"`
	FirstName    string  `json:"FirstName"`
	LastName     string  `json:"LastName"`
	MobileNumber string  `json:"MobileNo"`
	Limit        float64 `json:"Limit"`
	Description  string  `json:"description"`
}

type CreateDMTWalletRequestModel struct {
	MobileNumber string `json:"mobile_no" validate:"required"`
}

type CreateDMTWalletResponseModel struct {
	Error         int    `json:"error"`
	Message       string `json:"msg"`
	MobileNumber  string `json:"MobileNo"`
	RequestNumber string `json:"RequestNo"`
	Description   string `json:"description"`
}

type CreateDMTWalletVerifyRequestModel struct {
	RetailerID    string `json:"retailer_id" validate:"required"`
	MobileNumber  string `json:"mobile_no" validate:"required"`
	RequestNumber string `json:"request_no" validate:"required"`
	OTP           int    `json:"otp" validate:"required"`
	FirstName     string `json:"firstName" validate:"required"`
	LastName      string `json:"lastName" validate:"required"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	State         string `json:"state"`
	PinCode       string `json:"pinCode"`
}

type CreateDMTWalletVerifyResponseModel struct {
	Error        int    `json:"error"`
	Message      string `json:"msg"`
	MobileNumber string `json:"MobileNo"`
	Description  string `json:"description"`
}

type CreateDMTBeneficiaryRequestModel struct {
	MobileNumber    string `json:"mobile_no" validate:"required"`
	BeneficiaryName string `json:"beneficiaryName" validate:"required"`
	BankName        string `json:"bankName" validate:"required"`
	AccountNumber   string `json:"accountNo" validate:"required"`
	IFSCCode        string `json:"ifsc" validate:"required"`
}

type CreateDMTBeneficiaryResponseModel struct {
	Error         int    `json:"error"`
	Message       string `json:"msg"`
	MobileNumber  string `json:"MobileNo"`
	BeneficiaryID string `json:"BeneficiaryId"`
	Description   string `json:"description"`
}

type GetDMTBeneficiariesRequestModel struct {
	MobileNumber string `json:"mobile_no"`
}

type GetDMTBeneficiariesResponseModel struct {
	Error           int    `json:"error"`
	Message         string `json:"msg"`
	MobileNumber    string `json:"MobileNo"`
	Description     string `json:"Description"`
	BeneficiaryList any    `json:"BeneficiaryList"`
}

type DeleteDMTBeneficiaryRequestModel struct {
	MobileNumber string `json:"mobile_no" validate:"required"`
}

type DeleteDMTBeneficiaryResponseModel struct {
	Error         int    `json:"error"`
	Message       string `json:"msg"`
	MobileNumber  string `json:"MobileNo"`
	RequestNumber string `json:"RequestNo"`
	Description   string `json:"Description"`
}

type DeleteDMTBeneficiaryVerificationRequestModel struct {
	MobileNumber  string `json:"mobile_no" validate:"required"`
	RequestNumber string `json:"request_no" validate:"required"`
	OTP           int    `json:"otp" validate:"required"`
	BeneficiaryID string `json:"beneficiaryId" validate:"required"`
}

type DeleteDMTBeneficiaryVerificationResponseModel struct {
	Error        int    `json:"error"`
	Message      string `json:"msg"`
	MobileNumber string `json:"MobileNo"`
	Description  string `json:"description"`
}
