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

func (ch *commisionHandler) CreateCommision(c echo.Context) error {
	id, err := ch.commisionRepo.CreateCommision(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision created successfully",
		Data: map[string]int64{
			"commision_id": id,
		},
	})
}

func (ch *commisionHandler) GetCommisionByID(c echo.Context) error {
	data, err := ch.commisionRepo.GetCommisionByID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision fetched successfully",
		Data:    data,
	})
}

func (ch *commisionHandler) GetCommisionByUserID(c echo.Context) error {
	data, err := ch.commisionRepo.GetCommisionByUserID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision fetched successfully",
		Data:    data,
	})
}

func (ch *commisionHandler) GetAllCommisions(c echo.Context) error {
	data, err := ch.commisionRepo.GetAllCommisions(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commisions fetched successfully",
		Data:    data,
	})
}

func (ch *commisionHandler) UpdateCommision(c echo.Context) error {
	if err := ch.commisionRepo.UpdateCommision(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "commision updated successfully",
	})
}

func (ch *commisionHandler) DeleteCommision(c echo.Context) error {
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
