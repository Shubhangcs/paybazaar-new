package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type RetailerBeneficiaryHandler struct {
	repo repositories.RetailerBeneficiaryInterface
}

func NewRetailerBeneficiaryHandler(
	repo repositories.RetailerBeneficiaryInterface,
) *RetailerBeneficiaryHandler {
	return &RetailerBeneficiaryHandler{
		repo: repo,
	}
}

func (h *RetailerBeneficiaryHandler) CreateRetailerBeneficiary(
	c echo.Context,
) error {

	if err := h.repo.CreateRetailerBeneficiary(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiary created successfully",
	})
}

func (h *RetailerBeneficiaryHandler) GetRetailerBeneficiaryByID(
	c echo.Context,
) error {

	data, err := h.repo.GetRetailerBeneficiaryByID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status: "success",
		Data:   data,
	})
}

func (h *RetailerBeneficiaryHandler) ListRetailerBeneficiaries(
	c echo.Context,
) error {

	data, err := h.repo.ListRetailerBeneficiaries(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status: "success",
		Data:   data,
	})
}

func (h *RetailerBeneficiaryHandler) UpdateRetailerBeneficiary(
	c echo.Context,
) error {

	if err := h.repo.UpdateRetailerBeneficiary(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiary updated successfully",
	})
}

func (h *RetailerBeneficiaryHandler) UpdateRetailerBeneficiaryVerification(
	c echo.Context,
) error {

	if err := h.repo.UpdateRetailerBeneficiaryVerification(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiary verification updated",
	})
}

func (h *RetailerBeneficiaryHandler) DeleteRetailerBeneficiary(
	c echo.Context,
) error {

	if err := h.repo.DeleteRetailerBeneficiary(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "beneficiary deleted successfully",
	})
}
