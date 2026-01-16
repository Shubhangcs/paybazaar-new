package repositories

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type WalletTransactionInterface interface {
	CreateWalletTransaction(echo.Context) error
	GetAdminWalletTransactions(echo.Context) ([]models.GetWalletTransactionResponseModel, error)
	GetMasterDistributorWalletTransactions(echo.Context) ([]models.GetWalletTransactionResponseModel, error)
	GetDistributorWalletTransactions(echo.Context) ([]models.GetWalletTransactionResponseModel, error)
	GetRetailerWalletTransactions(echo.Context) ([]models.GetWalletTransactionResponseModel, error)
	GetAdminWalletBalance(echo.Context) (float64, error)
	GetMasterDistributorWalletBalance(echo.Context) (float64, error)
	GetDistributorWalletBalance(echo.Context) (float64, error)
	GetRetailerWalletBalance(echo.Context) (float64, error)
}

type walletTransactionRepository struct {
	db *database.Database
}

func NewWalletTransactionRepository(
	db *database.Database,
) *walletTransactionRepository {
	return &walletTransactionRepository{
		db: db,
	}
}

func (wr *walletTransactionRepository) CreateWalletTransaction(
	c echo.Context,
) error {

	var req models.CreateWalletTransactionRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return wr.db.CreateWalletTransactionQuery(ctx, req)
}

func (wr *walletTransactionRepository) GetAdminWalletTransactions(
	c echo.Context,
) ([]models.GetWalletTransactionResponseModel, error) {

	adminID := c.Param("admin_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return wr.db.GetAdminWalletTransactionsQuery(
		ctx,
		adminID,
		limit,
		offset,
	)
}

func (wr *walletTransactionRepository) GetMasterDistributorWalletTransactions(
	c echo.Context,
) ([]models.GetWalletTransactionResponseModel, error) {

	mdID := c.Param("master_distributor_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return wr.db.GetMasterDistributorWalletTransactionsQuery(
		ctx,
		mdID,
		limit,
		offset,
	)
}

func (wr *walletTransactionRepository) GetDistributorWalletTransactions(
	c echo.Context,
) ([]models.GetWalletTransactionResponseModel, error) {

	distributorID := c.Param("distributor_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return wr.db.GetDistributorWalletTransactionsQuery(
		ctx,
		distributorID,
		limit,
		offset,
	)
}

func (wr *walletTransactionRepository) GetRetailerWalletTransactions(
	c echo.Context,
) ([]models.GetWalletTransactionResponseModel, error) {

	retailerID := c.Param("retailer_id")
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	return wr.db.GetRetailerWalletTransactionsQuery(
		ctx,
		retailerID,
		limit,
		offset,
	)
}

func (wr *walletTransactionRepository) GetAdminWalletBalance(
	c echo.Context,
) (float64, error) {

	adminID := c.Param("admin_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return wr.db.GetAdminWalletBalanceQuery(ctx, adminID)
}

func (wr *walletTransactionRepository) GetMasterDistributorWalletBalance(
	c echo.Context,
) (float64, error) {

	mdID := c.Param("master_distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return wr.db.GetMasterDistributorWalletBalanceQuery(ctx, mdID)
}

func (wr *walletTransactionRepository) GetDistributorWalletBalance(
	c echo.Context,
) (float64, error) {

	distributorID := c.Param("distributor_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return wr.db.GetDistributorWalletBalanceQuery(ctx, distributorID)
}

func (wr *walletTransactionRepository) GetRetailerWalletBalance(
	c echo.Context,
) (float64, error) {

	retailerID := c.Param("retailer_id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	return wr.db.GetRetailerWalletBalanceQuery(ctx, retailerID)
}
