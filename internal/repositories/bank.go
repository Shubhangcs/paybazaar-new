package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type BankInterface interface {
	// Banks
	CreateBank(echo.Context) (int64, error)
	GetBankByID(echo.Context) (*models.GetBankModel, error)
	GetAllBanks(echo.Context) ([]models.GetBankModel, error)
	UpdateBank(echo.Context) error
	DeleteBank(echo.Context) error

	// Admin Banks
	CreateAdminBank(echo.Context) (int64, error)
	GetAdminBankByID(echo.Context) (*models.GetAdminBankModel, error)
	GetAdminBanksByAdminID(echo.Context) ([]models.GetAdminBankModel, error)
	UpdateAdminBank(echo.Context) error
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

func (br *bankRepository) CreateBank(c echo.Context) (int64, error) {
	var req models.CreateBankModel
	if err := bindAndValidate(c, &req); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.CreateBankQuery(ctx, req)
}

func (br *bankRepository) GetBankByID(
	c echo.Context,
) (*models.GetBankModel, error) {

	bankID, err := parseInt64Param(c, "bank_id")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.GetBankByIDQuery(ctx, bankID)
}

func (br *bankRepository) GetAllBanks(
	c echo.Context,
) ([]models.GetBankModel, error) {

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	data, err := br.db.GetAllBanksQuery(ctx)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []models.GetBankModel{}
	}

	return data, nil
}

func (br *bankRepository) UpdateBank(c echo.Context) error {
	bankID, err := parseInt64Param(c, "bank_id")
	if err != nil {
		return err
	}

	var req models.UpdateBankModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.UpdateBankQuery(ctx, bankID, req)
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
) (int64, error) {

	var req models.CreateAdminBankModel
	if err := bindAndValidate(c, &req); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.CreateAdminBankQuery(ctx, req)
}

func (br *bankRepository) GetAdminBankByID(
	c echo.Context,
) (*models.GetAdminBankModel, error) {

	adminBankID, err := parseInt64Param(c, "admin_bank_id")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.GetAdminBankByIDQuery(ctx, adminBankID)
}

func (br *bankRepository) GetAdminBanksByAdminID(
	c echo.Context,
) ([]models.GetAdminBankModel, error) {

	adminID := c.Param("admin_id")
	if adminID == "" {
		return nil, echo.NewHTTPError(400, "admin_id is required")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	data, err := br.db.GetAdminBanksByAdminIDQuery(ctx, adminID)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []models.GetAdminBankModel{}
	}

	return data, nil
}

func (br *bankRepository) UpdateAdminBank(c echo.Context) error {
	adminBankID, err := parseInt64Param(c, "admin_bank_id")
	if err != nil {
		return err
	}

	var req models.UpdateAdminBankModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return br.db.UpdateAdminBankQuery(ctx, adminBankID, req)
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
