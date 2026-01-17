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
	GetDistributorByID(echo.Context) (*models.DistributorModel, error)
	ListDistributors(echo.Context) ([]models.GetDistributorResponseModel, error)
	ListDistributorsByMasterDistributorID(echo.Context) ([]models.GetDistributorResponseModel, error)
	UpdateDistributor(echo.Context) error
	DeleteDistributor(echo.Context) error
	GetDistributorsByMasterDistributorIDForDropdown(echo.Context) ([]models.DropdownModel, error)
	LoginDistributor(echo.Context) (string, error)
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

func (dr *distributorRepository) GetDistributorByID(
	c echo.Context,
) (*models.DistributorModel, error) {

	distributorID := c.Param("distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return dr.db.GetDistributorByIDQuery(ctx, distributorID)
}

func (dr *distributorRepository) ListDistributors(
	c echo.Context,
) ([]models.GetDistributorResponseModel, error) {

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return dr.db.ListDistributorsQuery(ctx, limit, offset)
}

func (dr *distributorRepository) ListDistributorsByMasterDistributorID(
	c echo.Context,
) ([]models.GetDistributorResponseModel, error) {

	masterDistributorID := c.Param("master_distributor_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return dr.db.ListDistributorsByMasterDistributorIDQuery(
		ctx,
		masterDistributorID,
		limit,
		offset,
	)
}

func (dr *distributorRepository) UpdateDistributor(c echo.Context) error {
	distributorID := c.Param("distributor_id")

	var req models.UpdateDistributorRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return dr.db.UpdateDistributorQuery(ctx, distributorID, req)
}

func (dr *distributorRepository) DeleteDistributor(c echo.Context) error {
	distributorID := c.Param("distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return dr.db.DeleteDistributorQuery(ctx, distributorID)
}

func (dr *distributorRepository) GetDistributorsByMasterDistributorIDForDropdown(
	c echo.Context,
) ([]models.DropdownModel, error) {

	masterDistributorID := c.Param("master_distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return dr.db.GetDistributorsByMasterDistributorIDForDropdownQuery(
		ctx,
		masterDistributorID,
	)
}

func (dr *distributorRepository) LoginDistributor(c echo.Context) (string, error) {
	var req models.LoginDistributorModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	res, err := dr.db.GetDistributorByIDQuery(ctx, req.DistributorID)
	if err != nil {
		return "", err
	}
	if res.Password != req.DistributorPassword {
		return "", fmt.Errorf("incorrect password")
	}
	return dr.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		AdminID: res.AdminID,
		UserID:   res.DistributorID,
		UserName: res.Name,
		UserRole: "distributor",
	})
}
