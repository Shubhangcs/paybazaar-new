package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type adminHandler struct {
	adminRepository repositories.AdminInterface
}

func NewAdminHandler(adminRepository repositories.AdminInterface) *adminHandler {
	return &adminHandler{
		adminRepository: adminRepository,
	}
}

func (ah *adminHandler) CreateAdminRequest(c echo.Context) error {
	if err := ah.adminRepository.CreateAdmin(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "admin created successfully"},
	)
}

func (ah *adminHandler) UpdateAdminRequest(c echo.Context) error {
	if err := ah.adminRepository.UpdateAdmin(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "admin updated successfully"},
	)
}

func (ah *adminHandler) GetAdminByIDRequest(c echo.Context) error {
	admin, err := ah.adminRepository.GetAdminByID(c)
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
			Message: "admin fetched successfully",
			Data:    map[string]any{"admin": admin},
		},
	)
}

func (ah *adminHandler) ListAdminsRequest(c echo.Context) error {
	admins, err := ah.adminRepository.ListAdmins(c)
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
			Message: "admins fetched successfully",
			Data:    map[string]any{"admins": admins},
		},
	)
}

func (ah *adminHandler) DeleteAdminRequest(c echo.Context) error {
	if err := ah.adminRepository.DeleteAdmin(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "admin deleted successfully"},
	)
}

func (ah *adminHandler) GetAdminsForDropdownRequest(c echo.Context) error {
	admins, err := ah.adminRepository.GetAdminsForDropdown(c)
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
			Message: "admins fetched successfully",
			Data:    map[string]any{"admins": admins},
		},
	)
}

func (ah *adminHandler) LoginAdminRequest(c echo.Context) error {
	token, err := ah.adminRepository.LoginAdmin(c)
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
