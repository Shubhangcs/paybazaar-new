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

	r.Router.POST("/admin/login", adminHandler.AdminLoginRequest)
	r.Router.POST("/admin/create", adminHandler.CreateAdminRequest)
	arg := r.Router.Group("/admin", middlewares.AuthorizationMiddleware(jwtUtils))
	arg.GET("/get/all", adminHandler.GetAllAdminsRequest, middlewares.RequireRoles("admin"))
	arg.PUT("/update/details", adminHandler.UpdateAdminDetailsRequest, middlewares.RequireRoles("admin"))
	arg.PUT("/update/password", adminHandler.UpdateAdminPasswordRequest, middlewares.RequireRoles("admin"))
	arg.PUT("/update/wallet", adminHandler.UpdateAdminWalletRequest, middlewares.RequireRoles("admin"))
	arg.PUT("/update/block_status", adminHandler.UpdateAdminBlockStatusRequest, middlewares.RequireRoles("admin"))
	arg.DELETE("/delete/:admin_id", adminHandler.DeleteAdminRequest, middlewares.RequireRoles("admin"))
	arg.GET("/get/dropdown", adminHandler.GetAdminsForDropdownRequest, middlewares.RequireRoles("admin"))
	arg.GET("/get/:admin_id", adminHandler.GetAdminDetailsByAdminIDRequest, middlewares.RequireRoles("admin"))
	arg.GET("/portal/lock", repositories.NewAdminLockRepository().LockAPI, middlewares.RequireRoles("admin"))
	arg.GET("/portal/unlock" , repositories.NewAdminLockRepository().UnlockAPI , middlewares.RequireRoles("admin"))
}
