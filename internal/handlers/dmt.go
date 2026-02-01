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
		Message: "wallet check completed successfully",
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
		Message: "wallet creation request sent successfully",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) VerifyDMTWalletCreationRequest(c echo.Context) error {
	res, err := dh.dmtRepository.VerifyDMTWalletCreation(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "wallet created successfully",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) CreateDMTBeneficiaryRequest(c echo.Context) error {
	res, err := dh.dmtRepository.CreateDMTBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiary created successfully",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) GetDMTBeneficiariesRequest(c echo.Context) error {
	res, err := dh.dmtRepository.GetDMTBeneficieries(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiaries fetched successfully",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) DeleteDMTBeneficiaryRequest(c echo.Context) error {
	res, err := dh.dmtRepository.DeleteDMTBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiary deletion request sent successfully",
		Data:    map[string]any{"response": res},
	})
}

func (dh *dmtHandler) VerifyDMTBeneficiaryDeleteRequest(c echo.Context) error {
	res, err := dh.dmtRepository.VerifyDMTBeneficiaryDelete(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiary deleted successfully",
		Data:    map[string]any{"response": res},
	})
}