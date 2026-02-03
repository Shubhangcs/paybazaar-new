package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type dthRechargeHandler struct {
	dthRechargeRepository repositories.DTHRechargeInterface
}

func NewDTHRechargeHandler(dthRechargeRepository repositories.DTHRechargeInterface) *dthRechargeHandler {
	return &dthRechargeHandler{
		dthRechargeRepository,
	}
}

func (dh *dthRechargeHandler) CreateDTHRechargeRequest(c echo.Context) error {
	if err := dh.dthRechargeRepository.CreateDTHRecharge(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "dth recharge successfull"},
	)
}

func (dh *dthRechargeHandler) GetAllDTHOperatorsRequest(c echo.Context) error {
	res, err := dh.dthRechargeRepository.GetAllDTHOperators(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "operators fetched successfully",
		Data:    map[string]any{"operators": res},
	})
}

func (dh *dthRechargeHandler) GetAllDTHRechargesRequest(c echo.Context) error {
	res, err := dh.dthRechargeRepository.GetAllDTHRecharges(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "dth recharges fetched successfully",
		Data:    map[string]any{"recharges": res},
	})
}

func (dh *dthRechargeHandler) GetDTHRechargesByRetailerIDRequest(c echo.Context) error {
	res, err := dh.dthRechargeRepository.GetDTHRechargesByRetailerID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "dth recharges fetched successfully",
		Data:    map[string]any{"recharges": res},
	})
}

func (dh *dthRechargeHandler) DTHRechargeRefundRequest(c echo.Context) error {
	if err := dh.dthRechargeRepository.DTHRechargeRefund(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "dth recharge refund successfull"},
	)
}
