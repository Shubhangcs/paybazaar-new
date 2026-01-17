package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) AdminRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	adminRepo := repositories.NewAdminRepository(db, jwtUtils)
	adminHandler := handlers.NewAdminHandler(adminRepo)

	r.Router.POST("/admin/login", adminHandler.LoginAdminRequest)
	r.Router.POST("/admin/create", adminHandler.CreateAdminRequest)
	arg := r.Router.Group("/admin", middlewares.AuthorizationMiddleware(jwtUtils))
	arg.POST("/get/all", adminHandler.ListAdminsRequest, middlewares.RequireRoles("admin"))
	arg.PUT("/update/:admin_id", adminHandler.UpdateAdminRequest, middlewares.RequireRoles("admin"))
	arg.DELETE("/delete/:admin_id", adminHandler.DeleteAdminRequest, middlewares.RequireRoles("admin"))
	arg.GET("/get/dropdown", adminHandler.GetAdminsForDropdownRequest, middlewares.RequireRoles("admin"))
	arg.GET("/get/:admin_id", adminHandler.GetAdminByIDRequest, middlewares.RequireRoles("admin"))
	arg.POST("/wallet/topup", adminHandler.AdminWalletTopupRequest, middlewares.RequireRoles("admin"))
}
