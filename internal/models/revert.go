package models

import "time"

type CreateRevertRequest struct {
	FromID  string  `json:"from_id" validate:"required"`
	OnID    string  `json:"on_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required"`
	Remarks string  `json:"remarks"`
}

type GetRevertTransactionResponseModel struct {
	RevertTransactionID   int       `json:"revert_transaction_id"`
	RevertFromID          string    `json:"revert_from_id"`
	RevertOnID            string    `json:"revert_on_id"`
	RevertFromName        string    `json:"revert_from_name"`
	RevertOnName          string    `json:"revert_on_name"`
	RevertOnBusinessName  string    `json:"revert_on_business_name"`
	Amount                float64   `json:"amount"`
	Remarks               string    `json:"remarks"`
	CreatedAT             time.Time `json:"created_at"`
}


type GetRevertTransactionFilterRequestModel struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    *string    `json:"status,omitempty"`
	ID        string     `json:"id" validate:"required"`
}
