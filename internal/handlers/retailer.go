// package handlers

// import (
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// 	"github.com/levion-studio/paybazaar/internal/models"
// 	"github.com/levion-studio/paybazaar/internal/repositories"
// )

// type retailerHandler struct {
// 	retailerRepository repositories.RetailerInterface
// }

// func NewRetailerHandler(
// 	retailerRepository repositories.RetailerInterface,
// ) *retailerHandler {
// 	return &retailerHandler{
// 		retailerRepository: retailerRepository,
// 	}
// }

// func (rh *retailerHandler) CreateRetailerRequest(c echo.Context) error {
// 	if err := rh.retailerRepository.CreateRetailer(c); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{Status: "success", Message: "retailer created successfully"},
// 	)
// }

// func (rh *retailerHandler) GetRetailerByIDRequest(c echo.Context) error {
// 	res, err := rh.retailerRepository.GetRetailerByID(c)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusNotFound,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailer fetched successfully",
// 			Data:    map[string]any{"retailer": res},
// 		},
// 	)
// }

// func (rh *retailerHandler) ListRetailersRequest(c echo.Context) error {
// 	res, err := rh.retailerRepository.ListRetailers(c)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailers fetched successfully",
// 			Data:    map[string]any{"retailers": res},
// 		},
// 	)
// }

// func (rh *retailerHandler) ListRetailersByDistributorIDRequest(c echo.Context) error {
// 	res, err := rh.retailerRepository.ListRetailersByDistributorID(c)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailers fetched successfully",
// 			Data:    map[string]any{"retailers": res},
// 		},
// 	)
// }

// func (rh *retailerHandler) ListRetailersByMasterDistributorIDRequest(c echo.Context) error {
// 	res, err := rh.retailerRepository.ListRetailersByMasterDistributorID(c)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailers fetched successfully",
// 			Data:    map[string]any{"retailers": res},
// 		},
// 	)
// }

// func (rh *retailerHandler) UpdateRetailerRequest(c echo.Context) error {
// 	if err := rh.retailerRepository.UpdateRetailer(c); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{Status: "success", Message: "retailer updated successfully"},
// 	)
// }

// func (rh *retailerHandler) DeleteRetailerRequest(c echo.Context) error {
// 	if err := rh.retailerRepository.DeleteRetailer(c); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{Status: "success", Message: "retailer deleted successfully"},
// 	)
// }

// func (rh *retailerHandler) GetRetailersByDistributorIDForDropdownRequest(
// 	c echo.Context,
// ) error {
// 	res, err := rh.retailerRepository.GetRetailersByDistributorIDForDropdown(c)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailers fetched successfully",
// 			Data:    map[string]any{"retailers": res},
// 		},
// 	)
// }

// func (rh *retailerHandler) LoginRetailerRequest(c echo.Context) error {
// 	token, err := rh.retailerRepository.LoginRetailer(c)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusUnauthorized,
// 			models.ResponseModel{Status: "failed", Message: err.Error()},
// 		)
// 	}
// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "login successful",
// 			Data:    map[string]any{"access_token": token},
// 		},
// 	)
// }

// func (rh *retailerHandler) UpdateBlockStatus(
// 	c echo.Context,
// ) error {

// 	if err := rh.retailerRepository.UpdateBlockStatus(c); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{
// 				Status:  "failed",
// 				Message: err.Error(),
// 			},
// 		)
// 	}

// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailer block status updated successfully",
// 		},
// 	)
// }

// func (rh *retailerHandler) UpdateKYCStatus(
// 	c echo.Context,
// ) error {

// 	if err := rh.retailerRepository.UpdateKYCStatus(c); err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{
// 				Status:  "failed",
// 				Message: err.Error(),
// 			},
// 		)
// 	}

// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailer KYC status updated successfully",
// 		},
// 	)
// }

// func (rh *retailerHandler) UpdateMPIN(
// 	c echo.Context,
// ) error {

// 	mpin, err := rh.retailerRepository.UpdateMPIN(c)
// 	if err != nil {
// 		return c.JSON(
// 			http.StatusBadRequest,
// 			models.ResponseModel{
// 				Status:  "failed",
// 				Message: err.Error(),
// 			},
// 		)
// 	}

// 	return c.JSON(
// 		http.StatusOK,
// 		models.ResponseModel{
// 			Status:  "success",
// 			Message: "retailer MPIN updated successfully",
// 			Data: map[string]int64{
// 				"mpin": mpin,
// 			},
// 		},
// 	)
// }
package handlers