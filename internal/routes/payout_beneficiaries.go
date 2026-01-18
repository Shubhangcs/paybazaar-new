package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) RetailerBeneficiaryRoutes(
	db *database.Database,
	jwtUtils *pkg.JwtUtils,
) {

	retailerBeneficiaryRepo := repositories.NewRetailerBeneficiaryRepository(db)
	retailerBeneficiaryHandler := handlers.NewRetailerBeneficiaryHandler(
		retailerBeneficiaryRepo,
	)

	rb := r.Router.Group(
		"/retailer-beneficiary",
		middlewares.AuthorizationMiddleware(jwtUtils),
	)

	rb.POST("/create", retailerBeneficiaryHandler.CreateRetailerBeneficiary)
	rb.GET("/get/all", retailerBeneficiaryHandler.ListRetailerBeneficiaries)
	rb.GET("/get/:beneficiary_id", retailerBeneficiaryHandler.GetRetailerBeneficiaryByID)
	rb.PUT("/update/:beneficiary_id", retailerBeneficiaryHandler.UpdateRetailerBeneficiary)
	rb.PUT("/update/:beneficiary_id/verify", retailerBeneficiaryHandler.UpdateRetailerBeneficiaryVerification)
	rb.DELETE("/delete/:beneficiary_id", retailerBeneficiaryHandler.DeleteRetailerBeneficiary)
	rb.GET("/mobile/:mobile_number", retailerBeneficiaryHandler.GetBeneficiariesByMobileNumber)
}
