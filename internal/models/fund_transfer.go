package models

import "time"

type CreateFundTransferModel struct {
	FromID  string  `json:"from_id" validate:"required"`
	ToID    string  `json:"to_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required"`
	Remarks string  `json:"remarks" validate:"required"`
}

type GetFundTransferResponseModel struct {
	FundTransferID       int64       `json:"fund_transfer_id"`
	FundTransferFromID   string    `json:"fund_transfer_from_id"`
	FundTransferToID     string    `json:"fund_transfer_on_id"`
	FundTransferFromName string    `json:"fund_transfer_from_name"`
	FundTransferToName   string    `json:"fund_transfer_on_name"`
	Amount               float64   `json:"amount"`
	Remarks              string    `json:"remarks"`
	CreatedAT            time.Time `json:"created_at"`
}

type GetFundTransferFilterRequestModel struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    *string    `json:"status,omitempty"`
	ID        string     `json:"id" validate:"required"`
}
