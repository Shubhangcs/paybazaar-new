package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type revertHandler struct {
	revertRepository repositories.RevertInterface
}

func NewRevertHandler(revertRepository repositories.RevertInterface) *revertHandler {
	return &revertHandler{
		revertRepository,
	}
}

func (rh *revertHandler) CreateRevertRequest(c echo.Context) error {
	if err := rh.revertRepository.CreateRevert(c); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			models.ResponseModel{Status: "failed", Message: err.Error()},
		)
	}
	return c.JSON(
		http.StatusOK,
		models.ResponseModel{Status: "success", Message: "revert successful"},
	)
}

func (rh *revertHandler) GetRevertsByFromID(c echo.Context) error {
	data, err := rh.revertRepository.GetRevertsByFromID(c)
	if err != nil {
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
			Message: "revert transactions fetched",
			Data:    data,
		},
	)
}

func (rh *revertHandler) GetRevertsByOnID(c echo.Context) error {
	data, err := rh.revertRepository.GetRevertsByOnID(c)
	if err != nil {
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
			Message: "revert transactions fetched",
			Data:    data,
		},
	)
}
