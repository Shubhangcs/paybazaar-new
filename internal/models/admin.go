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

type UpdateAdminDetailsRequestModel struct {
	AdminID    string  `json:"admin_id" validate:"required"`
	AdminName  *string `json:"admin_name" validate:"omitempty,min=3,max=100"`
	AdminPhone *string `json:"admin_phone" validate:"omitempty,phone"`
	AdminEmail *string `json:"admin_email" validate:"omitempty,email"`
}

type UpdateAdminPasswordRequestModel struct {
	AdminID     string `json:"admin_id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required,strpwd"`
	NewPassword string `json:"new_password" validate:"required,strpwd"`
}

type UpdateAdminWalletRequestModel struct {
	AdminID string  `json:"admin_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,min=1"`
}

type UpdateAdminBlockStatusRequestModel struct {
	AdminID     string `json:"admin_id" validate:"required"`
	BlockStatus bool   `json:"block_status"`
}

type GetCompleteAdminDetailsResponseModel struct {
	AdminID            string    `json:"admin_id"`
	AdminName          string    `json:"admin_name"`
	AdminEmail         string    `json:"admin_email"`
	AdminPhone         string    `json:"admin_phone"`
	AdminWalletBalance float64   `json:"admin_wallet_balance"`
	IsAdminBlocked     bool      `json:"is_admin_blocked"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type GetAdminDetailsForLoginModel struct {
	AdminID        string `json:"admin_id"`
	AdminName      string `json:"admin_name"`
	AdminPassword  string `json:"admin_password"`
	IsAdminBlocked bool   `json:"is_admin_blocked"`
}

type GetAdminDetailsForDropdownModel struct {
	AdminName string `json:"admin_name"`
	AdminID   string `json:"admin_id"`
}

type AdminLoginRequestModel struct {
	AdminID       string `json:"admin_id" validate:"required"`
	AdminPassword string `json:"admin_password" validate:"required,strpwd"`
}

type RechargeKitWalletBalanceResponseModel struct {
	Error           int     `json:"error"`
	Message         string  `json:"msg"`
	WalletAmount    float64 `json:"wallet_amount"`
	DMRWalletAmount float64 `json:"dmr_wallet_amount"`
}
