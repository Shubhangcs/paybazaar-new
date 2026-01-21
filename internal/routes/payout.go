package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) PayoutRoutes(
	db *database.Database,
	jwtUtils *pkg.JwtUtils,
) {

	payoutRepo := repositories.NewPayoutRepository(db, jwtUtils)
	payoutHandler := handlers.NewPayoutHandler(payoutRepo)
	pr := r.Router.Group(
		"/payout",
		middlewares.AuthorizationMiddleware(jwtUtils),
	)
	pr.POST("/create", payoutHandler.CreatePayout, middlewares.RequireRoles("retailer"))
	pr.GET("/get", payoutHandler.GetAllPayoutsRequest, middlewares.RequireRoles("admin"))
	pr.GET("/get/:retailer_id", payoutHandler.GetPayoutsByRetailerIDRequest, middlewares.RequireRoles("retailer", "admin"))
	pr.GET("/get/ledger/:retailer_id", payoutHandler.GetRetailerPayoutLedgerWithWalletRequest, middlewares.RequireRoles("retailer"))
}
