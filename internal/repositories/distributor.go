package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/pkg"
)

type DistributorInterface interface {
	CreateDistributor(echo.Context) error
	GetDistributorDetailsByDistributorID(echo.Context) (*models.GetCompleteDistributorDetailsResponseModel, error)
	DistributorLogin(echo.Context) (string, error)
	GetDistributorsByAdminID(echo.Context) ([]models.GetCompleteDistributorDetailsResponseModel, error)
	GetDistributorsByMasterDistributorID(echo.Context) ([]models.GetCompleteDistributorDetailsResponseModel, error)
	GetDistributorsForDropdownByMasterDistributorID(echo.Context) ([]models.GetDistributorForDropdownModel, error)
	UpdateDistributorPassword(echo.Context) error
	UpdateDistributorBlockStatus(echo.Context) error
	UpdateDistributorKYCStatus(echo.Context) error
	UpdateDistributorDetails(echo.Context) error
	UpdateDistributorMasterDistributor(echo.Context) error
	UpdateDistributorMPIN(echo.Context) error
	DeleteDistributor(echo.Context) error
}

type distributorRepository struct {
	db       *database.Database
	jwtUtils *pkg.JwtUtils
}

func NewDistributorRepository(
	db *database.Database,
	jwtUtils *pkg.JwtUtils,
) *distributorRepository {
	return &distributorRepository{
		db:       db,
		jwtUtils: jwtUtils,
	}
}

func (dr *distributorRepository) CreateDistributor(c echo.Context) error {
	var req models.CreateDistributorRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return dr.db.CreateDistributorQuery(ctx, req)
}

func (dr *distributorRepository) GetDistributorDetailsByDistributorID(
	c echo.Context,
) (*models.GetCompleteDistributorDetailsResponseModel, error) {
	distributorID := c.Param("distributor_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return dr.db.GetDistributorDetailsByDistributorIDQuery(ctx, distributorID)
}

func (dr *distributorRepository) DistributorLogin(
	c echo.Context,
) (string, error) {
	var req models.DistributorLoginRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	details, err := dr.db.GetDistributorDetailsForLoginQuery(ctx, req.DistributorID)
	if err != nil {
		return "", err
	}

	if details.DistributorPassword != req.DistributorPassword {
		return "", fmt.Errorf("incorrect password")
	}

	if details.IsDistributorBlocked {
		return "", fmt.Errorf("distributor is blocked")
	}

	token, err := dr.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		AdminID:  details.AdminID,
		UserName: details.DistributorName,
		UserID:   details.DistributorID,
		UserRole: "distributor",
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (dr *distributorRepository) GetDistributorsByAdminID(
	c echo.Context,
) ([]models.GetCompleteDistributorDetailsResponseModel, error) {
	adminID := c.Param("admin_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return dr.db.GetDistributorsByAdminIDQuery(
		ctx,
		adminID,
		limit,
		offset,
	)
}

func (dr *distributorRepository) GetDistributorsByMasterDistributorID(
	c echo.Context,
) ([]models.GetCompleteDistributorDetailsResponseModel, error) {
	masterDistributorID := c.Param("master_distributor_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return dr.db.GetDistributorsByMasterDistributorIDQuery(
		ctx,
		masterDistributorID,
		limit,
		offset,
	)
}

func (dr *distributorRepository) GetDistributorsForDropdownByMasterDistributorID(
	c echo.Context,
) ([]models.GetDistributorForDropdownModel, error) {
	masterDistributorID := c.Param("master_distributor_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return dr.db.GetDistributorsForDropdownByMasterDistributorIDQuery(
		ctx,
		masterDistributorID,
	)
}

func (dr *distributorRepository) UpdateDistributorPassword(c echo.Context) error {
	var req models.UpdateDistributorPasswordRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return dr.db.UpdateDistributorPasswordQuery(ctx, req)
}

func (dr *distributorRepository) UpdateDistributorBlockStatus(c echo.Context) error {
	var req models.UpdateDistributorBlockStatusRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return dr.db.UpdateDistributorBlockStatusQuery(ctx, req)
}

func (dr *distributorRepository) UpdateDistributorKYCStatus(c echo.Context) error {
	var req models.UpdateDistributorKYCStatusRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return dr.db.UpdateDistributorKYCStatusQuery(ctx, req)
}

func (dr *distributorRepository) UpdateDistributorDetails(c echo.Context) error {
	var req models.UpdateDistributorDetailsRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return dr.db.UpdateDistributorDetailsQuery(ctx, req)
}

func (dr *distributorRepository) UpdateDistributorMasterDistributor(c echo.Context) error {
	var req models.UpdateDistributorMasterDistributorRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return dr.db.UpdateDistributorMasterDistributorQuery(ctx, req)
}

func (dr *distributorRepository) DeleteDistributor(
	c echo.Context,
) error {

	distributorID := c.Param("distributor_id")
	if distributorID == "" {
		return fmt.Errorf("distributor id is required")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return dr.db.DeleteDistributorQuery(ctx, distributorID)
}

func (dr *distributorRepository) UpdateDistributorMPIN(
	c echo.Context,
) error {

	var req models.UpdateDistributorMPINRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return dr.db.UpdateDistributorMPINQuery(ctx, req)
}
