package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type bbpsHandler struct {
	bbpsRepository repositories.BBPSInterface
}

func NewBBPSHandler(bbpsRepository repositories.BBPSInterface) *bbpsHandler {
	return &bbpsHandler{
		bbpsRepository,
	}
}

func (bh *bbpsHandler) CreatePostpaidMobileRechargeRequest(c echo.Context) error {
	if err := bh.bbpsRepository.CreatePostpaidMobileRecharge(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "postpaid mobile recharge successfull"},
	)
}

func (bh *bbpsHandler) GetPostpaidMobileRechargeBalanceRequest(c echo.Context) error {
	res, err := bh.bbpsRepository.GetPostpaidMobileRechargeBalance(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "postpaid mobile recharge balance fetched successfully", Data: map[string]any{"response": res}},
	)
}
