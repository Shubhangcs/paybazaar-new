package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type payoutHandler struct {
	payoutRepo repositories.PayoutInterface
}

func NewPayoutHandler(payoutRepo repositories.PayoutInterface) *payoutHandler {
	return &payoutHandler{
		payoutRepo: payoutRepo,
	}
}

func (ph *payoutHandler) CreatePayout(c echo.Context) error {
	if err := ph.payoutRepo.CreatePayout(c); err != nil {
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
			Message: "payout request submitted successfully",
		},
	)
}

func (ph *payoutHandler) GetAllPayoutsRequest(c echo.Context) error {
	res, err := ph.payoutRepo.GetAllPayouts(c)
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
			Message: "payout request submitted successfully",
			Data: map[string]any{
				"payout_transactions": res,
			},
		},
	)
}

func (ph *payoutHandler) GetPayoutsByRetailerIDRequest(c echo.Context) error {
	res, err := ph.payoutRepo.GetPayoutsByRetailerID(c)
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
			Message: "payout request submitted successfully",
			Data: map[string]any{
				"payout_transactions": res,
			},
		},
	)
}
