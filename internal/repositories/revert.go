package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type RevertInterface interface {
	CreateRevert(echo.Context) error
	GetRevertsByFromID(echo.Context) ([]models.GetRevertTransactionResponseModel, error)
	GetRevertsByOnID(echo.Context) ([]models.GetRevertTransactionResponseModel, error)
}

type revertRepository struct {
	db *database.Database
}

func NewRevertRepository(db *database.Database) *revertRepository {
	return &revertRepository{
		db,
	}
}

func (rr *revertRepository) CreateRevert(c echo.Context) error {
	var req models.CreateRevertRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return rr.db.CreateRevertQuery(ctx, req)
}

func (rr *revertRepository) GetRevertsByFromID(
	c echo.Context,
) ([]models.GetRevertTransactionResponseModel, error) {

	var req models.GetRevertTransactionFilterRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	data, err := rr.db.GetRevertTransactionsByFromIDQuery(
		ctx,
		req,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []models.GetRevertTransactionResponseModel{}
	}

	return data, nil
}

func (rr *revertRepository) GetRevertsByOnID(
	c echo.Context,
) ([]models.GetRevertTransactionResponseModel, error) {

	var req models.GetRevertTransactionFilterRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	data, err := rr.db.GetRevertTransactionsByOnIDQuery(
		ctx,
		req,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data = []models.GetRevertTransactionResponseModel{}
	}

	return data, nil
}
