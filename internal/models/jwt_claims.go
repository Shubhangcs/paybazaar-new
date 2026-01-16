package models

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	AdminID  string `json:"admin_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}
