package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type fundTransferHandler struct {
	repo repositories.FundTransferInterface
}

func NewFundTransferHandler(repo repositories.FundTransferInterface) *fundTransferHandler {
	return &fundTransferHandler{repo: repo}
}

func (fh *fundTransferHandler) CreateFundTransfer(c echo.Context) error {
	if err := fh.repo.CreateFundTransfer(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "fund transfer successful",
	})
}

func (fh *fundTransferHandler) GetFundTransfersByFromID(c echo.Context) error {
	data, err := fh.repo.GetFundTransfersByFromID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "fund transfers fetched successfully",
		Data:    data,
	})
}

func (fh *fundTransferHandler) GetFundTransfersByToID(c echo.Context) error {
	data, err := fh.repo.GetFundTransfersByToID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "fund transfers fetched successfully",
		Data:    data,
	})
}
