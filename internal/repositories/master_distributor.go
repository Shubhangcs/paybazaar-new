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

type MasterDistributorInterface interface {
	CreateMasterDistributor(echo.Context) error
	GetMasterDistributorDetailsByMasterDistributorID(echo.Context) (*models.GetCompleteMasterDistributorDetailsResponseModel, error)
	GetMasterDistributorsByAdminID(echo.Context) ([]models.GetCompleteMasterDistributorDetailsResponseModel, error)
	GetMasterDistributorsForDropdownByAdminID(echo.Context) ([]models.GetMasterDistributorForDropdownModel, error)
	UpdateMasterDistributorDetails(echo.Context) error
	UpdateMasterDistributorPassword(echo.Context) error
	UpdateMasterDistributorBlockStatus(echo.Context) error
	UpdateMasterDistributorKYCStatus(echo.Context) error
	UpdateMasterDistributorMPIN(echo.Context) error
	DeleteMasterDistributor(echo.Context) error
	MasterDistributorLogin(echo.Context) (string, error)
}

type masterDistributorRepository struct {
	db       *database.Database
	jwtUtils *pkg.JwtUtils
}

func NewMasterDistributorRepository(
	db *database.Database,
	jwtUtils *pkg.JwtUtils,
) *masterDistributorRepository {
	return &masterDistributorRepository{
		db:       db,
		jwtUtils: jwtUtils,
	}
}

func (mr *masterDistributorRepository) CreateMasterDistributor(c echo.Context) error {
	var req models.CreateMasterDistributorRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return mr.db.CreateMasterDistributorQuery(ctx, req)
}

func (mr *masterDistributorRepository) GetMasterDistributorDetailsByMasterDistributorID(c echo.Context) (*models.GetCompleteMasterDistributorDetailsResponseModel, error) {
	var masterDistributorID = c.Param("master_distributor_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.GetMasterDistributorDetailsByMasterDistributorIDQuery(ctx, masterDistributorID)
}

func (mr *masterDistributorRepository) GetMasterDistributorsByAdminID(c echo.Context) ([]models.GetCompleteMasterDistributorDetailsResponseModel, error) {
	var adminID = c.Param("admin_id")
	limit, offset := parsePagination(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.GetMasterDistributorsByAdminIDQuery(ctx, adminID, limit, offset)
}

func (mr *masterDistributorRepository) GetMasterDistributorsForDropdownByAdminID(c echo.Context) ([]models.GetMasterDistributorForDropdownModel, error) {
	var adminID = c.Param("admin_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.GetMasterDistributorsForDropdownByAdminIDQuery(ctx, adminID)
}

func (mr *masterDistributorRepository) UpdateMasterDistributorDetails(c echo.Context) error {
	var req models.UpdateMasterDistributorDetailsRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.UpdateMasterDistributorDetailsQuery(ctx, req)
}

func (mr *masterDistributorRepository) UpdateMasterDistributorPassword(c echo.Context) error {
	var req models.UpdateMasterDistributorPasswordRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.UpdateMasterDistributorPasswordQuery(ctx, req)
}

func (mr *masterDistributorRepository) UpdateMasterDistributorBlockStatus(c echo.Context) error {
	var req models.UpdateMasterDistributorBlockStatusRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.UpdateMasterDistributorBlockStatusQuery(ctx, req)
}

func (mr *masterDistributorRepository) UpdateMasterDistributorKYCStatus(c echo.Context) error {
	var req models.UpdateMasterDistributorKYCStatusRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.UpdateMasterDistributorKYCStatusQuery(ctx, req)
}

func (mr *masterDistributorRepository) UpdateMasterDistributorMPIN(c echo.Context) error {
	var req models.UpdateMasterDistributorMPINRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.UpdateMasterDistributorMPINQuery(ctx, req)
}

func (mr *masterDistributorRepository) DeleteMasterDistributor(c echo.Context) error {
	var masterDistributorID = c.Param("master_distributor_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return mr.db.DeleteMasterDistributorQuery(ctx, masterDistributorID)
}

func (mr *masterDistributorRepository) MasterDistributorLogin(c echo.Context) (string, error) {
	var req models.GetMasterDistributorDetailsForLoginModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	details, err := mr.db.GetMasterDistributorDetailsForLoginQuery(ctx, req.MasterDistributorID)
	if err != nil {
		return "", err
	}

	if details.Password != req.Password {
		return "", fmt.Errorf("incorrect password")
	}

	if details.IsBlocked {
		return "", fmt.Errorf("master distributor is blocked")
	}

	token, err := mr.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		AdminID:  details.AdminID,
		UserName: details.MasterDistributorName,
		UserID:   details.MasterDistributorID,
		UserRole: "master_distributor",
	})
	if err != nil {
		return "", err
	}
	return token, nil
}
