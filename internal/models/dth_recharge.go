package models

type CreateDTHRechargeRequestModel struct {
	CustomerID       string  `json:"customer_id" validate:"required"`
	OperatorName     string  `json:"operator_name"`
	OperatorCode     int   `json:"operator_code" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	PartnerRequestID string  `json:"partner_request_id"`
	Status           string  `json:"status"`
}
