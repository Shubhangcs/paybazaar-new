package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type CommisionInterface interface {
	CreateCommision(echo.Context) error
	GetCommisionDetailsByCommisionID(echo.Context) (*models.GetCommisionResponseModel, error)
	GetCommisionsByUserID(echo.Context) ([]models.GetCommisionResponseModel, error)
	GetCommisionByUserIDAndService(echo.Context) (*models.GetCommisionResponseModel, error)
	UpdateCommisionDetails(echo.Context) error
	DeleteCommision(echo.Context) error
	GetAllTDSCommision(echo.Context) ([]models.GetTDSCommisionResponseModel, error)
	GetTDSCommisionByUserID(echo.Context) ([]models.GetTDSCommisionResponseModel, error)
}

type commisionRepository struct {
	db *database.Database
}

func NewCommisionRepository(db *database.Database) *commisionRepository {
	return &commisionRepository{
		db: db,
	}
}

func (cr *commisionRepository) CreateCommision(
	c echo.Context,
) error {
	var req models.CreateCommisionRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return cr.db.CreateCommisionQuery(ctx, req)
}

func (cr *commisionRepository) GetCommisionDetailsByCommisionID(
	c echo.Context,
) (*models.GetCommisionResponseModel, error) {
	commisionID, err := parseInt64Param(c, "commision_id")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return cr.db.GetCommisionDetailsByCommisionIDQuery(ctx, commisionID)
}

func (cr *commisionRepository) GetCommisionsByUserID(
	c echo.Context,
) ([]models.GetCommisionResponseModel, error) {
	userID := c.Param("user_id")
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()
	return cr.db.GetCommisionsByUserIDQuery(ctx, userID)
}

func (cr *commisionRepository) GetCommisionByUserIDAndService(
	c echo.Context,
) (*models.GetCommisionResponseModel, error) {
	userID := c.Param("user_id")
	service := c.Param("service")
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()
	return cr.db.GetCommisionByUserIDAndServiceQuery(ctx, userID, service)
}

func (cr *commisionRepository) UpdateCommisionDetails(
	c echo.Context,
) error {
	var req models.UpdateCommisionRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()
	return cr.db.UpdateCommisionQuery(ctx, req)
}

func (cr *commisionRepository) DeleteCommision(
	c echo.Context,
) error {
	commisionID, err := parseInt64Param(c, "commision_id")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()
	return cr.db.DeleteCommisionQuery(ctx, commisionID)
}

func (cr *commisionRepository) GetAllTDSCommision(c echo.Context) ([]models.GetTDSCommisionResponseModel, error) {
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()
	limit, offset := parsePagination(c)
	return cr.db.GetAllTDSCommisionQuery(ctx, limit, offset)
}

func (cr *commisionRepository) GetTDSCommisionByUserID(c echo.Context) ([]models.GetTDSCommisionResponseModel, error) {
	var userID = c.Param("user_id")
	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()
	limit, offset := parsePagination(c)
	return cr.db.GetTDSCommisionByUserIDQuery(ctx, userID, limit, offset)
}
