package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) MobileRechargeRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	mobileRechargeRepo := repositories.NewMobileRechargeRepository(db)
	mobileRechargeHandler := handlers.NewMobileRechargeHandler(mobileRechargeRepo)

	mrrg := r.Router.Group("/mobile_recharge", middlewares.AuthorizationMiddleware(jwtUtils))
	mrrg.POST("/create", mobileRechargeHandler.CreateMobileRechargeRequest, middlewares.RequireRoles("retailer"))
	mrrg.GET("/get/operators", mobileRechargeHandler.GetMobileRechargeOperatorsRequest, middlewares.RequireRoles("retailer", "admin"))
	mrrg.GET("/get/circle", mobileRechargeHandler.GetMobileRechargeCirclesRequest, middlewares.RequireRoles("retailer", "admin"))
	mrrg.GET("/get/admin", mobileRechargeHandler.GetAllMobileRechargesRequest, middlewares.RequireRoles("admin"))
	mrrg.POST("/get/plans", mobileRechargeHandler.GetMobileRechargePlansRequest, middlewares.RequireRoles("admin", "retailer"))
	mrrg.GET("/get/:retailer_id", mobileRechargeHandler.GetMobileRechargesByRetailerIDRequest, middlewares.RequireRoles("retailer", "admin"))
}
