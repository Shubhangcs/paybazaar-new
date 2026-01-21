package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type beneficiaryHandler struct {
	repo repositories.Beneficiary
}

func NewBeneficiaryHandler(repo repositories.Beneficiary) *beneficiaryHandler {
	return &beneficiaryHandler{repo: repo}
}

func (bh *beneficiaryHandler) GetBeneficiaries(c echo.Context) error {
	res, err := bh.repo.GetBeneficiaries(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, models.ResponseModel{Message: "beneficiaries fetched successfully", Status: "success", Data: map[string]any{
		"beneficieries": res,
	}})
}

func (bh *beneficiaryHandler) AddNewBeneficiary(c echo.Context) error {
	err := bh.repo.AddNewBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, models.ResponseModel{Message: "beneficiaries added successfully", Status: "success"})
}

func (bh *beneficiaryHandler) VerifyBeneficiary(c echo.Context) error {
	err := bh.repo.VerifyBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, models.ResponseModel{Message: "beneficiaries verification successfully", Status: "success"})
}

func (bh *beneficiaryHandler) DeleteBeneficiary(c echo.Context) error {
	err := bh.repo.DeleteBeneficiary(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{Message: err.Error(), Status: "falied"})
	}
	return c.JSON(http.StatusOK, models.ResponseModel{Message: "beneficiaries deleted successfully", Status: "success"})
}
