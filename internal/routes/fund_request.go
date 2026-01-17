package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) FundRequestRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	fundReqRepo := repositories.NewFundRequestRepository(db)
	fundReqHandler := handlers.NewFundRequestHandler(fundReqRepo)

	frr := r.Router.Group("/fund_request", middlewares.AuthorizationMiddleware(jwtUtils))
	frr.POST("/create", fundReqHandler.CreateFundRequestRequest, middlewares.RequireRoles("master_distributor", "distributor", "retailer"))
	frr.GET("/get/all", fundReqHandler.GetAllFundRequestsRequest, middlewares.RequireRoles("admin"))
	frr.POST("/get/requester", fundReqHandler.GetFundRequestsByRequesterIDRequest, middlewares.RequireRoles("master_distributor", "distributor", "retailer"))
	frr.POST("/get/request_to", fundReqHandler.GetFundRequestsByRequestToIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
	frr.GET("/get/:fund_request_id", fundReqHandler.GetFundRequestByIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor", "retailer"))
	frr.PUT("/accept/:fund_request_id", fundReqHandler.AcceptFundRequestRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
	frr.PUT("/reject/:fund_request_id", fundReqHandler.RejectFundRequestRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
}
