package models

import "time"

type GetPostpaidMobileRechargeBillFetchAPIRequestModel struct {
	MobileNumber string `json:"mobile_no" validate:"required"`
	OperatorCode int    `json:"operator_code" validate:"required"`
}

type GetPostpaidMobileRechargeBillFetchAPIResponseModel struct {
	Error      int    `json:"error"`
	Message    string `json:"msg"`
	Status     int    `json:"status"`
	BillAmount any    `json:"billAmount"`
}

type CreatePostpaidMobileRechargeAPIRequestModel struct {
	RetailerID       string  `json:"retailer_id" validate:"required"`
	MobileNumber     string  `json:"mobile_no" validate:"required"`
	OperatorCode     int     `json:"operator_code" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	OperatorCircle   int     `json:"circle" validate:"required"`
	PartnerRequestID string  `json:"partner_request_id,omitempty"`
	OperatorName     string  `json:"operator_name" validate:"required"`
	CircleName       string  `json:"circle_name" validate:"required"`
}

type GetPostpaidMobileRechargeAPIResponseModel struct {
	Error                 int    `json:"error"`
	Message               string `json:"msg"`
	Status                int    `json:"status"`
	OrderID               string `json:"orderid"`
	OperatorTransactionID string `json:"optransid"`
}

type GetPostpaidMobileRechargeHistoryResponseModel struct {
	PostpaidRechargeTransactionID int       `json:"postpaid_recharge_transaction_id"`
	RetailerID                    string    `json:"retailer_id"`
	RetailerName                  string    `json:"retailer_name"`
	RetailerBusinessName          string    `json:"retailer_business_name"`
	PartnerRequestID              string    `json:"partner_request_id"`
	OperatorTransactionID         string    `json:"operator_transaction_id"`
	OrderID                       string    `json:"order_id"`
	MobileNumber                  string    `json:"mobile_number"`
	OperatorCode                  string    `json:"operator_code"`
	Amount                        float64   `json:"amount"`
	BeforeBalance                 float64   `json:"before_balance"`
	AfterBalance                  float64   `json:"after_balance"`
	CircleCode                    string    `json:"circle_code"`
	CircleName                    string    `json:"circle_name"`
	OperatorName                  string    `json:"operator_name"`
	RechargeType                  string    `json:"recharge_type"`
	RechargeStatus                string    `json:"recharge_status"`
	Commission                    float64   `json:"commission"`
	CreatedAt                     time.Time `json:"created_at"`
}

type GetElectricityBillFetchRequestModel struct {
	CustomerID   string `json:"customer_id" validate:"required"`
	OperatorCode int    `json:"operator_code" validate:"required"`
}

type GetElectricityBillFetchResponseModel struct {
	Error      int    `json:"error"`
	Message    string `json:"msg"`
	Status     int    `json:"status"`
	BillAmount any    `json:"billAmount"`
}

type CreateElectricityBillPaymentRequestModel struct {
	RetailerID       string  `json:"retailer_id"`
	CustomerID       string  `json:"customer_id" validate:"required"`
	CustomerEmail    string  `json:"customer_email" validate:"required"`
	OperatorCode     string  `json:"operator_code" validate:"required"`
	OperatorName     string  `json:"operator_name" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	PartnerRequestID string  `json:"partner_request_id,omitempty"`
}

type GetElectricityBillPaymentAPIResponseModel struct {
	Error                 int    `json:"error"`
	Status                int    `json:"status"`
	Message               string `json:"msg"`
	OrderID               string `json:"orderid"`
	PartnerRequestID      string `json:"partnerreqid"`
	OperatorTransactionID string `json:"optransid"`
}

type GetElectricityOperatorResponseModel struct {
	OperatorName string `json:"operator_name"`
	OperatorCode int    `json:"operator_code"`
}
