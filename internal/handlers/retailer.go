package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type retailerHandler struct {
	retailerRepository repositories.RetailerInterface
}

func NewRetailerHandler(
	retailerRepository repositories.RetailerInterface,
) *retailerHandler {
	return &retailerHandler{
		retailerRepository: retailerRepository,
	}
}

func (rh *retailerHandler) CreateRetailerRequest(c echo.Context) error {
	if err := rh.retailerRepository.CreateRetailer(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer created successfully"},
	)
}

func (rh *retailerHandler) GetRetailerDetailsByRetailerIDRequest(c echo.Context) error {
	res, err := rh.retailerRepository.GetRetailerDetailsByRetailerID(c)
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
			Message: "retailer fetched successfully",
			Data:    map[string]any{"retailer": res},
		},
	)
}

func (rh *retailerHandler) GetRetailersByAdminIDRequest(c echo.Context) error {
	res, err := rh.retailerRepository.GetRetailersByAdminID(c)
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
			Message: "retailers fetched successfully",
			Data:    map[string]any{"retailers": res},
		},
	)
}

func (rh *retailerHandler) GetRetailersByMasterDistributorIDRequest(c echo.Context) error {
	res, err := rh.retailerRepository.GetRetailersByMasterDistributorID(c)
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
			Message: "retailers fetched successfully",
			Data:    map[string]any{"retailers": res},
		},
	)
}

func (rh *retailerHandler) GetRetailersByDistributorIDRequest(c echo.Context) error {
	res, err := rh.retailerRepository.GetRetailersByDistributorID(c)
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
			Message: "retailers fetched successfully",
			Data:    map[string]any{"retailers": res},
		},
	)
}

func (rh *retailerHandler) GetRetailersForDropdownByDistributorIDRequest(c echo.Context) error {
	res, err := rh.retailerRepository.GetRetailersForDropdownByDistributorID(c)
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
			Message: "retailers fetched successfully",
			Data:    map[string]any{"retailers": res},
		},
	)
}

func (rh *retailerHandler) UpdateRetailerDetailsRequest(c echo.Context) error {
	if err := rh.retailerRepository.UpdateRetailerDetails(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer details updated successfully"},
	)
}

func (rh *retailerHandler) UpdateRetailerKYCStatusRequest(c echo.Context) error {
	if err := rh.retailerRepository.UpdateRetailerKYCStatus(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer kyc status updated successfully"},
	)
}

func (rh *retailerHandler) UpdateRetailerBlockStatusRequest(c echo.Context) error {
	if err := rh.retailerRepository.UpdateRetailerBlockStatus(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer block status updated successfully"},
	)
}

func (rh *retailerHandler) UpdateRetailerPasswordRequest(c echo.Context) error {
	if err := rh.retailerRepository.UpdateRetailerPassword(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer password updated successfully"},
	)
}

func (rh *retailerHandler) UpdateRetailerMPINRequest(c echo.Context) error {
	if err := rh.retailerRepository.UpdateRetailerMPIN(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer mpin updated successfully"},
	)
}

func (rh *retailerHandler) UpdateRetailerDistributorRequest(c echo.Context) error {
	if err := rh.retailerRepository.UpdateRetailerDistributor(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer distributor updated successfully"},
	)
}

func (rh *retailerHandler) DeleteRetailerRequest(c echo.Context) error {
	if err := rh.retailerRepository.DeleteRetailer(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "retailer deleted successfully"},
	)
}

func (rh *retailerHandler) RetailerLoginRequest(c echo.Context) error {
	token, err := rh.retailerRepository.RetailerLogin(c)
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
