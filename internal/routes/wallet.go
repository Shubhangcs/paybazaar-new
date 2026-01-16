package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) WalletTransactionRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	walletRepo := repositories.NewWalletTransactionRepository(db)
	walletHandler := handlers.NewWalletTransactionHandler(walletRepo)

	wtr := r.Router.Group("/wallet", middlewares.AuthorizationMiddleware(jwtUtils))
	wtr.POST("/create", walletHandler.CreateWalletTransactionRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor", "retailer"))
	wtr.GET("/get/balance/admin/:admin_id", walletHandler.GetAdminWalletBalanceRequest, middlewares.RequireRoles("admin"))
	wtr.GET("/get/balance/md/:master_distributor_id", walletHandler.GetMasterDistributorWalletBalanceRequest, middlewares.RequireRoles("master_distributor"))
	wtr.GET("/get/balance/distributor/:distributor_id", walletHandler.GetDistributorWalletBalanceRequest, middlewares.RequireRoles("distributor"))
	wtr.GET("/get/balance/retailer/:retailer_id", walletHandler.GetRetailerWalletBalanceRequest, middlewares.RequireRoles("retailer"))
	wtr.GET("/get/transactions/admin/:admin_id", walletHandler.GetAdminWalletTransactionsRequest, middlewares.RequireRoles("admin"))
	wtr.GET("/get/transactions/md/:master_distributor_id", walletHandler.GetMasterDistributorWalletTransactionsRequest, middlewares.RequireRoles("master_distributor"))
	wtr.GET("/get/transactions/distributor/:distributor_id", walletHandler.GetDistributorWalletTransactionsRequest, middlewares.RequireRoles("distributor"))
	wtr.GET("/get/transaction/retailer/:retailer_id", walletHandler.GetRetailerWalletTransactionsRequest, middlewares.RequireRoles("retailer"))
}
