package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) DTHRechargeRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	dthRechargeRepo := repositories.NewDTHRechargeRepository(db)
	dthRechargeHandler := handlers.NewDTHRechargeHandler(dthRechargeRepo)

	mrrg := r.Router.Group("/dth_recharge", middlewares.AuthorizationMiddleware(jwtUtils))
	mrrg.POST("/create", dthRechargeHandler.CreateDTHRechargeRequest, middlewares.RequireRoles("retailer"))
	mrrg.GET("/get/operators", dthRechargeHandler.GetAllDTHOperatorsRequest, middlewares.RequireRoles("retailer", "admin"))
	mrrg.GET("/get/admin", dthRechargeHandler.GetAllDTHRechargesRequest, middlewares.RequireRoles("admin"))
	mrrg.GET("/get/:retailer_id", dthRechargeHandler.GetDTHRechargesByRetailerIDRequest, middlewares.RequireRoles("retailer", "admin"))
}
