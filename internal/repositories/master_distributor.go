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
	GetMasterDistributorByID(echo.Context) (*models.MasterDistributorModel, error)
	ListMasterDistributors(echo.Context) ([]models.GetMasterDistributorResponseModel, error)
	ListMasterDistributorsByAdminID(echo.Context) ([]models.GetMasterDistributorResponseModel, error)
	UpdateMasterDistributor(echo.Context) error
	DeleteMasterDistributor(echo.Context) error
	GetMasterDistributorsByAdminIDForDropdown(echo.Context) ([]models.DropdownModel, error)
	LoginMasterDistributor(echo.Context) (string, error)
	UpdateBlockStatus(echo.Context) error
	UpdateKYCStatus(echo.Context) error
	UpdateMPIN(echo.Context) (int64, error)
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

func (mr *masterDistributorRepository) GetMasterDistributorByID(
	c echo.Context,
) (*models.MasterDistributorModel, error) {

	mdID := c.Param("master_distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return mr.db.GetMasterDistributorByIDQuery(ctx, mdID)
}

func (mr *masterDistributorRepository) ListMasterDistributors(
	c echo.Context,
) ([]models.GetMasterDistributorResponseModel, error) {

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return mr.db.ListMasterDistributorsQuery(ctx, limit, offset)
}

func (mr *masterDistributorRepository) ListMasterDistributorsByAdminID(
	c echo.Context,
) ([]models.GetMasterDistributorResponseModel, error) {

	adminID := c.Param("admin_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return mr.db.ListMasterDistributorsByAdminIDQuery(ctx, adminID, limit, offset)
}

func (mr *masterDistributorRepository) UpdateMasterDistributor(c echo.Context) error {
	mdID := c.Param("master_distributor_id")

	var req models.UpdateMasterDistributorRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return mr.db.UpdateMasterDistributorQuery(ctx, mdID, req)
}

func (mr *masterDistributorRepository) DeleteMasterDistributor(c echo.Context) error {
	mdID := c.Param("master_distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return mr.db.DeleteMasterDistributorQuery(ctx, mdID)
}

func (mr *masterDistributorRepository) GetMasterDistributorsByAdminIDForDropdown(
	c echo.Context,
) ([]models.DropdownModel, error) {

	adminID := c.Param("admin_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return mr.db.GetMasterDistributorsByAdminIDForDropdownQuery(ctx, adminID)
}

// Master Distributor Login
func (mdr *masterDistributorRepository) LoginMasterDistributor(c echo.Context) (string, error) {
	var req models.LoginMasterDistributorModel
	if err := bindAndValidate(c, &req); err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	res, err := mdr.db.GetMasterDistributorByIDQuery(ctx, req.MasterDistributorID)
	if err != nil {
		return "", err
	}
	if res.Password != req.MasterDistributorPassword || res.IsBlocked {
		return "", fmt.Errorf("incorrect password or md is blocked")
	}
	return mdr.jwtUtils.GenerateToken(ctx, models.AccessTokenClaims{
		AdminID:  res.AdminID,
		UserID:   res.MasterDistributorID,
		UserName: res.Name,
		UserRole: "master_distributor",
	})
}

func (mr *masterDistributorRepository) UpdateBlockStatus(
	c echo.Context,
) error {

	mdID := c.Param("master_distributor_id")
	if mdID == "" {
		return fmt.Errorf("master_distributor_id is required")
	}

	var req struct {
		IsBlocked bool `json:"is_blocked" validate:"required"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return mr.db.UpdateMasterDistributorBlockStatusQuery(
		ctx,
		mdID,
		req.IsBlocked,
	)
}

func (mr *masterDistributorRepository) UpdateKYCStatus(
	c echo.Context,
) error {

	mdID := c.Param("master_distributor_id")
	if mdID == "" {
		return fmt.Errorf("master_distributor_id is required")
	}

	var req struct {
		KYCStatus bool `json:"kyc_status" validate:"required"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return mr.db.UpdateMasterDistributorKYCStatusQuery(
		ctx,
		mdID,
		req.KYCStatus,
	)
}

func (mr *masterDistributorRepository) UpdateMPIN(
	c echo.Context,
) (int64, error) {

	mdID := c.Param("master_distributor_id")
	if mdID == "" {
		return 0, fmt.Errorf("master_distributor_id is required")
	}

	var req struct {
		MPIN int64 `json:"mpin" validate:"required,min=1000,max=9999"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(
		c.Request().Context(),
		30*time.Second,
	)
	defer cancel()

	return mr.db.UpdateMasterDistributorMPINQuery(
		ctx,
		mdID,
		req.MPIN,
	)
}
