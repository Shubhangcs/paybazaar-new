package repositories

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type LimitInterface interface {
	CreateLimit(echo.Context) error
	UpdateLimit(echo.Context) error
	DeleteLimit(echo.Context) error
	GetAllLimits(c echo.Context) ([]models.GetLimitResponseModel, error)
	GetLimitByRetailerIDAndService(echo.Context) (*models.GetLimitResponseModel, error)
}

type limitRepository struct {
	db *database.Database
}

func NewLimitRepository(db *database.Database) *limitRepository {
	return &limitRepository{
		db,
	}
}

func (lr *limitRepository) CreateLimit(c echo.Context) error {
	var req models.CreateTransactionLimitRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	return lr.db.CreateLimitQuery(ctx, req)
}

func (lr *limitRepository) UpdateLimit(c echo.Context) error {
	var req models.UpdateTransactionLimitRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	return lr.db.UpdateLimitQuery(ctx, req)
}

func (lr *limitRepository) DeleteLimit(c echo.Context) error {
	var limitId = c.Param("limit_id")
	limId, err := strconv.ParseInt(limitId, 10, 64)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	return lr.db.DeleteLimitQuery(ctx, int(limId))
}

func (lr *limitRepository) GetAllLimits(c echo.Context) ([]models.GetLimitResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	return lr.db.GetAllLimitsQuery(ctx)
}

func (lr *limitRepository) GetLimitByRetailerIDAndService(c echo.Context) (*models.GetLimitResponseModel, error) {
	var retailerId = c.Param("retailer_id")
	var service = c.Param("service")
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	return lr.db.GetLimitByRetailerIDServiceQuery(ctx, retailerId, service)
}
