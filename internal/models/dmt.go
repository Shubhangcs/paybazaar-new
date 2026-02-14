package models

type DMTWalletCheckRequestModel struct {
	MobileNumber string `json:"mobile_no" validate:"required"`
}

type DMTWalletCheckResponseModel struct {
	Error         int    `json:"error"`
	Message       string `json:"msg"`
	AccountExists int    `json:"AccountExists"`
	Description   string `json:"Description"`
}

type DMTCreateWalletRequestModel struct {
	RetailerID   string `json:"retailer_id" validate:"required"`
	MobileNumber string `json:"mobile_no" validate:"required"`
	Latitude     string `json:"lat" validate:"required"`
	Longitude    string `json:"long" validate:"required"`
	AadharNumber string `json:"aadhaar_number"`
	PidData      string `json:"pid_data" validate:"required"`
	IsIris       int    `json:"is_iris" validate:"required"`
}

type DMTCreateWalletResponseModel struct {
	Error         int    `json:"error"`
	Message       string `json:"msg"`
	AccountExists int    `json:"AccountExists"`
	Description   string `json:"Description"`
}

type DMTWalletVerificationRequestModel struct {
	RetailerID       string `json:"retailer_id" validate:"required"`
	MobileNumber     string `json:"mobile_no" validate:"required"`
	OTP              string `json:"otp" validate:"required"`
	EKycID           string `json:"ekyc_id" validate:"required"`
	StateResp        string `json:"stateresp" validate:"required"`
	PartnerRequestID string `json:"partner_request_id"`
}

type DMTWalletVerificationResponseModel struct {
	Error            int    `json:"error"`
	Message          string `json:"msg"`
	OrderID          string `json:"orderid"`
	PartnerRequestID string `json:"partnerreqid"`
	Description      string `json:"description"`
}

type DMTAddBeneficiaryRequestModel struct {
	MobileNumber     string `json:"mobile_no" validate:"required"`
	BeneficiaryName  string `json:"beneficiaryName" validate:"required"`
	AccountNumber    string `json:"accountNo" validate:"required"`
	IFSCCode         string `json:"ifsc" validate:"required"`
	BankID           string `json:"bankId" validate:"required"`
	PartnerRequestID string `json:"partner_request_id"`
}

type DMTAddBeneficiaryResponseModel struct {
	Error                 int     `json:"error"`
	Message               string  `json:"msg"`
	Status                int     `json:"status"`
	OrderID               *string `json:"orderid"`
	OperatorTransactionID *string `json:"optransid"`
	PartnerRequestID      string  `json:"partnerreqid"`
	Description           string  `json:"description"`
}

type DMTBankListResponseModel struct {
	Error    int    `json:"error"`
	Message  string `json:"message"`
	BankList []any  `json:"bankList"`
}

type DMTGetBeneficiaryRequestModel struct {
	MobileNumber string `json:"mobile_no" validate:"required"`
}

type DMTGetBeneficiaryResponseModel struct {
	Error           int    `json:"error"`
	Message         string `json:"msg"`
	MobileNumber    string `json:"MobileNo"`
	BeneficiaryList any    `json:"BeneficiaryList"`
}

type DMTTransactionRequestModel struct {
	RetailerID        string `json:"retailer_id"`
	MobileNumber      string `json:"mobile_number"`
	TransferType      string `json:"transfer_type"`
	Amount            int    `json:"amount"`
	BeneficiaryID     int    `json:"beneficiary_id"`
	BeneficiaryName   string `json:"beneficiary_name"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	IFSCCode          string `json:"ifsc_code"`
	Pincode           string `json:"pincode"`
	Address           string `json:"address"`
	PartnerRequestID  string `json:"partner_request_id"`
	TransactionStatus string `json:"transaction_status"`
}
