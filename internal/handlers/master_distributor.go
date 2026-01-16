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

func (mdh *masterDistributorHandler) GetMasterDistributorByIDRequest(c echo.Context) error {
	res, err := mdh.masterDistributorRepository.GetMasterDistributorByID(c)
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
			Message: "master distributor fetched successfully",
			Data:    map[string]any{"master_distributor": res},
		},
	)
}

func (mdh *masterDistributorHandler) ListMasterDistributorsRequest(c echo.Context) error {
	res, err := mdh.masterDistributorRepository.ListMasterDistributors(c)
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

func (mdh *masterDistributorHandler) ListMasterDistributorsByAdminIDRequest(c echo.Context) error {
	res, err := mdh.masterDistributorRepository.ListMasterDistributorsByAdminID(c)
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

func (mdh *masterDistributorHandler) UpdateMasterDistributorRequest(c echo.Context) error {
	if err := mdh.masterDistributorRepository.UpdateMasterDistributor(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "master distributor updated successfully"},
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

func (mdh *masterDistributorHandler) GetMasterDistributorsByAdminIDForDropdownRequest(
	c echo.Context,
) error {
	res, err := mdh.masterDistributorRepository.GetMasterDistributorsByAdminIDForDropdown(c)
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

func (mdh *masterDistributorHandler) LoginMasterDistributorRequest(c echo.Context) error {
	token, err := mdh.masterDistributorRepository.LoginMasterDistributor(c)
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
