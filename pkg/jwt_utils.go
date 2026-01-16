package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/levion-studio/paybazaar/internal/models"
)

type JwtUtils struct {
	secret string
	expiry time.Duration
}

type JwtConfig struct {
	SecretKey string
	Expiry    time.Duration
}

func NewJwtUtils(cfg JwtConfig) *JwtUtils {
	return &JwtUtils{
		secret: cfg.SecretKey,
		expiry: cfg.Expiry,
	}
}

func (ju *JwtUtils) GenerateToken(ctx context.Context, req models.AccessTokenClaims) (string, error) {
	claims := models.AccessTokenClaims{
		AdminID:  req.AdminID,
		UserID:   req.UserID,
		UserName: req.UserName,
		UserRole: req.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ju.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ju.secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return tokenString, nil
}

func (ju *JwtUtils) ValidateToken(tokenString string) (*models.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token: wrong signing method")
		}
		return []byte(ju.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
