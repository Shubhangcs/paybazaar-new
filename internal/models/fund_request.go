package models

import "time"

type FundRequestModel struct {
	FundRequestID int64
	RequesterID   string
	RequestToID   string
	Amount        string
	BankName      string
	RequestDate   string
	UTRNumber     string
	RequestStatus string
	Remarks       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateFundRequestModel struct {
	RequesterID string `json:"requester_id" validate:"required"`
	RequestToID string `json:"request_to_id" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	BankName    string `json:"bank_name" validate:"required"`
	RequestDate string `json:"request_date" validate:"required"`
	UTRNumber   string `json:"utr_number" validate:"required"`
	Remarks     string `json:"remarks" validate:"required"`
}

type GetFundRequestFilterRequestModel struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    *string    `json:"status,omitempty"`
	ID        string     `json:"id" validate:"required"`
}

type GetFundRequestResponseModel struct {
	FundRequestID int64     `json:"fund_request_id"`
	RequesterID   string    `json:"requester_id"`
	RequestToID   string    `json:"request_to_id"`
	Amount        string    `json:"amount"`
	BankName      string    `json:"bank_name"`
	RequestDate   string    `json:"request_date"`
	UTRNumber     string    `json:"utr_number"`
	RequestStatus string    `json:"request_status"`
	Remarks       string    `json:"remarks"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
