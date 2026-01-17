package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/pkg"
)

type AdminInterface interface {
	CreateAdmin(echo.Context) error
	GetAdminByID(echo.Context) (*models.AdminModel, error)
	ListAdmins(echo.Context) ([]models.GetAdminResponseModel, error)
	UpdateAdmin(echo.Context) error
	DeleteAdmin(echo.Context) error
	GetAdminsForDropdown(echo.Context) ([]models.DropdownModel, error)
	LoginAdmin(echo.Context) (string, error)
	AdminWalletTopup(echo.Context) (float64, error)
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

func (ar *adminRepository) UpdateAdmin(c echo.Context) error {
	adminID := c.Param("admin_id")

	var req models.UpdateAdminRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return ar.db.UpdateAdminQuery(ctx, adminID, req)
}

func (ar *adminRepository) GetAdminByID(
	c echo.Context,
) (*models.AdminModel, error) {

	adminID := c.Param("admin_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return ar.db.GetAdminByIDQuery(ctx, adminID)
}

func (ar *adminRepository) ListAdmins(
	c echo.Context,
) ([]models.GetAdminResponseModel, error) {

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return ar.db.ListAdminsQuery(ctx, limit, offset)
}

func (ar *adminRepository) DeleteAdmin(c echo.Context) error {
	adminID := c.Param("admin_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return ar.db.DeleteAdminQuery(ctx, adminID)
}

func (ar *adminRepository) GetAdminsForDropdown(
	c echo.Context,
) ([]models.DropdownModel, error) {

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return ar.db.GetAllAdminsForDropdownQuery(ctx)
}

func (ar *adminRepository) LoginAdmin(c echo.Context) (string, error) {
	var req models.LoginAdminModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	res, err := ar.db.GetAdminByIDQuery(ctx, req.AdminID)
	if err != nil {
		return "", err
	}
	if res.AdminPassword != req.AdminPassword {
		return "", fmt.Errorf("incorrect password")
	}
	return ar.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		UserID:   res.AdminID,
		UserName: res.AdminName,
		UserRole: "admin",
	})
}

func (ar *adminRepository) AdminWalletTopup(c echo.Context) (float64, error) {
	var req models.AdminWalletTopupModel
	if err := bindAndValidate(c, &req); err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return ar.db.AdminWalletTopupQuery(ctx, req)
}
