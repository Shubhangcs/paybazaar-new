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

/* =========================================================
   BANKS
========================================================= */

// POST /banks
func (bh *bankHandler) CreateBank(c echo.Context) error {
	id, err := bh.bankRepo.CreateBank(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "bank created successfully",
		Data: map[string]int64{
			"bank_id": id,
		},
	})
}

// GET /banks/:bank_id
func (bh *bankHandler) GetBankByID(c echo.Context) error {
	data, err := bh.bankRepo.GetBankByID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "bank fetched successfully",
		Data:    data,
	})
}

// GET /banks
func (bh *bankHandler) GetAllBanks(c echo.Context) error {
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
		Data:    data,
	})
}

// PUT /banks/:bank_id
func (bh *bankHandler) UpdateBank(c echo.Context) error {
	if err := bh.bankRepo.UpdateBank(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "bank updated successfully",
	})
}

// DELETE /banks/:bank_id
func (bh *bankHandler) DeleteBank(c echo.Context) error {
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

/* =========================================================
   ADMIN BANKS
========================================================= */

// POST /admin-banks
func (bh *bankHandler) CreateAdminBank(c echo.Context) error {
	id, err := bh.bankRepo.CreateAdminBank(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin bank created successfully",
		Data: map[string]int64{
			"admin_bank_id": id,
		},
	})
}

// GET /admin-banks/:admin_bank_id
func (bh *bankHandler) GetAdminBankByID(c echo.Context) error {
	data, err := bh.bankRepo.GetAdminBankByID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin bank fetched successfully",
		Data:    data,
	})
}

// GET /admins/:admin_id/admin-banks
func (bh *bankHandler) GetAdminBanksByAdminID(c echo.Context) error {
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
		Data:    data,
	})
}

// PUT /admin-banks/:admin_bank_id
func (bh *bankHandler) UpdateAdminBank(c echo.Context) error {
	if err := bh.bankRepo.UpdateAdminBank(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "admin bank updated successfully",
	})
}

// DELETE /admin-banks/:admin_bank_id
func (bh *bankHandler) DeleteAdminBank(c echo.Context) error {
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
