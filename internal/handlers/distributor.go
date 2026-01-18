package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type distributorHandler struct {
	distributorRepository repositories.DistributorInterface
}

func NewDistributorHandler(
	distributorRepository repositories.DistributorInterface,
) *distributorHandler {
	return &distributorHandler{
		distributorRepository: distributorRepository,
	}
}

func (dh *distributorHandler) CreateDistributorRequest(c echo.Context) error {
	if err := dh.distributorRepository.CreateDistributor(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "distributor created successfully"},
	)
}

func (dh *distributorHandler) GetDistributorByIDRequest(c echo.Context) error {
	res, err := dh.distributorRepository.GetDistributorByID(c)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "distributor fetched successfully",
			Data:    map[string]any{"distributor": res},
		},
	)
}

func (dh *distributorHandler) ListDistributorsRequest(c echo.Context) error {
	res, err := dh.distributorRepository.ListDistributors(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "distributors fetched successfully",
			Data:    map[string]any{"distributors": res},
		},
	)
}

func (dh *distributorHandler) ListDistributorsByMasterDistributorIDRequest(c echo.Context) error {
	res, err := dh.distributorRepository.ListDistributorsByMasterDistributorID(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "distributors fetched successfully",
			Data:    map[string]any{"distributors": res},
		},
	)
}

func (dh *distributorHandler) UpdateDistributorRequest(c echo.Context) error {
	if err := dh.distributorRepository.UpdateDistributor(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "distributor updated successfully"},
	)
}

func (dh *distributorHandler) DeleteDistributorRequest(c echo.Context) error {
	if err := dh.distributorRepository.DeleteDistributor(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "distributor deleted successfully"},
	)
}

func (dh *distributorHandler) GetDistributorsByMasterDistributorIDForDropdownRequest(
	c echo.Context,
) error {
	res, err := dh.distributorRepository.GetDistributorsByMasterDistributorIDForDropdown(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "distributors fetched successfully",
			Data:    map[string]any{"distributors": res},
		},
	)
}

func (dh *distributorHandler) LoginDistributorRequest(c echo.Context) error {
	token, err := dh.distributorRepository.LoginDistributor(c)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "login successful",
			Data:    map[string]any{"access_token": token},
		},
	)
}

func (dh *distributorHandler) UpdateBlockStatus(
	c echo.Context,
) error {

	if err := dh.distributorRepository.UpdateBlockStatus(c); err != nil {
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
			Status:  "success",
			Message: "distributor block status updated successfully",
		},
	)
}

func (dh *distributorHandler) UpdateKYCStatus(
	c echo.Context,
) error {

	if err := dh.distributorRepository.UpdateKYCStatus(c); err != nil {
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
			Status:  "success",
			Message: "distributor KYC status updated successfully",
		},
	)
}

func (dh *distributorHandler) UpdateMPIN(
	c echo.Context,
) error {

	mpin, err := dh.distributorRepository.UpdateMPIN(c)
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
			Status:  "success",
			Message: "distributor MPIN updated successfully",
			Data: map[string]int64{
				"mpin": mpin,
			},
		},
	)
}
