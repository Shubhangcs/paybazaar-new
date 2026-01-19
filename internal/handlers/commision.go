package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type commisionHandler struct {
	commisionRepo repositories.CommisionInterface
}

func NewCommisionHandler(
	commisionRepo repositories.CommisionInterface,
) *commisionHandler {
	return &commisionHandler{
		commisionRepo: commisionRepo,
	}
}

func (ch *commisionHandler) CreateCommisionRequest(c echo.Context) error {
	err := ch.commisionRepo.CreateCommision(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision created successfully",
	})
}

func (ch *commisionHandler) GetCommisionDetailsByCommisionIDRequest(c echo.Context) error {
	data, err := ch.commisionRepo.GetCommisionDetailsByCommisionID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision fetched successfully",
		Data: map[string]any{
			"commision": data,
		},
	})
}

func (ch *commisionHandler) GetCommisionsByUserIDRequest(c echo.Context) error {
	data, err := ch.commisionRepo.GetCommisionsByUserID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commisions fetched successfully",
		Data: map[string]any{
			"commisions": data,
		},
	})
}

func (ch *commisionHandler) GetCommisionByUserIDAndServiceRequest(c echo.Context) error {
	data, err := ch.commisionRepo.GetCommisionByUserIDAndService(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision fetched successfully",
		Data: map[string]any{
			"commision": data,
		},
	})
}

func (ch *commisionHandler) UpdateCommisionDetailsRequest(c echo.Context) error {
	if err := ch.commisionRepo.UpdateCommisionDetails(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision details updated successfully",
	})
}

func (ch *commisionHandler) DeleteCommisionRequest(c echo.Context) error {
	if err := ch.commisionRepo.DeleteCommision(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision deleted successfully",
	})
}
