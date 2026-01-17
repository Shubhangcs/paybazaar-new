package models

import "time"

type WalletTransactionModel struct {
	WalletTransactionID string
	UserID              string
	ReferenceID         string

	CreditAmount  *float64
	DebitAmount   *float64
	BeforeBalance float64
	AfterBalance  float64

	TransactionReason string
	Remarks           string

	CreatedAt time.Time
}

type CreateWalletTransactionRequestModel struct {
	UserID      string `json:"user_id" validate:"required"`
	ReferenceID string `json:"reference_id" validate:"required"`

	CreditAmount *float64 `json:"credit_amount"`
	DebitAmount  *float64 `json:"debit_amount"`

	BeforeBalance float64 `json:"before_balance" validate:"required"`
	AfterBalance  float64 `json:"after_balance" validate:"required"`

	TransactionReason string `json:"transaction_reason" validate:"required"`
	Remarks           string `json:"remarks" validate:"required"`
}

type GetWalletTransactionResponseModel struct {
	WalletTransactionID string `json:"wallet_transaction_id"`
	UserID              string `json:"user_id"`
	ReferenceID         string `json:"reference_id"`

	CreditAmount  *float64 `json:"credit_amount,omitempty"`
	DebitAmount   *float64 `json:"debit_amount,omitempty"`
	BeforeBalance float64  `json:"before_balance"`
	AfterBalance  float64  `json:"after_balance"`

	TransactionReason string    `json:"transaction_reason"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
}
