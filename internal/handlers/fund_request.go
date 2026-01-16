package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type fundRequestHandler struct {
	fundRequestRepository repositories.FundRequestInterface
}

func NewFundRequestHandler(
	fundRequestRepository repositories.FundRequestInterface,
) *fundRequestHandler {
	return &fundRequestHandler{
		fundRequestRepository: fundRequestRepository,
	}
}

func (fh *fundRequestHandler) CreateFundRequestRequest(
	c echo.Context,
) error {

	id, err := fh.fundRequestRepository.CreateFundRequest(c)
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
			Message: "fund request created successfully",
			Data:    map[string]any{"fund_request_id": id},
		},
	)
}

func (fh *fundRequestHandler) GetFundRequestByIDRequest(
	c echo.Context,
) error {

	res, err := fh.fundRequestRepository.GetFundRequestByID(c)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
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
			Message: "fund request fetched successfully",
			Data:    map[string]any{"fund_request": res},
		},
	)
}

func (fh *fundRequestHandler) GetAllFundRequestsRequest(
	c echo.Context,
) error {

	res, err := fh.fundRequestRepository.GetAllFundRequests(c)
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
			Message: "fund requests fetched successfully",
			Data:    map[string]any{"fund_requests": res},
		},
	)
}

func (fh *fundRequestHandler) GetFundRequestsByRequesterIDRequest(
	c echo.Context,
) error {

	res, err := fh.fundRequestRepository.GetFundRequestsByRequesterID(c)
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
			Message: "fund requests fetched successfully",
			Data:    map[string]any{"fund_requests": res},
		},
	)
}

func (fh *fundRequestHandler) GetFundRequestsByRequestToIDRequest(
	c echo.Context,
) error {

	res, err := fh.fundRequestRepository.GetFundRequestsByRequestToID(c)
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
			Message: "fund requests fetched successfully",
			Data:    map[string]any{"fund_requests": res},
		},
	)
}

func (fh *fundRequestHandler) AcceptFundRequestRequest(
	c echo.Context,
) error {

	if err := fh.fundRequestRepository.AcceptFundRequest(c); err != nil {
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
			Message: "fund request accepted successfully",
		},
	)
}

func (fh *fundRequestHandler) RejectFundRequestRequest(
	c echo.Context,
) error {

	if err := fh.fundRequestRepository.RejectFundRequest(c); err != nil {
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
			Message: "fund request rejected successfully",
		},
	)
}
