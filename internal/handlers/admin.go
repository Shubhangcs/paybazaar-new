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

func (ah *adminHandler) UpdateAdminDetailsRequest(c echo.Context) error {
	if err := ah.adminRepository.UpdateAdminDetails(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "admin details updated successfully"},
	)
}

func (ah *adminHandler) UpdateAdminPasswordRequest(c echo.Context) error {
	if err := ah.adminRepository.UpdateAdminPassword(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "admin password updated successfully"},
	)
}

func (ah *adminHandler) UpdateAdminWalletRequest(c echo.Context) error {
	if err := ah.adminRepository.UpdateAdminWallet(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "admin wallet updated successfully"},
	)
}

func (ah *adminHandler) UpdateAdminBlockStatusRequest(c echo.Context) error {
	if err := ah.adminRepository.UpdateAdminBlockStatus(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "admin block status updated successfully"},
	)
}

func (ah *adminHandler) GetAdminDetailsByAdminIDRequest(c echo.Context) error {
	admin, err := ah.adminRepository.GetAdminDetailsByAdminID(c)
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
			Message: "admin details fetched successfully",
			Data:    map[string]any{"admin": admin},
		},
	)
}

func (ah *adminHandler) GetAllAdminsRequest(c echo.Context) error {
	admins, err := ah.adminRepository.GetAllAdmins(c)
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
			Message: "dropdown admins fetched successfully",
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

func (ah *adminHandler) AdminLoginRequest(c echo.Context) error {
	token, err := ah.adminRepository.AdminLogin(c)
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
