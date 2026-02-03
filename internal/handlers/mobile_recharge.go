package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type mobileRechargeHandler struct {
	mobileRechargeRepository repositories.MobileRechargeInterface
}

func NewMobileRechargeHandler(mobileRechargeRepository repositories.MobileRechargeInterface) *mobileRechargeHandler {
	return &mobileRechargeHandler{
		mobileRechargeRepository,
	}
}

func (mrh *mobileRechargeHandler) CreateMobileRechargeRequest(c echo.Context) error {
	if err := mrh.mobileRechargeRepository.CreateMobileRecharge(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "mobile recharge successfull"},
	)
}

func (mrh *mobileRechargeHandler) GetMobileRechargeCirclesRequest(c echo.Context) error {
	res, err := mrh.mobileRechargeRepository.GetAllMobileRechargeCircles(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "circles fetched successfully",
		Data:    map[string]any{"circles": res},
	})
}

func (mrh *mobileRechargeHandler) GetMobileRechargeOperatorsRequest(c echo.Context) error {
	res, err := mrh.mobileRechargeRepository.GetAllMobileRechargeOperators(c)
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

func (mrh *mobileRechargeHandler) GetAllMobileRechargesRequest(c echo.Context) error {
	res, err := mrh.mobileRechargeRepository.GetAllMobileRecharges(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "mobile recharges fetched successfully",
		Data:    map[string]any{"recharges": res},
	})
}

func (mrh *mobileRechargeHandler) GetMobileRechargesByRetailerIDRequest(c echo.Context) error {
	res, err := mrh.mobileRechargeRepository.GetMobileRechargesByRetailerID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "mobile recharges fetched successfully",
		Data:    map[string]any{"recharges": res},
	})
}

func (mrh *mobileRechargeHandler) GetMobileRechargePlansRequest(c echo.Context) error {
	res, err := mrh.mobileRechargeRepository.GetAllPlansBasedOnCircleAndOperator(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "plans fetched successfully",
		Data:    res,
	})
}

func (mrh *mobileRechargeHandler) MobileRechargeRefundRequest(c echo.Context) error {
	if err := mrh.mobileRechargeRepository.MobileRechargeRefund(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "mobile recharge refund successfull"},
	)
}
