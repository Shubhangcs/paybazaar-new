package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type limitHandler struct {
	limitRepository repositories.LimitInterface
}

func NewLimitHandler(limitRepository repositories.LimitInterface) *limitHandler {
	return &limitHandler{
		limitRepository,
	}
}

func (lh *limitHandler) CreateLimitRequest(c echo.Context) error {
	if err := lh.limitRepository.CreateLimit(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "limit created successfully",
		},
	)
}

func (lh *limitHandler) UpdateLimitRequest(c echo.Context) error {
	if err := lh.limitRepository.UpdateLimit(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "limit updated successfully",
		},
	)
}

func (lh *limitHandler) DeleteLimitRequest(c echo.Context) error {
	if err := lh.limitRepository.DeleteLimit(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "limit deleted successfully",
		},
	)
}

func (lh *limitHandler) GetAllLimitsRequest(c echo.Context) error {
	res , err := lh.limitRepository.GetAllLimits(c); 
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{
			Status:  "success",
			Message: "limit deleted successfully",
			Data: map[string]any{"limits": res},
		},
	)
}
