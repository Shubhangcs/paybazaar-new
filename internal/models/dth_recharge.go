package models

type CreateDTHRechargeRequestModel struct {
	CustomerID       string  `json:"customer_id"`
	OperatorCode     string  `json:"operator_code"`
	Amount           float64 `json:"amount"`
	PartnerRequestID string  `json:"partner_request_id"`
}
