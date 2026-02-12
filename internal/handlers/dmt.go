package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type dmtHandler struct {
	dmtRepository repositories.DMTInterface
}

func NewDMTHandler(dmtRepository repositories.DMTInterface) *dmtHandler {
	return &dmtHandler{
		dmtRepository,
	}
}

func (dh *dmtHandler) CheckDMTWalletExistsRequest(c echo.Context) error {
	res, err := dh.dmtRepository.CheckDMTWalletExists(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "dmt request success",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) CreateDMTWalletRequest(c echo.Context) error {
	res, err := dh.dmtRepository.CreateDMTWallet(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "dmt request success",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) VerifyDMTWalletRequest(c echo.Context) error {
	res, err := dh.dmtRepository.VerifyDMTWallet(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "dmt request success",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) AddDMTBeneficiaryRequest(c echo.Context) error {
	res, err := dh.dmtRepository.AddDMTBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "dmt request success",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) GetDMTBankListRequest(c echo.Context) error {
	res, err := dh.dmtRepository.GetDMTBankList(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "dmt request success",
		Data:    map[string]any{"response": res},
	})
}
