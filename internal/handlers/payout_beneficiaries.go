package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type PayoutBeneficiaryHandler struct {
	repo repositories.PayoutBeneficiaryInterface
}

func NewPayoutBeneficiaryHandler(
	repo repositories.PayoutBeneficiaryInterface,
) *PayoutBeneficiaryHandler {
	return &PayoutBeneficiaryHandler{
		repo: repo,
	}
}

func (h *PayoutBeneficiaryHandler) CreatePayoutBeneficiary(
	c echo.Context,
) error {

	if err := h.repo.CreatePayoutBeneficiary(c); err != nil {
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

func (h *PayoutBeneficiaryHandler) GetPayoutBeneficiaryByID(
	c echo.Context,
) error {

	data, err := h.repo.GetPayoutBeneficiaryByID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Message: "bene fetched successfully",
		Status:  "success",
		Data:    data,
	})
}

func (h *PayoutBeneficiaryHandler) ListPayoutBeneficiaries(
	c echo.Context,
) error {

	data, err := h.repo.ListPayoutBeneficiaries(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Message: "bene fetched successfully",
		Status:  "success",
		Data:    data,
	})
}

func (h *PayoutBeneficiaryHandler) UpdatePayoutBeneficiary(
	c echo.Context,
) error {

	if err := h.repo.UpdatePayoutBeneficiary(c); err != nil {
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

func (h *PayoutBeneficiaryHandler) UpdatePayoutBeneficiaryVerification(
	c echo.Context,
) error {

	if err := h.repo.UpdatePayoutBeneficiaryVerification(c); err != nil {
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

func (h *PayoutBeneficiaryHandler) DeletePayoutBeneficiary(
	c echo.Context,
) error {

	if err := h.repo.DeletePayoutBeneficiary(c); err != nil {
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

func (h *PayoutBeneficiaryHandler) GetBeneficiariesByMobileNumber(
	c echo.Context,
) error {

	list, err := h.repo.GetPayoutBeneficiariesByMobileNumber(c)
	if err != nil {
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
			Status: "success",
			Data:   list,
		},
	)
}
