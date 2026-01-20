package models

type CreatePayoutRequestModel struct {
	AdminID                  string  `json:"admin_id" validate:"required"`
	RetailerID               string  `json:"retailer_id" validate:"required"`
	MobileNumber             string  `json:"mobile_number" validate:"required"`
	BeneficiaryBankName      string  `json:"beneficiary_bank_name" validate:"required"`
	BeneficiaryName          string  `json:"beneficiary_name" validate:"required"`
	BeneficiaryAccountNumber string  `json:"beneficiary_account_number" validate:"required"`
	BeneficiaryIFSCCode      string  `json:"beneficiary_ifsc_code" validate:"required"`
	Amount                   float64 `json:"amount" validate:"required"`
	TransferType             string  `json:"transfer_type" validate:"required"`
	MPIN                     int64   `json:"mpin" validate:"required"`
	PartnerRequestID         string  `json:"partner_req_id,omitempty"`
}

type PayoutAPIResponseModel struct {
	Status                int    `json:"status"`
	OrderID               string `json:"orderid"`
	OperatorTransactionID string `json:"optransid"`
	PartnerRequestID      string `json:"partnerreqid"`
}
