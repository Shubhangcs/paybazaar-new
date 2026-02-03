package models

import "time"

type CreateMobileRechargeRequestModel struct {
	RetailerID       string  `json:"retailer_id"`
	MobileNumber     int64   `json:"mobile_number" validate:"required"`
	OperatorCode     int     `json:"operator_code" validate:"required"`
	OperatorName     string  `json:"operator_name" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	CircleCode       int     `json:"circle_code" validate:"required"`
	CircleName       string  `json:"circle_name" validate:"required"`
	RechargeType     string  `json:"recharge_type,omitempty"`
	PartnerRequestID string  `json:"partner_request_id,omitempty"`
	Commision        float64 `json:"commision"`
	Status           string  `json:"status"`
}

type GetMobileRechargeHistoryResponseModel struct {
	RetailerID                  string    `json:"retailer_id"`
	RetailerName                string    `json:"retailer_name"`
	BusinessName                string    `json:"business_name"`
	MobileRechargeTransactionID int       `json:"mobile_recharge_transaction_id"`
	MobileNumber                int64     `json:"mobile_number"`
	OperatorCode                int       `json:"operator_code"`
	OperatorName                string    `json:"operator_name"`
	Amount                      float64   `json:"amount"`
	CircleCode                  int       `json:"circle_code"`
	CircleName                  string    `json:"circle_name"`
	RechargeType                string    `json:"recharge_type"`
	PartnerRequestID            string    `json:"partner_request_id"`
	CreatedAt                   time.Time `json:"created_at"`
	Commision                   float64   `json:"commision"`
	BeforeBalance               float64   `json:"before_balance"`
	AfterBalance                float64   `json:"after_balance"`
	Status                      string    `json:"status"`
}

type GetMobileRechargeOperatorsResponseModel struct {
	OperatorCode int    `json:"operator_code"`
	OperatorName string `json:"operator_name"`
}

type GetMobileRechargeCircleResponseModel struct {
	CircleCode int    `json:"circle_code"`
	CircleName string `json:"circle_name"`
}

type GetMobileRechargePlansRequestModel struct {
	OperatorCode int `json:"operator_code" validate:"required"`
	Circle       int `json:"circle" validate:"required"`
}

type GetMobileRechargePlansResponseModel struct {
	Error    int    `json:"error"`
	Message  string `json:"msg"`
	Status   int    `json:"status"`
	PlanData any    `json:"planData"`
}
