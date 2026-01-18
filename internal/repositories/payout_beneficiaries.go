package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type RetailerBeneficiaryInterface interface {
	CreateRetailerBeneficiary(echo.Context) error
	GetRetailerBeneficiaryByID(echo.Context) (*models.GetRetailerBeneficiaryResponseModel, error)
	ListRetailerBeneficiaries(echo.Context) ([]models.GetRetailerBeneficiaryResponseModel, error)
	UpdateRetailerBeneficiary(echo.Context) error
	UpdateRetailerBeneficiaryVerification(echo.Context) error
	DeleteRetailerBeneficiary(echo.Context) error
}

type retailerBeneficiaryRepository struct {
	db *database.Database
}

func NewRetailerBeneficiaryRepository(
	db *database.Database,
) *retailerBeneficiaryRepository {
	return &retailerBeneficiaryRepository{
		db: db,
	}
}

func (rb *retailerBeneficiaryRepository) CreateRetailerBeneficiary(
	c echo.Context,
) error {

	var req models.CreateRetailerBeneficiaryModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	_, err := rb.db.CreateRetailerBeneficiaryQuery(ctx, req)
	return err
}

func (rb *retailerBeneficiaryRepository) GetRetailerBeneficiaryByID(
	c echo.Context,
) (*models.GetRetailerBeneficiaryResponseModel, error) {

	beneficiaryID, err := parseInt64Param(c, "beneficiary_id")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.GetRetailerBeneficiaryByIDQuery(ctx, beneficiaryID)
}

func (rb *retailerBeneficiaryRepository) ListRetailerBeneficiaries(
	c echo.Context,
) ([]models.GetRetailerBeneficiaryResponseModel, error) {

	retailerID := c.QueryParam("retailer_id")
	if retailerID == "" {
		return nil, echo.NewHTTPError(400, "retailer_id is required")
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.GetRetailerBeneficiariesByRetailerIDQuery(
		ctx,
		retailerID,
		limit,
		offset,
	)
}

func (rb *retailerBeneficiaryRepository) UpdateRetailerBeneficiary(
	c echo.Context,
) error {

	beneficiaryID, err := parseInt64Param(c, "beneficiary_id")
	if err != nil {
		return err
	}

	var req models.UpdateRetailerBeneficiaryModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.UpdateRetailerBeneficiaryQuery(ctx, beneficiaryID, req)
}

func (rb *retailerBeneficiaryRepository) UpdateRetailerBeneficiaryVerification(
	c echo.Context,
) error {

	beneficiaryID, err := parseInt64Param(c, "beneficiary_id")
	if err != nil {
		return err
	}

	var req struct {
		IsVerified bool `json:"is_verified" validate:"required"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.UpdateRetailerBeneficiaryVerificationQuery(
		ctx,
		beneficiaryID,
		req.IsVerified,
	)
}

func (rb *retailerBeneficiaryRepository) DeleteRetailerBeneficiary(
	c echo.Context,
) error {

	beneficiaryID, err := parseInt64Param(c, "beneficiary_id")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.DeleteRetailerBeneficiaryQuery(ctx, beneficiaryID)
}
