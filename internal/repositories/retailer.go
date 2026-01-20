package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/pkg"
)

type RetailerInterface interface {
	CreateRetailer(echo.Context) error
	GetRetailerDetailsByRetailerID(echo.Context) (*models.GetCompleteRetailerDetailsResponseModel, error)
	GetRetailersByAdminID(echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error)
	GetRetailersByMasterDistributorID(echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error)
	GetRetailersByDistributorID(echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error)
	UpdateRetailerDetails(echo.Context) error
	UpdateRetailerPassword(echo.Context) error
	UpdateRetailerKYCStatus(echo.Context) error
	UpdateRetailerBlockStatus(echo.Context) error
	UpdateRetailerMPIN(echo.Context) error
	UpdateRetailerDistributor(echo.Context) error
	DeleteRetailer(echo.Context) error
	RetailerLogin(echo.Context) (string, error)
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

func (rr *retailerRepository) GetRetailerDetailsByRetailerID(c echo.Context) (*models.GetCompleteRetailerDetailsResponseModel, error) {
	var retailerID = c.Param("retailer_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.GetRetailerDetailsByRetailerIDQuery(ctx, retailerID)
}

func (rr *retailerRepository) GetRetailersByAdminID(c echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error) {
	var adminID = c.Param("admin_id")
	limit, offset := parsePagination(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.GetRetailersByAdminIDQuery(ctx, adminID, limit, offset)
}

func (rr *retailerRepository) GetRetailersByMasterDistributorID(c echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error) {
	var masterDistributorID = c.Param("master_distributor_id")
	limit, offset := parsePagination(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.GetRetailersByMasterDistributorIDQuery(ctx, masterDistributorID, limit, offset)
}

func (rr *retailerRepository) GetRetailersByDistributorID(c echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error) {
	var distributorID = c.Param("distributor_id")
	limit, offset := parsePagination(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.GetRetailersByDistributorIDQuery(ctx, distributorID, limit, offset)
}

func (rr *retailerRepository) UpdateRetailerDetails(c echo.Context) error {
	var req models.UpdateRetailerDetailsRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.UpdateRetailerDetailsQuery(ctx,req)
}
