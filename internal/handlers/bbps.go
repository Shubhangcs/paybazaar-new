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

func (bh *bbpsHandler) GetAllPostpaidMobileRechargeRequest(c echo.Context) error {
	res, err := bh.bbpsRepository.GetAllPostpaidMobileRecharge(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "postpaid mobile recharges fetched successfully", Data: map[string]any{"history": res}},
	)
}

func (bh *bbpsHandler) GetPostpaidMobileRechargeByRetailerIDRequest(c echo.Context) error {
	res, err := bh.bbpsRepository.GetPostpaidMobileRechargeByRetailerID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "postpaid mobile recharges fetched successfully", Data: map[string]any{"history": res}},
	)
}

func (bh *bbpsHandler) CreateElectricityBillPaymentRequest(c echo.Context) error {
	if err := bh.bbpsRepository.CreateElectricityBillPayment(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "electricity bill paid successfully"},
	)
}

func (bh *bbpsHandler) GetAllElectricityBillOperatorsRequest(c echo.Context) error {
	res, err := bh.bbpsRepository.GetAllElectricityOperators(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "operators fetched successfully", Data: map[string]any{"operators": res}},
	)
}

func (bh *bbpsHandler) GetElectricityBillBalanceRequest(c echo.Context) error {
	res, err := bh.bbpsRepository.GetElectricityBillFetchBalance(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "balance fetched successfully", Data: map[string]any{"response": res}},
	)
}

func (bh *bbpsHandler) GetAllElectricityBillHistoryRequest(c echo.Context) error {
	res, err := bh.bbpsRepository.GetAllElectricityBillPaymentTransactions(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "transactions fetched successfully", Data: map[string]any{"transactions": res}},
	)
}

func (bh *bbpsHandler) GetElectricityBillHistoryByRetailerIDRequest(c echo.Context) error {
	res, err := bh.bbpsRepository.GetElectricityBillPaymentTransactionsByRetailerID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "transactions fetched successfully", Data: map[string]any{"transactions": res}},
	)
}

func (bh *bbpsHandler) ElectricityBillPaymentRefundRequest(c echo.Context) error {
	if err := bh.bbpsRepository.ElectricityBillPaymentRefund(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "electricity bill refund successfull"},
	)
}

func (bh *bbpsHandler) MobileRechargePostpaidRefundRequest(c echo.Context) error {
	if err := bh.bbpsRepository.PostpaidMobileRechargeRefund(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "postpaid mobile recharge refund successfull"},
	)
}
