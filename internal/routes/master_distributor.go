package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) MasterDistributorRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	mdRepo := repositories.NewMasterDistributorRepository(db, jwtUtils)
	mdHandler := handlers.NewMasterDistributorHandler(mdRepo)

	r.Router.POST("/md/login", mdHandler.MasterDistributorLoginRequest)
	mdrg := r.Router.Group("/md", middlewares.AuthorizationMiddleware(jwtUtils))
	mdrg.POST("/create", mdHandler.CreateMasterDistributorRequest, middlewares.RequireRoles("admin"))
	mdrg.PUT("/update/details", mdHandler.UpdateMasterDistributorDetailsRequest, middlewares.RequireRoles("admin"))
	mdrg.PUT("/update/password", mdHandler.UpdateMasterDistributorPasswordRequest, middlewares.RequireRoles("admin", "master_distributor"))
	mdrg.PUT("/update/kyc", mdHandler.UpdateMasterDistributorKYCStatusRequest, middlewares.RequireRoles("admin"))
	mdrg.PUT("/update/block", mdHandler.UpdateMasterDistributorBlockStatusRequest, middlewares.RequireRoles("admin"))
	mdrg.PUT("/update/mpin", mdHandler.UpdateMasterDistributorMPINRequest, middlewares.RequireRoles("admin", "master_distributor"))
	mdrg.DELETE("/delete/:master_distributor_id", mdHandler.DeleteMasterDistributorRequest, middlewares.RequireRoles("admin"))
	mdrg.GET("/get/md/:master_distributor_id", mdHandler.GetMasterDistributorDetailsByMasterDistributorIDRequest, middlewares.RequireRoles("master_distributor", "admin"))
	mdrg.GET("/get/admin/:admin_id", mdHandler.GetMasterDistributorsByAdminIDRequest, middlewares.RequireRoles("admin"))
	mdrg.GET("/get/dropdown/:admin_id", mdHandler.GetMasterDistributorsForDropdownByAdminIDRequest, middlewares.RequireRoles("admin"))
}
