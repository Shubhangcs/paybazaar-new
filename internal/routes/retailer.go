package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) RetailerRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	retRepo := repositories.NewRetailerRepository(db, jwtUtils)
	retHandler := handlers.NewRetailerHandler(retRepo)

	r.Router.POST("/retailer/login", retHandler.RetailerLoginRequest)
	rrg := r.Router.Group("/retailer", middlewares.AuthorizationMiddleware(jwtUtils))
	rrg.POST("/create", retHandler.CreateRetailerRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
	rrg.PUT("/update/details", retHandler.UpdateRetailerDetailsRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
	rrg.PUT("/update/mpin", retHandler.UpdateRetailerMPINRequest, middlewares.RequireRoles("retailer"))
	rrg.PUT("/update/block", retHandler.UpdateRetailerBlockStatusRequest, middlewares.RequireRoles("admin"))
	rrg.PUT("/update/kyc", retHandler.UpdateRetailerKYCStatusRequest, middlewares.RequireRoles("admin"))
	rrg.PUT("/update/password", retHandler.UpdateRetailerPasswordRequest, middlewares.RequireRoles("retailer"))
	rrg.PUT("/update/distributor", retHandler.UpdateRetailerDistributorRequest, middlewares.RequireRoles("admin"))
	rrg.DELETE("/delete/:retailer_id", retHandler.DeleteRetailerRequest, middlewares.RequireRoles("admin"))
	rrg.GET("/get/retailer/:retailer_id", retHandler.GetRetailerDetailsByRetailerIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor", "retailer"))
	rrg.GET("/get/md/:master_distributor_id", retHandler.GetRetailersByMasterDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor"))
	rrg.GET("/get/distributor/:distributor_id", retHandler.GetRetailersByDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
	rrg.GET("/get/dropdown/:distributor_id", retHandler.GetRetailersForDropdownByDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
}
