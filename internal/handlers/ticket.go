package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/models"
	"github.com/levion-studio/paybazaar/internal/repositories"
)

type ticketHandler struct {
	repo repositories.TicketInterface
}

func NewTicketHandler(repo repositories.TicketInterface) *ticketHandler {
	return &ticketHandler{repo: repo}
}

func (th *ticketHandler) CreateTicket(c echo.Context) error {
	id, err := th.repo.CreateTicket(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "ticket created successfully",
		Data: map[string]int64{
			"ticket_id": id,
		},
	})
}

func (th *ticketHandler) GetTicketByID(c echo.Context) error {
	ticket, err := th.repo.GetTicketByID(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "ticket fetched successfully",
		Data:    ticket,
	})
}

func (th *ticketHandler) GetTicketsByAdminID(c echo.Context) error {
	tickets, err := th.repo.GetTicketsByAdminID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "tickets fetched successfully",
		Data:    tickets,
	})
}

func (th *ticketHandler) GetTicketsByUserID(c echo.Context) error {
	tickets, err := th.repo.GetTicketsByUserID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "tickets fetched successfully",
		Data:    tickets,
	})
}

func (th *ticketHandler) GetAllTickets(c echo.Context) error {
	tickets, err := th.repo.GetAllTickets(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "tickets fetched successfully",
		Data:    tickets,
	})
}

func (th *ticketHandler) UpdateTicket(c echo.Context) error {
	if err := th.repo.UpdateTicket(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "ticket updated successfully",
	})
}

func (th *ticketHandler) UpdateTicketStatus(c echo.Context) error {
	if err := th.repo.UpdateTicketStatus(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "ticket status updated successfully",
	})
}

func (th *ticketHandler) DeleteTicket(c echo.Context) error {
	if err := th.repo.DeleteTicket(c); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "failed",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "success",
		Message: "ticket deleted successfully",
	})
}
