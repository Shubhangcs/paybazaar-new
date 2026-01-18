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

type RetailerInterface interface {
	CreateRetailer(echo.Context) error
	GetRetailerByID(echo.Context) (*models.RetailerModel, error)
	ListRetailers(echo.Context) ([]models.GetRetailerResponseModel, error)
	ListRetailersByDistributorID(echo.Context) ([]models.GetRetailerResponseModel, error)
	ListRetailersByMasterDistributorID(echo.Context) ([]models.GetRetailerResponseModel, error)
	UpdateRetailer(echo.Context) error
	DeleteRetailer(echo.Context) error
	GetRetailersByDistributorIDForDropdown(echo.Context) ([]models.DropdownModel, error)
	LoginRetailer(echo.Context) (string, error)
}

type retailerRepository struct {
	db       *database.Database
	jwtUtils *pkg.JwtUtils
}

func NewRetailerRepository(
	db *database.Database,
	jwtUtils *pkg.JwtUtils,
) *retailerRepository {
	return &retailerRepository{
		db:       db,
		jwtUtils: jwtUtils,
	}
}

func (rr *retailerRepository) CreateRetailer(c echo.Context) error {
	var req models.CreateRetailerRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return rr.db.CreateRetailerQuery(ctx, req)
}

func (rr *retailerRepository) GetRetailerByID(
	c echo.Context,
) (*models.RetailerModel, error) {

	retailerID := c.Param("retailer_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return rr.db.GetRetailerByIDQuery(ctx, retailerID)
}

func (rr *retailerRepository) ListRetailers(
	c echo.Context,
) ([]models.GetRetailerResponseModel, error) {

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return rr.db.ListRetailersQuery(ctx, limit, offset)
}

func (rr *retailerRepository) ListRetailersByDistributorID(
	c echo.Context,
) ([]models.GetRetailerResponseModel, error) {

	distributorID := c.Param("distributor_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return rr.db.ListRetailersByDistributorIDQuery(
		ctx,
		distributorID,
		limit,
		offset,
	)
}

func (rr *retailerRepository) ListRetailersByMasterDistributorID(
	c echo.Context,
) ([]models.GetRetailerResponseModel, error) {

	masterDistributorID := c.Param("master_distributor_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return rr.db.ListRetailersByMasterDistributorIDQuery(
		ctx,
		masterDistributorID,
		limit,
		offset,
	)
}

func (rr *retailerRepository) UpdateRetailer(c echo.Context) error {
	retailerID := c.Param("retailer_id")

	var req models.UpdateRetailerRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return rr.db.UpdateRetailerQuery(ctx, retailerID, req)
}

func (rr *retailerRepository) DeleteRetailer(c echo.Context) error {
	retailerID := c.Param("retailer_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return rr.db.DeleteRetailerQuery(ctx, retailerID)
}

func (rr *retailerRepository) GetRetailersByDistributorIDForDropdown(
	c echo.Context,
) ([]models.DropdownModel, error) {

	distributorID := c.Param("distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return rr.db.GetRetailersByDistributorIDForDropdownQuery(
		ctx,
		distributorID,
	)
}

func (rr *retailerRepository) LoginRetailer(c echo.Context) (string, error) {
	var req models.LoginRetailerModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	res, err := rr.db.GetRetailerByPhoneQuery(ctx, req.RetailerPhoneNumber)
	if err != nil {
		return "", err
	}
	if res.Password != req.RetailerPassword || res.IsBlocked {
		return "", fmt.Errorf("incorrect password or retailer is blocked")
	}
	return rr.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		AdminID:  res.AdminID,
		UserID:   res.RetailerID,
		UserName: res.Name,
		UserRole: "retailer",
	})
}
