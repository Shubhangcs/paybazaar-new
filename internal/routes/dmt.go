package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) DMTRoutes(db *database.Database, jutUtils *pkg.JwtUtils) {
	dmtRepo := repositories.NewDMTRepository(db)
	dmtHandler := handlers.NewDMTHandler(dmtRepo)

	drg := r.Router.Group("/dmt", middlewares.AuthorizationMiddleware(jutUtils))
	drg.POST("/check/wallet", dmtHandler.CheckDMTWalletExistsRequest)
	drg.POST("/create/wallet", dmtHandler.CreateDMTWalletRequest)
	drg.POST("/verify/wallet", dmtHandler.VerifyDMTWalletRequest)
	drg.POST("/add/beneficiary", dmtHandler.AddDMTBeneficiaryRequest)
	drg.GET("/get/banks", dmtHandler.GetDMTBankListRequest)
	drg.POST("/get/beneficiary", dmtHandler.GetDMTBeneficiariesRequest)
}
