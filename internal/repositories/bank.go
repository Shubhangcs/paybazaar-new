package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type BankInterface interface {
	CreateBank(echo.Context) error
	GetBankDetailsByBankID(echo.Context) (*models.GetBankDetailsResponseModel, error)
	GetAllBanks(echo.Context) ([]models.GetBankDetailsResponseModel, error)
	UpdateBankDetails(echo.Context) error
	DeleteBank(echo.Context) error
	CreateAdminBank(echo.Context) error
	GetAdminBankDetailsByAdminBankID(echo.Context) (*models.GetAdminBankDetailsResponseModel, error)
	GetAdminBanksByAdminID(echo.Context) ([]models.GetAdminBankDetailsResponseModel, error)
	UpdateAdminBankDetails(echo.Context) error
	DeleteAdminBank(echo.Context) error
}

type bankRepository struct {
	db *database.Database
}

func NewBankRepository(db *database.Database) *bankRepository {
	return &bankRepository{
		db: db,
	}
}

func (br *bankRepository) CreateBank(c echo.Context) error {
	var req models.CreateBankRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.CreateBankQuery(ctx, req)
}

func (br *bankRepository) GetBankDetailsByBankID(
	c echo.Context,
) (*models.GetBankDetailsResponseModel, error) {
	bankID, err := parseInt64Param(c, "bank_id")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.GetBankDetailsByBankIDQuery(ctx, bankID)
}

func (br *bankRepository) GetAllBanks(
	c echo.Context,
) ([]models.GetBankDetailsResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()
	return br.db.GetAllBanksQuery(ctx)
}

func (br *bankRepository) UpdateBankDetails(c echo.Context) error {
	var req models.UpdateBankDetailsRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.UpdateBankQuery(ctx, req)
}

func (br *bankRepository) DeleteBank(c echo.Context) error {
	bankID, err := parseInt64Param(c, "bank_id")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()
	return br.db.DeleteBankQuery(ctx, bankID)
}

func (br *bankRepository) CreateAdminBank(
	c echo.Context,
) error {

	var req models.CreateAdminBankRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.CreateAdminBankQuery(ctx, req)
}

func (br *bankRepository) GetAdminBankDetailsByAdminBankID(
	c echo.Context,
) (*models.GetAdminBankDetailsResponseModel, error) {
	adminBankID, err := parseInt64Param(c, "admin_bank_id")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.GetAdminBankDetailsByAdminBankIDQuery(ctx, adminBankID)
}

func (br *bankRepository) GetAdminBanksByAdminID(
	c echo.Context,
) ([]models.GetAdminBankDetailsResponseModel, error) {
	adminID := c.Param("admin_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()
	return br.db.GetAdminBanksByAdminIDQuery(ctx, adminID)
}

func (br *bankRepository) UpdateAdminBankDetails(c echo.Context) error {
	var req models.UpdateAdminBankDetailsRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.UpdateAdminBankQuery(ctx, req)
}

func (br *bankRepository) DeleteAdminBank(c echo.Context) error {
	adminBankID, err := parseInt64Param(c, "admin_bank_id")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.DeleteAdminBankQuery(ctx, adminBankID)
}
