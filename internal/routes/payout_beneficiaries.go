package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) PayoutBeneficiaryRoutes(
	db *database.Database,
	jwtUtils *pkg.JwtUtils,
) {
	var benRepo = repositories.NewBeneficiaryRepo(db)
	var benHandler = handlers.NewBeneficiaryHandler(benRepo)

	rg := r.Router.Group("/bene", middlewares.AuthorizationMiddleware(jwtUtils))
	rg.GET("/get/beneficiaries/:phone", benHandler.GetBeneficiaries)
	rg.POST("/verify/beneficiaries", benHandler.VerifyBeneficiary)
	rg.POST("/add/beneficiary", benHandler.AddNewBeneficiary)
	rg.GET("/delete/beneficiary/:ben_id", benHandler.DeleteBeneficiary)
}
