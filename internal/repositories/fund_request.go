package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type FundRequestInterface interface {
	CreateFundRequest(echo.Context) (int64, error)
	GetFundRequestByID(echo.Context) (*models.GetFundRequestResponseModel, error)
	GetAllFundRequests(echo.Context) ([]models.GetFundRequestResponseModel, error)
	GetFundRequestsByRequesterID(echo.Context) ([]models.GetFundRequestResponseModel, error)
	GetFundRequestsByRequestToID(echo.Context) ([]models.GetFundRequestResponseModel, error)
	AcceptFundRequest(echo.Context) error
	RejectFundRequest(echo.Context) error
}

type fundRequestRepository struct {
	db *database.Database
}

func NewFundRequestRepository(db *database.Database) *fundRequestRepository {
	return &fundRequestRepository{db: db}
}

func (fr *fundRequestRepository) CreateFundRequest(
	c echo.Context,
) (int64, error) {

	var req models.CreateFundRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	return fr.db.CreateFundRequestQuery(ctx, req)
}

func (fr *fundRequestRepository) GetFundRequestByID(
	c echo.Context,
) (*models.GetFundRequestResponseModel, error) {

	fundRequestID, err := parseInt64Param(c, "fund_request_id")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	return fr.db.GetFundRequestQuery(ctx, fundRequestID)
}

func (fr *fundRequestRepository) GetAllFundRequests(
	c echo.Context,
) ([]models.GetFundRequestResponseModel, error) {

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return fr.db.GetAllFundRequestsQuery(ctx, limit, offset)
}

func (fr *fundRequestRepository) GetFundRequestsByRequesterID(
	c echo.Context,
) ([]models.GetFundRequestResponseModel, error) {

	var req models.GetFundRequestFilterRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return fr.db.GetFundRequestsByRequesterIDQuery(ctx, req, limit, offset)
}

func (fr *fundRequestRepository) GetFundRequestsByRequestToID(
	c echo.Context,
) ([]models.GetFundRequestResponseModel, error) {

	var req models.GetFundRequestFilterRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return fr.db.GetFundRequestsByRequestToIDQuery(ctx, req, limit, offset)
}

func (fr *fundRequestRepository) AcceptFundRequest(
	c echo.Context,
) error {

	fundRequestID, err := parseInt64Param(c, "fund_request_id")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return fr.db.AcceptFundRequestQuery(ctx, fundRequestID)
}

func (fr *fundRequestRepository) RejectFundRequest(
	c echo.Context,
) error {

	fundRequestID, err := parseInt64Param(c, "fund_request_id")
	if err != nil {
		return err
	}

	var req struct {
		RejectRemarks string `json:"reject_remarks" validate:"required,min=3"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	return fr.db.RejectFundRequestQuery(
		ctx,
		fundRequestID,
		req.RejectRemarks,
	)
}
