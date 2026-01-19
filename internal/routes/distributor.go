package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) DistributorRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	disRepo := repositories.NewDistributorRepository(db, jwtUtils)
	disHandler := handlers.NewDistributorHandler(disRepo)

	r.Router.POST("/distributor/login", disHandler.LoginDistributorRequest)
	drg := r.Router.Group("/distributor", middlewares.AuthorizationMiddleware(jwtUtils))
	drg.POST("/create", disHandler.CreateDistributorRequest, middlewares.RequireRoles("admin", "master_distributor"))
	drg.PUT("/update", disHandler.UpdateDistributorDetailsRequest, middlewares.RequireRoles("admin", "distributor"))
	drg.DELETE("/delete/:distributor_id", disHandler.DeleteDistributorRequest, middlewares.RequireRoles("admin"))
	drg.GET("/get/admin/:admin_id", disHandler.GetDistributorsByAdminIDRequest, middlewares.RequireRoles("admin"))
	drg.GET("/get/distributor/:distributor_id", disHandler.GetDistributorDetailsByDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
	drg.GET("/get/md/:master_distributor_id", disHandler.GetDistributorsByMasterDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor"))
	drg.GET("/get/dropdown/:master_distributor_id", disHandler.GetDistributorsByMasterDistributorIDForDropdownRequest, middlewares.RequireRoles("admin", "master_distributor"))
	drg.PUT("/update/kyc/:distributor_id", disHandler.UpdateDistributorKYCStatusRequest, middlewares.RequireRoles("admin"))
	drg.PUT("/update/block/:distributor_id", disHandler.UpdateDistributorBlockStatusRequest, middlewares.RequireRoles("admin"))
	drg.PUT("/update/mpin/:distributor_id", disHandler.UpdateDistributorMPINRequest, middlewares.RequireRoles("distributor"))
}
