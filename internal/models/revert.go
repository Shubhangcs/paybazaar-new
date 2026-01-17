package models

type CreateRevertRequest struct {
	FromID  string  `json:"from_id" validate:"required"`
	OnID    string  `json:"on_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required"`
	Remarks string  `json:"remarks"`
}
