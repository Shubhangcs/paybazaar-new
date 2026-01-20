package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type masterDistributorHandler struct {
	masterDistributorRepository repositories.MasterDistributorInterface
}

func NewMasterDistributorHandler(
	masterDistributorRepository repositories.MasterDistributorInterface,
) *masterDistributorHandler {
	return &masterDistributorHandler{
		masterDistributorRepository: masterDistributorRepository,
	}
}

func (mdh *masterDistributorHandler) CreateMasterDistributorRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.CreateMasterDistributor(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor created successfully"},
	)
}

func (mdh *masterDistributorHandler) GetMasterDistributorDetailsByMasterDistributorIDRequest(c echo.Context) error {
	res, err := mdh.masterDistributorRepository.GetMasterDistributorDetailsByMasterDistributorID(c)
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
			Message: "master distributor details fetched successfully",
			Data:    map[string]any{"master_distributor": res},
		},
	)
}

func (mdh *masterDistributorHandler) GetMasterDistributorsByAdminIDRequest(c echo.Context) error {
	res, err := mdh.masterDistributorRepository.GetMasterDistributorsByAdminID(c)
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
			Message: "master distributors fetched successfully",
			Data:    map[string]any{"master_distributors": res},
		},
	)
}

func (mdh *masterDistributorHandler) GetMasterDistributorsForDropdownByAdminIDRequest(c echo.Context) error {
	res, err := mdh.masterDistributorRepository.GetMasterDistributorsForDropdownByAdminID(c)
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
			Message: "master distributors fetched successfully",
			Data:    map[string]any{"master_distributors": res},
		},
	)
}

func (mdh *masterDistributorHandler) UpdateMasterDistributorDetailsRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.UpdateMasterDistributorDetails(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor details updated successfully"},
	)
}

func (mdh *masterDistributorHandler) UpdateMasterDistributorPasswordRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.UpdateMasterDistributorPassword(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor password updated successfully"},
	)
}

func (mdh *masterDistributorHandler) UpdateMasterDistributorKYCStatusRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.UpdateMasterDistributorKYCStatus(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor kyc status updated successfully"},
	)
}

func (mdh *masterDistributorHandler) UpdateMasterDistributorMPINRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.UpdateMasterDistributorMPIN(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor mpin updated successfully"},
	)
}

func (mdh *masterDistributorHandler) UpdateMasterDistributorBlockStatusRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.UpdateMasterDistributorBlockStatus(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor block status updated successfully"},
	)
}

func (mdh *masterDistributorHandler) DeleteMasterDistributorRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.DeleteMasterDistributor(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor deleted successfully"},
	)
}

func (mdh *masterDistributorHandler) MasterDistributorLoginRequest(c echo.Context) error {
	token, err := mdh.masterDistributorRepository.MasterDistributorLogin(c)
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
