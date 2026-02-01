package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) DMTRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	dmtRepo := repositories.NewDMTRepository(db)
	dmtHandlers := handlers.NewDMTHandler(dmtRepo)

	dmtrg := r.Router.Group("/dmt", middlewares.AuthorizationMiddleware(jwtUtils))
	dmtrg.POST("/check/exist", dmtHandlers.CheckDMTWalletExistsRequest)
	dmtrg.POST("/create/wallet", dmtHandlers.CreateDMTWalletRequest)
	dmtrg.POST("/verify/create/wallet", dmtHandlers.VerifyDMTWalletCreationRequest)
	dmtrg.POST("/create/beneficiary", dmtHandlers.CreateDMTBeneficiaryRequest)
	dmtrg.POST("/get/beneficiaries", dmtHandlers.GetDMTBeneficiariesRequest)
	dmtrg.POST("/delete/beneficieries", dmtHandlers.DeleteDMTBeneficiaryRequest)
	dmtrg.POST("/verify/delete/beneficieries", dmtHandlers.VerifyDMTBeneficiaryDeleteRequest)
}
