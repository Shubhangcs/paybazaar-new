package models

import "time"

type CreateDTHRechargeRequestModel struct {
	RetailerID       string  `json:"retailer_id" validate:"required"`
	CustomerID       string  `json:"customer_id" validate:"required"`
	OperatorName     string  `json:"operator_name" validate:"required"`
	OperatorCode     int     `json:"operator_code" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	PartnerRequestID string  `json:"partner_request_id"`
	Status           string  `json:"status"`
	Commision        float64 `json:"commision"`
}

type GetDTHRechargeHistoryResponseModel struct {
	DTHTransactionID int       `json:"dth_transaction_id"`
	RetailerID       string    `json:"retailer_id"`
	CustomerID       string    `json:"customer_id"`
	OperatorName     string    `json:"operator_name"`
	OperatorCode     int       `json:"operator_code"`
	Amount           float64   `json:"amount"`
	PartnerRequestID string    `json:"partner_request_id"`
	Status           string    `json:"status"`
	BeforeBalance    float64   `json:"before_balance"`
	AfterBalance     float64   `json:"after_balance"`
	CreatedAt        time.Time `json:"created_at"`
	Commision        float64   `json:"commision"`
}

type GetDTHOperatorsResponseModel struct {
	OperatorCode string `json:"operator_code"`
	OperatorName string `json:"operator_name"`
}
