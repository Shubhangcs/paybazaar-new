package models

import "time"

type WalletTransactionModel struct {
	WalletTransactionID string
	UserID              string
	ReferenceID         string

	CreditAmount  *string
	DebitAmount   *string
	BeforeBalance string
	AfterBalance  string

	TransactionReason string
	Remarks           string

	CreatedAt time.Time
}

type CreateWalletTransactionRequestModel struct {
	UserID      string `json:"user_id" validate:"required"`
	ReferenceID string `json:"reference_id" validate:"required"`

	CreditAmount *string `json:"credit_amount"`
	DebitAmount  *string `json:"debit_amount"`

	BeforeBalance string `json:"before_balance" validate:"required"`
	AfterBalance  string `json:"after_balance" validate:"required"`

	TransactionReason string `json:"transaction_reason" validate:"required"`
	Remarks           string `json:"remarks" validate:"required"`
}

type GetWalletTransactionResponseModel struct {
	WalletTransactionID string `json:"wallet_transaction_id"`
	UserID              string `json:"user_id"`
	ReferenceID         string `json:"reference_id"`

	CreditAmount  *string `json:"credit_amount,omitempty"`
	DebitAmount   *string `json:"debit_amount,omitempty"`
	BeforeBalance string  `json:"before_balance"`
	AfterBalance  string  `json:"after_balance"`

	TransactionReason string    `json:"transaction_reason"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
}
