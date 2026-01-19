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
	brg.POST("/create", bankHandler.CreateBankRequest, middlewares.RequireRoles("admin"))
	brg.POST("/create/admin", bankHandler.CreateAdminBankRequest, middlewares.RequireRoles("admin"))
	brg.GET("/get/all", bankHandler.GetAllBanksRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor", "retailer"))
	brg.GET("/get/admin/:admin_id", bankHandler.GetAdminBanksByAdminIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor", "retailer"))
	brg.DELETE("/delete/:bank_id", bankHandler.DeleteBankRequest, middlewares.RequireRoles("admin"))
	brg.DELETE("/delete/admin/:admin_bank_id", bankHandler.DeleteAdminBankRequest, middlewares.RequireRoles("admin"))
	brg.PUT("/update", bankHandler.UpdateBankDetailsRequest, middlewares.RequireRoles("admin"))
	brg.PUT("/update/admin", bankHandler.UpdateAdminBankDetailsRequest, middlewares.RequireRoles("admin"))
}
