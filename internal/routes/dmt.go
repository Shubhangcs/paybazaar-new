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
	drg.POST("/check/wallet", dmtHandler.CheckDMTWalletExistsRequest, middlewares.RequireRoles("retailer"))
	drg.POST("/create/wallet", dmtHandler.CreateDMTWalletRequest, middlewares.RequireRoles("retailer"))
	drg.POST("/verify/wallet", dmtHandler.VerifyDMTWalletRequest, middlewares.RequireRoles("retailer"))
	drg.POST("/add/beneficiary", dmtHandler.AddDMTBeneficiaryRequest, middlewares.RequireRoles("retailer"))
	drg.GET("/get/banks", dmtHandler.GetDMTBankListRequest, middlewares.RequireRoles("retailer"))
	drg.POST("/get/beneficiary", dmtHandler.GetDMTBeneficiariesRequest, middlewares.RequireRoles("retailer"))
}
