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
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "distributor created successfully"},
	)
}

func (dh *distributorHandler) GetDistributorDetailsByDistributorIDRequest(c echo.Context) error {
	res, err := dh.distributorRepository.GetDistributorDetailsByDistributorID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "distributor fetched successfully",
		Data:    map[string]any{"distributor": res},
	})
}

func (dh *distributorHandler) GetDistributorsByAdminIDRequest(c echo.Context) error {
	res, err := dh.distributorRepository.GetDistributorsByAdminID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "distributors fetched successfully",
		Data:    map[string]any{"distributors": res},
	})
}

func (dh *distributorHandler) GetDistributorsByMasterDistributorIDRequest(c echo.Context) error {
	res, err := dh.distributorRepository.GetDistributorsByMasterDistributorID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "distributors fetched successfully",
		Data:    map[string]any{"distributors": res},
	})
}

func (dh *distributorHandler) GetDistributorsByMasterDistributorIDForDropdownRequest(
	c echo.Context,
) error {
	res, err := dh.distributorRepository.GetDistributorsForDropdownByMasterDistributorID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "distributors fetched successfully",
		Data:    map[string]any{"distributors": res},
	})
}

func (dh *distributorHandler) UpdateDistributorDetailsRequest(c echo.Context) error {
	if err := dh.distributorRepository.UpdateDistributorDetails(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "distributor updated successfully"},
	)
}

func (dh *distributorHandler) LoginDistributorRequest(c echo.Context) error {
	token, err := dh.distributorRepository.DistributorLogin(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "login successful",
		Data:    map[string]any{"access_token": token},
	})
}

func (dh *distributorHandler) UpdateDistributorBlockStatusRequest(c echo.Context) error {
	if err := dh.distributorRepository.UpdateDistributorBlockStatus(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "distributor block status updated successfully"},
	)
}

func (dh *distributorHandler) UpdateDistributorKYCStatusRequest(c echo.Context) error {
	if err := dh.distributorRepository.UpdateDistributorKYCStatus(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "distributor kyc status updated successfully"},
	)
}

func (dh *distributorHandler) UpdateDistributorMasterDistributorRequest(c echo.Context) error {
	if err := dh.distributorRepository.UpdateDistributorMasterDistributor(c); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor updated successfully"},
	)
}

func (dh *distributorHandler) UpdateDistributorPasswordRequest(
	c echo.Context,
) error {

	if err := dh.distributorRepository.UpdateDistributorPassword(c); err != nil {
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
			Message: "distributor password updated successfully",
		},
	)
}

func (dh *distributorHandler) DeleteDistributorRequest(
	c echo.Context,
) error {

	if err := dh.distributorRepository.DeleteDistributor(c); err != nil {
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
			Message: "distributor deleted successfully",
		},
	)
}

func (dh *distributorHandler) UpdateDistributorMPINRequest(
	c echo.Context,
) error {

	if err := dh.distributorRepository.UpdateDistributorMPIN(c); err != nil {
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
		},
	)
}
