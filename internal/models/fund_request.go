package models

import "time"

type FundRequestModel struct {
	FundRequestID  int64
	RequesterID    string
	RequestToID    string
	Amount         float64
	BankName       string
	RequestDate    time.Time
	UTRNumber      string
	AdvancePayment bool
	RequestStatus  string
	Remarks        string
	RejectRemarks  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateFundRequestModel struct {
	RequesterID    string    `json:"requester_id" validate:"required"`
	RequestToID    string    `json:"request_to_id" validate:"required"`
	Amount         float64   `json:"amount" validate:"required,gt=0"`
	BankName       string    `json:"bank_name" validate:"required"`
	RequestDate    time.Time `json:"request_date" validate:"required"`
	UTRNumber      string    `json:"utr_number"`
	AdvancePayment bool      `json:"advance_payment"`
	Remarks        string    `json:"remarks" validate:"required"`
}

type RejectFundRequestModel struct {
	RejectRemarks string `json:"reject_remarks" validate:"required"`
}

type GetFundRequestFilterRequestModel struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    *string    `json:"status,omitempty"`
	ID        string     `json:"id" validate:"required"`
}

type GetFundRequestResponseModel struct {
	FundRequestID  int64     `json:"fund_request_id"`
	RequesterID    string    `json:"requester_id"`
	RequestToID    string    `json:"request_to_id"`
	BusinessName   string    `json:"business_name"`
	Amount         float64   `json:"amount"`
	BankName       string    `json:"bank_name"`
	RequestDate    time.Time `json:"request_date"`
	UTRNumber      *string   `json:"utr_number"`
	AdvancePayment bool      `json:"advance_payment"`
	RequestStatus  string    `json:"request_status"`
	Remarks        string    `json:"remarks"`
	RejectRemarks  string    `json:"reject_remarks"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
