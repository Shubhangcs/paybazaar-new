package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type payoutHandler struct {
	payoutRepository repositories.PayoutInterface
}

func NewPayoutHandler(payoutRepository repositories.PayoutInterface) *payoutHandler {
	return &payoutHandler{
		payoutRepository,
	}
}

func (ph *payoutHandler) CreatePayoutRequest(c echo.Context) error {
	err := ph.payoutRepository.CreatePayoutTransaction(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK, models.ResponseModel{Message: "payout transaction successfull", Status: "success"})
}

func (ph *payoutHandler) GetAllPayoutTransactionsRequest(c echo.Context) error {
	res, err := ph.payoutRepository.GetAllPayoutTransactions(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK, models.ResponseModel{Message: "payout transaction fetched successfully", Status: "success", Data: map[string]any{"transactions": res}})
}

func (ph *payoutHandler) GetPayoutTransactionsByRetailerIdRequest(c echo.Context) error {
	res, err := ph.payoutRepository.GetPayoutTransactionsByRetailerId(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK, models.ResponseModel{Message: "payout transaction fetched successfully", Status: "success", Data: map[string]any{"transactions": res}})
}
