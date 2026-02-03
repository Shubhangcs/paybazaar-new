package models

import "time"

type CreatePayoutRequestModel struct {
	RetailerId            string  `json:"retailer_id" validate:"required"`
	MobileNumber          string  `json:"mobile_number" validate:"required"`
	IFSCCode              string  `json:"ifsc_code" validate:"required"`
	BankName              string  `json:"bank_name" validate:"required"`
	AccountNumber         string  `json:"account_number" validate:"required"`
	BeneficiaryName       string  `json:"beneficiary_name" validate:"required"`
	Amount                float64 `json:"amount" validate:"required"`
	TransferType          int     `json:"transfer_type" validate:"required"`
	PartnerRequestId      string  `json:"partner_request_id"`
	OrderId               string  `json:"order_id"`
	OperatorTransactionId string  `json:"operator_transaction_id"`
	TransactionStatus     string  `json:"transaction_status"`
}

type GetPayoutCommisionModel struct {
	TotalCommision             float64
	AdminCommision             float64
	MasterDistributorCommision float64
	DistributorCommision       float64
	RetailerCommision          float64
}

type GetAllPayoutTransactionsResponseModel struct {
	PayoutTransactionId        string    `json:"payout_transaction_id"`
	OperatorTransactionId      *string   `json:"operator_transaction_id"`
	PartnerRequestId           string    `json:"partner_request_id"`
	OrderId                    *string   `json:"order_id"`
	RetailerId                 string    `json:"retailer_id"`
	RetailerName               string    `json:"retailer_name"`
	RetailerBusinessName       string    `json:"retailer_business_name"`
	MobileNumber               string    `json:"mobile_number"`
	BankName                   string    `json:"bank_name"`
	BeneficiaryName            string    `json:"beneficiary_name"`
	AccountNumber              string    `json:"account_number"`
	IFSCCode                   string    `json:"ifsc_code"`
	Amount                     float64   `json:"amount"`
	TransferType               string    `json:"transfer_type"`
	TransactionStatus          string    `json:"transaction_status"`
	AdminCommision             float64   `json:"admin_commision"`
	MasterDistributorCommision float64   `json:"master_distributor_commision"`
	DistributorCommision       float64   `json:"distributor_commision"`
	RetailerCommision          float64   `json:"retailer_commision"`
	BeforeBalance              float64   `json:"before_balance"`
	AfterBalance               float64   `json:"after_balance"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type GetRetailerPayoutTransactionsResponseModel struct {
	PayoutTransactionId   string    `json:"payout_transaction_id"`
	OperatorTransactionId *string   `json:"operator_transaction_id"`
	PartnerRequestId      string    `json:"partner_request_id"`
	OrderId               *string   `json:"order_id"`
	RetailerId            string    `json:"retailer_id"`
	RetailerName          string    `json:"retailer_name"`
	RetailerBusinessName  string    `json:"retailer_business_name"`
	MobileNumber          string    `json:"mobile_number"`
	BankName              string    `json:"bank_name"`
	BeneficiaryName       string    `json:"beneficiary_name"`
	AccountNumber         string    `json:"account_number"`
	IFSCCode              string    `json:"ifsc_code"`
	Amount                float64   `json:"amount"`
	TransferType          string    `json:"transfer_type"`
	TransactionStatus     string    `json:"transaction_status"`
	RetailerCommision     float64   `json:"retailer_commision"`
	BeforeBalance         float64   `json:"before_balance"`
	AfterBalance          float64   `json:"after_balance"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
