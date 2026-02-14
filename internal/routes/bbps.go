package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) BBPSRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	bbpsRepo := repositories.NewBBPSRepository(db)
	bbpsHandler := handlers.NewBBPSHandler(bbpsRepo)

	bbpsrg := r.Router.Group("/bbps", middlewares.AuthorizationMiddleware(jwtUtils))

	bbpsrg.POST("/create/postpaid", bbpsHandler.CreatePostpaidMobileRechargeRequest, middlewares.RequireRoles("retailer"))
	bbpsrg.POST("/get/postpaid/balance", bbpsHandler.GetPostpaidMobileRechargeBalanceRequest, middlewares.RequireRoles("retailer"))
	bbpsrg.GET("/recharge/get/all", bbpsHandler.GetAllPostpaidMobileRechargeRequest, middlewares.RequireRoles("admin"))
	bbpsrg.GET("/recharge/get/:retailer_id", bbpsHandler.GetPostpaidMobileRechargeByRetailerIDRequest)
	bbpsrg.POST("/create/electricity", bbpsHandler.CreateElectricityBillPaymentRequest, middlewares.RequireRoles("retailer"))
	bbpsrg.GET("/get/electricity/operators", bbpsHandler.GetAllElectricityBillOperatorsRequest, middlewares.RequireRoles("retailer"))
	bbpsrg.POST("/get/electricity/balance", bbpsHandler.GetElectricityBillBalanceRequest, middlewares.RequireRoles("retailer"))
	bbpsrg.GET("/get/all/electricity/transactions", bbpsHandler.GetAllElectricityBillHistoryRequest, middlewares.RequireRoles("admin"))
	bbpsrg.GET("/get/electricity/transactions/:retailer_id", bbpsHandler.GetElectricityBillHistoryByRetailerIDRequest, middlewares.RequireRoles("retailer", "admin"))
	bbpsrg.GET("/electricity/transaction/refund/:transaction_id", bbpsHandler.ElectricityBillPaymentRefundRequest, middlewares.RequireRoles("admin"))
	bbpsrg.GET("/postpaid/recharge/transaction/refund/:transaction_id", bbpsHandler.MobileRechargePostpaidRefundRequest, middlewares.RequireRoles("admin"))
}
