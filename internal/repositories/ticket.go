package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type TicketInterface interface {
	CreateTicket(echo.Context) (int64, error)
	GetTicketByID(echo.Context) (*models.TicketResponseModel, error)
	GetTicketsByAdminID(echo.Context) ([]models.TicketResponseModel, error)
	GetTicketsByUserID(echo.Context) ([]models.TicketResponseModel, error)
	GetAllTickets(echo.Context) ([]models.TicketResponseModel, error)
	UpdateTicket(echo.Context) error
	UpdateTicketStatus(echo.Context) error
	DeleteTicket(echo.Context) error
}

type ticketRepository struct {
	db *database.Database
}

func NewTicketRepository(db *database.Database) *ticketRepository {
	return &ticketRepository{db: db}
}

func (tr *ticketRepository) CreateTicket(c echo.Context) (int64, error) {
	var req models.CreateTicketModel
	if err := bindAndValidate(c, &req); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.CreateTicketQuery(ctx, req)
}

func (tr *ticketRepository) GetTicketByID(c echo.Context) (*models.TicketResponseModel, error) {
	idStr := c.Param("ticket_id")
	ticketID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ticket_id")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.GetTicketByIDQuery(ctx, ticketID)
}

func (tr *ticketRepository) GetTicketsByAdminID(c echo.Context) ([]models.TicketResponseModel, error) {
	adminID := c.Param("admin_id")
	if adminID == "" {
		return nil, fmt.Errorf("admin_id is required")
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.GetTicketsByAdminIDQuery(ctx, adminID, limit, offset)
}

func (tr *ticketRepository) GetTicketsByUserID(c echo.Context) ([]models.TicketResponseModel, error) {
	userID := c.Param("user_id")
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.GetTicketsByUserIDQuery(ctx, userID, limit, offset)
}

func (tr *ticketRepository) GetAllTickets(c echo.Context) ([]models.TicketResponseModel, error) {
	limit, offset := parsePagination(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.GetAllTicketsQuery(ctx, limit, offset)
}

func (tr *ticketRepository) UpdateTicket(c echo.Context) error {
	idStr := c.Param("ticket_id")
	ticketID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ticket_id")
	}

	var req models.UpdateTicketModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.UpdateTicketQuery(ctx, ticketID, req)
}

func (tr *ticketRepository) UpdateTicketStatus(c echo.Context) error {
	idStr := c.Param("ticket_id")
	ticketID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ticket_id")
	}

	var req models.UpdateTicketStatusModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.UpdateTicketStatusQuery(ctx, ticketID, req.IsTicketCleared)
}

func (tr *ticketRepository) DeleteTicket(c echo.Context) error {
	idStr := c.Param("ticket_id")
	ticketID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ticket_id")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	return tr.db.DeleteTicketQuery(ctx, ticketID)
}
