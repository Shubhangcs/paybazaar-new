package models

import "time"

type CreateTransactionLimitRequestModel struct {
	RetailerID  string  `json:"retailer_id" validate:"required"`
	LimitAmount float64 `json:"limit_amount" validate:"required"`
	Service     string  `json:"service" validate:"required"`
}

type UpdateTransactionLimitRequestModel struct {
	LimitID     int     `json:"limit_id" validate:"required"`
	LimitAmount float64 `json:"limit_amount"`
	Service     string  `json:"service"`
}

type GetLimitRequestModel struct {
	RetailerID string `json:"retailer_id"`
	Service    string `json:"service"`
}

type GetLimitResponseModel struct {
	LimitID     int       `json:"limit_id"`
	RetailerID  string    `json:"retailer_id"`
	LimitAmount float64   `json:"limit_amount"`
	Service     string    `json:"service"`
	CreatedAT   time.Time `json:"created_at"`
	UpdatedAT   time.Time `json:"updated_at"`
}
