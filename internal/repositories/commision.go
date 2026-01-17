package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type CommisionInterface interface {
	CreateCommision(echo.Context) (int64, error)
	GetCommisionByID(echo.Context) (*models.GetCommisionModel, error)
	GetCommisionByUserID(echo.Context) (*models.GetCommisionModel, error)
	GetAllCommisions(echo.Context) ([]models.GetCommisionModel, error)
	UpdateCommision(echo.Context) error
	DeleteCommision(echo.Context) error
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
) (int64, error) {

	var req models.CreateCommisionModel
	if err := bindAndValidate(c, &req); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return cr.db.CreateCommisionQuery(ctx, req)
}

func (cr *commisionRepository) GetCommisionByID(
	c echo.Context,
) (*models.GetCommisionModel, error) {

	commisionID, err := parseInt64Param(c, "commision_id")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return cr.db.GetCommisionByIDQuery(ctx, commisionID)
}

func (cr *commisionRepository) GetCommisionByUserID(
	c echo.Context,
) (*models.GetCommisionModel, error) {

	userID := c.Param("user_id")
	if userID == "" {
		return nil, echo.NewHTTPError(400, "user_id is required")
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return cr.db.GetCommisionByUserIDQuery(ctx, userID)
}

func (cr *commisionRepository) GetAllCommisions(
	c echo.Context,
) ([]models.GetCommisionModel, error) {

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	data, err := cr.db.GetAllCommisionsQuery(
		ctx,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []models.GetCommisionModel{}
	}

	return data, nil
}

func (cr *commisionRepository) UpdateCommision(
	c echo.Context,
) error {

	commisionID, err := parseInt64Param(c, "commision_id")
	if err != nil {
		return err
	}

	var req models.UpdateCommisionModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return cr.db.UpdateCommisionQuery(ctx, commisionID, req)
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
