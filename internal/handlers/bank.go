package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type bankHandler struct {
	bankRepo repositories.BankInterface
}

func NewBankHandler(bankRepo repositories.BankInterface) *bankHandler {
	return &bankHandler{
		bankRepo: bankRepo,
	}
}

func (bh *bankHandler) CreateBankRequest(c echo.Context) error {
	err := bh.bankRepo.CreateBank(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "bank created successfully",
	})
}

func (bh *bankHandler) GetBankDetailsByBankIDRequest(c echo.Context) error {
	data, err := bh.bankRepo.GetBankDetailsByBankID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "bank fetched successfully",
		Data: map[string]any{
			"bank": data,
		},
	})
}

func (bh *bankHandler) GetAllBanksRequest(c echo.Context) error {
	data, err := bh.bankRepo.GetAllBanks(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "banks fetched successfully",
		Data: map[string]any{
			"banks": data,
		},
	})
}

func (bh *bankHandler) UpdateBankDetailsRequest(c echo.Context) error {
	if err := bh.bankRepo.UpdateBankDetails(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "bank details updated successfully",
	})
}

func (bh *bankHandler) DeleteBankRequest(c echo.Context) error {
	if err := bh.bankRepo.DeleteBank(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "bank deleted successfully",
	})
}

func (bh *bankHandler) CreateAdminBankRequest(c echo.Context) error {
	err := bh.bankRepo.CreateAdminBank(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin bank created successfully",
	})
}

func (bh *bankHandler) GetAdminBankDetailsByAdminBankIDRequest(c echo.Context) error {
	data, err := bh.bankRepo.GetAdminBankDetailsByAdminBankID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin bank fetched successfully",
		Data: map[string]any{
			"admin_bank": data,
		},
	})
}

func (bh *bankHandler) GetAdminBanksByAdminIDRequest(c echo.Context) error {
	data, err := bh.bankRepo.GetAdminBanksByAdminID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin banks fetched successfully",
		Data: map[string]any{
			"admin_banks": data,
		},
	})
}

func (bh *bankHandler) UpdateAdminBankDetailsRequest(c echo.Context) error {
	if err := bh.bankRepo.UpdateAdminBankDetails(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin bank details updated successfully",
	})
}

func (bh *bankHandler) DeleteAdminBankRequest(c echo.Context) error {
	if err := bh.bankRepo.DeleteAdminBank(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin bank deleted successfully",
	})
}
