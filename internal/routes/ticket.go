package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) TicketRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	ticketRepo := repositories.NewTicketRepository(db)
	ticketHandler := handlers.NewTicketHandler(ticketRepo)

	trr := r.Router.Group("/ticket", middlewares.AuthorizationMiddleware(jwtUtils))

	trr.POST("/create", ticketHandler.CreateTicket)
	trr.GET("/get/:ticket_id", ticketHandler.GetTicketByID, middlewares.RequireRoles("admin"))
	trr.GET("/admin/:admin_id", ticketHandler.GetTicketsByAdminID, middlewares.RequireRoles("admin"))
	trr.GET("/user/:user_id", ticketHandler.GetTicketsByUserID, middlewares.RequireRoles("admin"))
	trr.GET("/get/all", ticketHandler.GetAllTickets, middlewares.RequireRoles("admin"))
	trr.PUT("/update/:ticket_id", ticketHandler.UpdateTicket)
	trr.PUT("/update/:ticket_id/status", ticketHandler.UpdateTicketStatus, middlewares.RequireRoles("admin"))
	trr.DELETE("/delete/:ticket_id", ticketHandler.DeleteTicket, middlewares.RequireRoles("admin"))
}
