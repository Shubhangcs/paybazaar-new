package models

import (
	"time"
)

type CreateAdminRequestModel struct {
	AdminName     string `json:"admin_name" validate:"required,min=3,max=100"`
	AdminEmail    string `json:"admin_email" validate:"required,email"`
	AdminPhone    string `json:"admin_phone" validate:"required,phone"`
	AdminPassword string `json:"admin_password" validate:"required,strpwd"`
}

type UpdateAdminRequestModel struct {
	AdminName      *string  `json:"admin_name" validate:"omitempty,min=3,max=100"`
	AdminPhone     *string  `json:"admin_phone" validate:"omitempty,phone"`
	AdminPassword  *string  `json:"admin_password" validate:"omitempty,strpwd"`
	IsAdminBlocked *bool    `json:"is_admin_blocked"`
	WalletBalance  *float64 `json:"admin_wallet_balance" validate:"omitempty,gte=0"`
}

type GetAdminResponseModel struct {
	AdminID            string    `json:"admin_id"`
	AdminName          string    `json:"admin_name"`
	AdminEmail         string    `json:"admin_email"`
	AdminPhone         string    `json:"admin_phone"`
	AdminWalletBalance float64   `json:"admin_wallet_balance"`
	IsAdminBlocked     bool      `json:"is_admin_blocked"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type AdminModel struct {
	AdminID            string
	AdminName          string
	AdminEmail         string
	AdminPhone         string
	AdminPassword      string
	AdminWalletBalance float64
	IsAdminBlocked     bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type LoginAdminModel struct {
	AdminID       string `json:"admin_id" validate:"required"`
	AdminPassword string `json:"admin_password" validate:"required"`
}

type AdminWalletTopupModel struct {
	AdminID string  `json:"admin_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required"`
}
