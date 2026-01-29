package models

import "time"

type CreateMobileRechargeRequestModel struct {
	RetailerID       string  `json:"retailer_id"`
	MobileNumber     string  `json:"mobile_number" validate:"required,phone"`
	OperatorCode     int64   `json:"operator_code" validate:"required"`
	OperatorName     string  `json:"operator_name" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	CircleCode       int64   `json:"circle_code" validate:"required"`
	CircleName       string  `json:"circle_name" validate:"required"`
	RechargeType     string  `json:"recharge_type,omitempty"`
	PartnerRequestID string  `json:"partner_request_id,omitempty"`
	Commision        float64 `json:"commision"`
	Status           string  `json:"status"`
}

type GetMobileRechargeHistoryResponseModel struct {
	RetailerID                  string    `json:"retailer_id"`
	MobileRechargeTransactionID int64     `json:"mobile_recharge_transaction_id"`
	MobileNumber                string    `json:"mobile_number"`
	OperatorCode                int64     `json:"operator_code"`
	OperatorName                string    `json:"operator_name"`
	Amount                      float64   `json:"amount"`
	CircleCode                  int64     `json:"circle_code"`
	CircleName                  string    `json:"circle_name"`
	RechargeType                string    `json:"recharge_type"`
	PartnerRequestID            string    `json:"partner_request_id"`
	CreatedAt                   time.Time `json:"created_at"`
	Commision                   float64   `json:"commision"`
	Status                      string    `json:"status"`
}

type GetMobileRechargeOperatorsResponseModel struct {
	OperatorCode int64  `json:"operator_code"`
	OperatorName string `json:"operator_name"`
}

type GetMobileRechargeCircleResponseModel struct {
	CircleCode int64  `json:"circle_code"`
	CircleName string `json:"circle_name"`
}

type GetMobileRechargePlansRequestModel struct {
	OperatorCode int64 `json:"operator_code" validate:"required"`
	Circle       int64 `json:"circle" validate:"required"`
}

type GetMobileRechargePlansResponseModel struct {
	Error    int64         `json:"error"`
	Message  string        `json:"msg"`
	Status   int64         `json:"status"`
	PlanData PlanDataModel `json:"planData"`
}

type PlanDataModel struct {
	CircleCode int64  `json:"circle_id"`
	ID         string `json:"_id"`
	Plans      []any  `json:"plan"`
}
