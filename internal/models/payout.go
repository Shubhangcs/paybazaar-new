package models

import "time"

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

type GetPayoutTransactionModel struct {
	PayoutTransactionID string `json:"payout_transaction_id"`
	PartnerRequestID    string `json:"partner_request_id"`
	OperatorTxnID       string `json:"operator_transaction_id"`

	RetailerID          string `json:"retailer_id"`
	RetailerName        string `json:"retailer_name"`
	RetailerBusinessName string `json:"retailer_business_name"`

	OrderID      string `json:"order_id"`
	MobileNumber string `json:"mobile_number"`

	BeneficiaryBankName  string `json:"beneficiary_bank_name"`
	BeneficiaryName      string `json:"beneficiary_name"`
	BeneficiaryAccountNo string `json:"beneficiary_account_number"`
	BeneficiaryIFSCCode  string `json:"beneficiary_ifsc_code"`

	Amount       float64 `json:"amount"`
	TransferType string  `json:"transfer_type"`

	AdminCommision             float64 `json:"admin_commision"`
	MasterDistributorCommision float64 `json:"master_distributor_commision"`
	DistributorCommision       float64 `json:"distributor_commision"`
	RetailerCommision          float64 `json:"retailer_commision"`

	PayoutStatus string    `json:"payout_transaction_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GetRetailerPayoutModel struct {
	PayoutTransactionID string `json:"payout_transaction_id"`
	PartnerRequestID    string `json:"partner_request_id"`
	OperatorTxnID       string `json:"operator_transaction_id"`

	RetailerID          string `json:"retailer_id"`
	RetailerName        string `json:"retailer_name"`
	RetailerBusinessName string `json:"retailer_business_name"`

	OrderID      string `json:"order_id"`
	MobileNumber string `json:"mobile_number"`

	BeneficiaryBankName string `json:"beneficiary_bank_name"`
	BeneficiaryName     string `json:"beneficiary_name"`
	AccountNumber       string `json:"beneficiary_account_number"`
	IFSCCode            string `json:"beneficiary_ifsc_code"`

	Amount       float64 `json:"amount"`
	TransferType string  `json:"transfer_type"`

	RetailerCommision float64 `json:"retailer_commision"`

	Status    string    `json:"payout_transaction_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PayoutLedgerWithWalletResponseModel struct {
	// ---------- PAYOUT ----------
	PayoutTransactionID   string  `json:"payout_transaction_id"`
	PartnerRequestID      string  `json:"partner_request_id"`
	OperatorTransactionID string  `json:"operator_transaction_id"`

	RetailerID           string  `json:"retailer_id"`
	RetailerName         *string `json:"retailer_name,omitempty"`
	RetailerBusinessName *string `json:"retailer_business_name,omitempty"`

	OrderID                  string  `json:"order_id"`
	MobileNumber             string  `json:"mobile_number"`
	BeneficiaryBankName      string  `json:"beneficiary_bank_name"`
	BeneficiaryName          string  `json:"beneficiary_name"`
	BeneficiaryAccountNumber string  `json:"beneficiary_account_number"`
	BeneficiaryIFSCCode      string  `json:"beneficiary_ifsc_code"`
	Amount                   float64 `json:"amount"`
	TransferType             string  `json:"transfer_type"`

	AdminCommission             float64 `json:"admin_commision"`
	MasterDistributorCommission float64 `json:"master_distributor_commision"`
	DistributorCommission       float64 `json:"distributor_commision"`
	RetailerCommission          float64 `json:"retailer_commision"`

	PayoutStatus    string    `json:"payout_transaction_status"`
	PayoutCreatedAt time.Time `json:"payout_created_at"`
	PayoutUpdatedAt time.Time `json:"payout_updated_at"`

	// ---------- TDS ----------
	TDSCommissionID  *int64     `json:"tds_commision_id,omitempty"`
	TDSTransactionID *string    `json:"tds_transaction_id,omitempty"`
	TDSUserID        *string    `json:"tds_user_id,omitempty"`
	TDSUserName      *string    `json:"tds_user_name,omitempty"`
	TDSCommission    *float64   `json:"tds_commision,omitempty"`
	TDSAmount        *float64   `json:"tds,omitempty"`
	CommissionNet    *float64   `json:"paid_commision,omitempty"`
	PANNumber        *string    `json:"pan_number,omitempty"`
	TDSStatus        *string    `json:"tds_status,omitempty"`
	TDSCreatedAt     *time.Time `json:"tds_created_at,omitempty"`

	// ---------- WALLET ----------
	WalletTransactionID *int64     `json:"wallet_transaction_id,omitempty"`
	WalletUserID        *string    `json:"wallet_user_id,omitempty"`
	WalletReferenceID   *string    `json:"wallet_reference_id,omitempty"`
	CreditAmount        *float64   `json:"credit_amount,omitempty"`
	DebitAmount         *float64   `json:"debit_amount,omitempty"`
	BeforeBalance       *float64   `json:"before_balance,omitempty"`
	AfterBalance        *float64   `json:"after_balance,omitempty"`
	TransactionReason   *string    `json:"transaction_reason,omitempty"`
	Remarks             *string    `json:"remarks,omitempty"`
	WalletCreatedAt     *time.Time `json:"wallet_created_at,omitempty"`
}
