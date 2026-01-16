package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type walletTransactionHandler struct {
	walletRepository repositories.WalletTransactionInterface
}

func NewWalletTransactionHandler(
	walletRepository repositories.WalletTransactionInterface,
) *walletTransactionHandler {
	return &walletTransactionHandler{
		walletRepository: walletRepository,
	}
}

// ============================
// CREATE WALLET TRANSACTION
// ============================
func (wh *walletTransactionHandler) CreateWalletTransactionRequest(
	c echo.Context,
) error {

	if err := wh.walletRepository.CreateWalletTransaction(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet transaction created successfully",
		},
	)
}

// ============================
// GET WALLET TRANSACTIONS
// ============================

func (wh *walletTransactionHandler) GetAdminWalletTransactionsRequest(
	c echo.Context,
) error {

	res, err := wh.walletRepository.GetAdminWalletTransactions(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet transactions fetched successfully",
			Data:    map[string]any{"transactions": res},
		},
	)
}

func (wh *walletTransactionHandler) GetMasterDistributorWalletTransactionsRequest(
	c echo.Context,
) error {

	res, err := wh.walletRepository.GetMasterDistributorWalletTransactions(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet transactions fetched successfully",
			Data:    map[string]any{"transactions": res},
		},
	)
}

func (wh *walletTransactionHandler) GetDistributorWalletTransactionsRequest(
	c echo.Context,
) error {

	res, err := wh.walletRepository.GetDistributorWalletTransactions(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet transactions fetched successfully",
			Data:    map[string]any{"transactions": res},
		},
	)
}

func (wh *walletTransactionHandler) GetRetailerWalletTransactionsRequest(
	c echo.Context,
) error {

	res, err := wh.walletRepository.GetRetailerWalletTransactions(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet transactions fetched successfully",
			Data:    map[string]any{"transactions": res},
		},
	)
}

// ============================
// GET WALLET BALANCE
// ============================

func (wh *walletTransactionHandler) GetAdminWalletBalanceRequest(
	c echo.Context,
) error {

	balance, err := wh.walletRepository.GetAdminWalletBalance(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet balance fetched successfully",
			Data:    map[string]any{"wallet_balance": balance},
		},
	)
}

func (wh *walletTransactionHandler) GetMasterDistributorWalletBalanceRequest(
	c echo.Context,
) error {

	balance, err := wh.walletRepository.GetMasterDistributorWalletBalance(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet balance fetched successfully",
			Data:    map[string]any{"wallet_balance": balance},
		},
	)
}

func (wh *walletTransactionHandler) GetDistributorWalletBalanceRequest(
	c echo.Context,
) error {

	balance, err := wh.walletRepository.GetDistributorWalletBalance(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet balance fetched successfully",
			Data:    map[string]any{"wallet_balance": balance},
		},
	)
}

func (wh *walletTransactionHandler) GetRetailerWalletBalanceRequest(
	c echo.Context,
) error {

	balance, err := wh.walletRepository.GetRetailerWalletBalance(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{
				Status:  "failed",
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "wallet balance fetched successfully",
			Data:    map[string]any{"wallet_balance": balance},
		},
	)
}
