package models

import "time"

type CreateTicketModel struct {
	AdminID           string `json:"admin_id" validate:"required"`
	UserID            string `json:"user_id" validate:"required"`
	TicketTitle       string `json:"ticket_title" validate:"required,min=3,max=200"`
	TicketDescription string `json:"ticket_description" validate:"required,min=5"`
}

type UpdateTicketModel struct {
	TicketTitle       *string `json:"ticket_title,omitempty" validate:"omitempty,min=3,max=200"`
	TicketDescription *string `json:"ticket_description,omitempty" validate:"omitempty,min=5"`
}

type UpdateTicketStatusModel struct {
	IsTicketCleared bool `json:"is_ticket_cleared" validate:"required"`
}

type TicketResponseModel struct {
	TicketID          int64     `json:"ticket_id"`
	AdminID           string    `json:"admin_id"`
	UserID            string    `json:"user_id"`
	TicketTitle       string    `json:"ticket_title"`
	TicketDescription string    `json:"ticket_description"`
	IsTicketCleared   bool      `json:"is_ticket_cleared"`
	CreatedAt         time.Time `json:"created_at"`
}

type TicketFilterRequestModel struct {
	AdminID   *string    `json:"admin_id,omitempty"`
	UserID    *string    `json:"user_id,omitempty"`
	IsCleared *bool      `json:"is_ticket_cleared,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}
