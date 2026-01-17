package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) BankRouter(db *database.Database, jwtUtils *pkg.JwtUtils) {
	bankRepo := repositories.NewBankRepository(db)
	bankHandler := handlers.NewBankHandler(bankRepo)

	brg := r.Router.Group("/bank", middlewares.AuthorizationMiddleware(jwtUtils))
	brg.POST("/create", bankHandler.CreateBank, middlewares.RequireRoles("admin"))
	brg.POST("/create/admin", bankHandler.CreateAdminBank, middlewares.RequireRoles("admin"))
	brg.GET("/get/all", bankHandler.GetAllBanks)
	brg.GET("/get/admin/:admin_id", bankHandler.GetAdminBanksByAdminID)
	brg.DELETE("/delete/:bank_id", bankHandler.DeleteBank, middlewares.RequireRoles("admin"))
	brg.DELETE("/delete/admin/:admin_bank_id", bankHandler.DeleteAdminBank, middlewares.RequireRoles("admin"))
	brg.PUT("/update/:bank_id", bankHandler.UpdateBank, middlewares.RequireRoles("admin"))
	brg.PUT("/update/admin/:admin_bank_id", bankHandler.UpdateAdminBank, middlewares.RequireRoles("admin"))
}
