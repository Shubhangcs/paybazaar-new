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

	payoutBeneficiaryRepo := repositories.NewPayoutBeneficiaryRepository(db)
	payoutBeneficiaryHandler := handlers.NewPayoutBeneficiaryHandler(
		payoutBeneficiaryRepo,
	)

	rb := r.Router.Group(
		"/payout_beneficiary",
		middlewares.AuthorizationMiddleware(jwtUtils),
	)

	rb.POST("/create", payoutBeneficiaryHandler.CreatePayoutBeneficiary, middlewares.RequireRoles("retailer", "admin"))
	rb.GET("/get/all/:payout_id", payoutBeneficiaryHandler.ListPayoutBeneficiaries, middlewares.RequireRoles("retailer", "admin"))
	rb.GET("/get/:beneficiary_id", payoutBeneficiaryHandler.GetPayoutBeneficiaryByID, middlewares.RequireRoles("retailer", "admin"))
	rb.PUT("/update/:beneficiary_id", payoutBeneficiaryHandler.UpdatePayoutBeneficiary, middlewares.RequireRoles("retailer", "admin"))
	rb.PUT("/update/:beneficiary_id/verify", payoutBeneficiaryHandler.UpdatePayoutBeneficiaryVerification, middlewares.RequireRoles("retailer", "admin"))
	rb.DELETE("/delete/:beneficiary_id", payoutBeneficiaryHandler.DeletePayoutBeneficiary, middlewares.RequireRoles("retailer", "admin"))
	rb.GET("/mobile/:mobile_number", payoutBeneficiaryHandler.GetBeneficiariesByMobileNumber, middlewares.RequireRoles("retailer", "admin"))
}
