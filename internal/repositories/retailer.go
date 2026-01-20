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
	GetRetailerDetailsByRetailerID(echo.Context) (*models.GetCompleteRetailerDetailsResponseModel, error)
	GetRetailersByAdminID(echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error)
	GetRetailersByMasterDistributorID(echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error)
	GetRetailersByDistributorID(echo.Context) ([]models.GetCompleteRetailerDetailsResponseModel, error)
	GetRetailersForDropdownByDistributorID(echo.Context) ([]models.GetRetailerForDropdownModel, error)
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
	return rr.db.UpdateRetailerDetailsQuery(ctx, req)
}

func (rr *retailerRepository) UpdateRetailerPassword(c echo.Context) error {
	var req models.UpdateRetailerPasswordRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.UpdateRetailerPasswordQuery(ctx, req)
}

func (rr *retailerRepository) UpdateRetailerKYCStatus(c echo.Context) error {
	var req models.UpdateRetailerKYCStatusRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.UpdateRetailerKYCStatusQuery(ctx, req)
}

func (rr *retailerRepository) UpdateRetailerBlockStatus(c echo.Context) error {
	var req models.UpdateRetailerBlockStatusRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.UpdateRetailerBlockStatusQuery(ctx, req)
}

func (rr *retailerRepository) UpdateRetailerMPIN(c echo.Context) error {
	var req models.UpdateRetailerMPINRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.UpdateRetailerMPINQuery(ctx, req)
}

func (rr *retailerRepository) UpdateRetailerDistributor(c echo.Context) error {
	var req models.UpdateRetailerDistributorRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.UpdateRetailerDistributorQuery(ctx, req)
}

func (rr *retailerRepository) DeleteRetailer(c echo.Context) error {
	var retailerID = c.Param("retailer_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return rr.db.DeleteRetailerQuery(ctx, retailerID)
}

func (rr *retailerRepository) RetailerLogin(c echo.Context) (string, error) {
	var req models.RetailerLoginRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	details, err := rr.db.GetRetailerDetailsForLoginQuery(ctx, req.RetailerID)
	if err != nil {
		return "", err
	}

	if details.Password != req.RetailerPassword {
		return "", fmt.Errorf("incorrect password")
	}

	if details.IsBlocked {
		return "", fmt.Errorf("distributor is blocked")
	}

	token, err := rr.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		AdminID:  details.AdminID,
		UserName: details.RetailerName,
		UserID:   details.RetailerID,
		UserRole: "retailer",
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (rr *retailerRepository) GetRetailersForDropdownByDistributorID(c echo.Context) ([]models.GetRetailerForDropdownModel, error) {
	var distributorID = c.Param("distributor_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return rr.db.GetRetailersForDropdownByDistributorIDQuery(ctx, distributorID)
}
