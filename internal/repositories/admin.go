package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/pkg"
)

type AdminInterface interface {
	CreateAdmin(echo.Context) error
	GetAdminDetailsByAdminID(echo.Context) (*models.GetCompleteAdminDetailsResponseModel, error)
	GetAllAdmins(echo.Context) ([]models.GetCompleteAdminDetailsResponseModel, error)
	GetAdminsForDropdown(echo.Context) ([]models.GetAdminDetailsForDropdownModel, error)
	UpdateAdminDetails(echo.Context) error
	UpdateAdminPassword(echo.Context) error
	UpdateAdminWallet(echo.Context) error
	UpdateAdminBlockStatus(echo.Context) error
	DeleteAdmin(echo.Context) error
	AdminLogin(echo.Context) (string, error)
	GetRechargeKitWalletBalance(c echo.Context) (*models.RechargeKitWalletBalanceResponseModel, error)
}

type adminRepository struct {
	db       *database.Database
	jwtUtils *pkg.JwtUtils
}

func NewAdminRepository(db *database.Database, jwtUtils *pkg.JwtUtils) *adminRepository {
	return &adminRepository{
		db:       db,
		jwtUtils: jwtUtils,
	}
}

func (ar *adminRepository) CreateAdmin(c echo.Context) error {
	var req models.CreateAdminRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return ar.db.CreateAdminQuery(ctx, req)
}

func (ar *adminRepository) GetAdminDetailsByAdminID(c echo.Context) (*models.GetCompleteAdminDetailsResponseModel, error) {
	var adminID = c.Param("admin_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.GetAdminDetailsByAdminID(ctx, adminID)
}

func (ar *adminRepository) GetAllAdmins(c echo.Context) ([]models.GetCompleteAdminDetailsResponseModel, error) {
	limit, offset := parsePagination(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.GetAllAdminsQuery(ctx, offset, limit)
}

func (ar *adminRepository) GetAdminsForDropdown(c echo.Context) ([]models.GetAdminDetailsForDropdownModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.GetAllAdminsForDropdownQuery(ctx)
}

func (ar *adminRepository) UpdateAdminDetails(c echo.Context) error {
	var req models.UpdateAdminDetailsRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.UpdateAdminDetailsQuery(ctx, req)
}

func (ar *adminRepository) UpdateAdminPassword(c echo.Context) error {
	var req models.UpdateAdminPasswordRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.UpdateAdminPasswordQuery(ctx, req)
}

func (ar *adminRepository) UpdateAdminWallet(c echo.Context) error {
	var req models.UpdateAdminWalletRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.UpdateAdminWalletQuery(ctx, req)
}

func (ar *adminRepository) UpdateAdminBlockStatus(c echo.Context) error {
	var req models.UpdateAdminBlockStatusRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.UpdateAdminBlockStatusQuery(ctx, req)
}

func (ar *adminRepository) DeleteAdmin(c echo.Context) error {
	var adminID = c.Param("admin_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.DeleteAdminQuery(ctx, adminID)
}

func (ar *adminRepository) AdminLogin(c echo.Context) (string, error) {
	var req models.AdminLoginRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	details, err := ar.db.GetAdminDetailsForLogin(ctx, req.AdminID)
	if err != nil {
		return "", err
	}

	if details.AdminPassword != req.AdminPassword {
		return "", fmt.Errorf("incorrect password")
	}

	if details.IsAdminBlocked {
		return "", fmt.Errorf("admin is blocked")
	}

	token, err := ar.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		AdminID:  details.AdminID,
		UserName: details.AdminName,
		UserRole: "admin",
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (ar *adminRepository) GetRechargeKitWalletBalance(c echo.Context) (*models.RechargeKitWalletBalanceResponseModel, error) {
	apiUrl := `https://v2a.rechargkit.biz/recharge/balanceCheck`

	apiRequest, err := http.NewRequest(
		http.MethodGet,
		apiUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse models.RechargeKitWalletBalanceResponseModel

	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
