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
