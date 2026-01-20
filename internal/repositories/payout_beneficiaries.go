package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type PayoutBeneficiaryInterface interface {
	CreatePayoutBeneficiary(echo.Context) error
	GetPayoutBeneficiaryByID(echo.Context) (*models.GetPayoutBeneficiaryResponseModel, error)
	ListPayoutBeneficiaries(echo.Context) ([]models.GetPayoutBeneficiaryResponseModel, error)
	UpdatePayoutBeneficiary(echo.Context) error
	UpdatePayoutBeneficiaryVerification(echo.Context) error
	DeletePayoutBeneficiary(echo.Context) error
	GetPayoutBeneficiariesByMobileNumber(echo.Context) ([]models.GetPayoutBeneficiaryResponseModel, error)
}

type payoutBeneficiaryRepository struct {
	db *database.Database
}

func NewPayoutBeneficiaryRepository(
	db *database.Database,
) *payoutBeneficiaryRepository {
	return &payoutBeneficiaryRepository{
		db: db,
	}
}

func (rb *payoutBeneficiaryRepository) CreatePayoutBeneficiary(
	c echo.Context,
) error {

	var req models.CreatePayoutBeneficiaryModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	_, err := rb.db.CreatePayoutBeneficiaryQuery(ctx, req)
	return err
}

func (rb *payoutBeneficiaryRepository) GetPayoutBeneficiaryByID(
	c echo.Context,
) (*models.GetPayoutBeneficiaryResponseModel, error) {

	beneficiaryID, err := parseInt64Param(c, "beneficiary_id")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.GetPayoutBeneficiaryByIDQuery(ctx, beneficiaryID)
}

func (rb *payoutBeneficiaryRepository) ListPayoutBeneficiaries(
	c echo.Context,
) ([]models.GetPayoutBeneficiaryResponseModel, error) {

	payoutID := c.Param("payout_id")
	if payoutID == "" {
		return nil, echo.NewHTTPError(400, "payout_id is required")
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.GetPayoutBeneficiariesByRetailerIDQuery(
		ctx,
		payoutID,
		limit,
		offset,
	)
}

func (rb *payoutBeneficiaryRepository) UpdatePayoutBeneficiary(
	c echo.Context,
) error {

	beneficiaryID, err := parseInt64Param(c, "beneficiary_id")
	if err != nil {
		return err
	}

	var req models.UpdatePayoutBeneficiaryModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.UpdatePayoutBeneficiaryQuery(ctx, beneficiaryID, req)
}

func (rb *payoutBeneficiaryRepository) UpdatePayoutBeneficiaryVerification(
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

	return rb.db.UpdatePayoutBeneficiaryVerificationQuery(
		ctx,
		beneficiaryID,
		req.IsVerified,
	)
}

func (rb *payoutBeneficiaryRepository) DeletePayoutBeneficiary(
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

	return rb.db.DeletePayoutBeneficiaryQuery(ctx, beneficiaryID)
}

func (rb *payoutBeneficiaryRepository) GetPayoutBeneficiariesByMobileNumber(
	c echo.Context,
) ([]models.GetPayoutBeneficiaryResponseModel, error) {

	mobileNumber := c.Param("mobile_number")
	if mobileNumber == "" {
		return nil, fmt.Errorf("mobile_number is required")
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return rb.db.GetPayoutBeneficiariesByMobileNumberQuery(
		ctx,
		mobileNumber,
		limit,
		offset,
	)
}
