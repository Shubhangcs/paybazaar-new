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
	RetailerID   string  `json:"retailer_id" validate:"required"`
	MobileNumber string  `json:"mobile_no" validate:"required"`
	Latitude     float64 `json:"lat" validate:"required"`
	Longitude    float64 `json:"long" validate:"required"`
	AadharNumber string  `json:"aadhar_number"`
	PidData      string  `json:"pid_data" validate:"required"`
	IsIris       int     `json:"is_iris" validate:"required"`
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
	StateResp        string `json:"stateresp"`
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
