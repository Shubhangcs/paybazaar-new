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

	bbpsrg.POST("/create/postpaid", bbpsHandler.CreatePostpaidMobileRechargeRequest)
	bbpsrg.POST("/get/postpaid/balance", bbpsHandler.GetPostpaidMobileRechargeBalanceRequest)
}
