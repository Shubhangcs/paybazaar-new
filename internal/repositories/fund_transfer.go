package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type FundTransferInterface interface {
	CreateFundTransfer(echo.Context) error
	GetFundTransfersByFromID(echo.Context) ([]models.GetFundTransferResponseModel, error)
	GetFundTransfersByToID(echo.Context) ([]models.GetFundTransferResponseModel, error)
}

type fundTransferRepository struct {
	db *database.Database
}

func NewFundTransferRepository(db *database.Database) *fundTransferRepository {
	return &fundTransferRepository{db: db}
}

func (fr *fundTransferRepository) CreateFundTransfer(c echo.Context) error {
	var req models.CreateFundTransferModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return fr.db.CreateFundTransferQuery(ctx, req)
}

func (fr *fundTransferRepository) GetFundTransfersByFromID(
	c echo.Context,
) ([]models.GetFundTransferResponseModel, error) {

	var req models.GetFundTransferFilterRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return fr.db.GetFundTransfersByFromIDQuery(
		ctx,
		req,
		limit,
		offset,
	)
}

func (fr *fundTransferRepository) GetFundTransfersByToID(
	c echo.Context,
) ([]models.GetFundTransferResponseModel, error) {

	var req models.GetFundTransferFilterRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return fr.db.GetFundTransfersByToIDQuery(
		ctx,
		req,
		limit,
		offset,
	)
}
