package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) FundTransferRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	fundTransferRepo := repositories.NewFundTransferRepository(db)
	fundTransferHandler := handlers.NewFundTransferHandler(fundTransferRepo)

	ftr := r.Router.Group("/fund_transfer", middlewares.AuthorizationMiddleware(jwtUtils))
	ftr.POST("/create", fundTransferHandler.CreateFundTransfer)
	ftr.GET("/from", fundTransferHandler.GetFundTransfersByFromID)
	ftr.GET("/to", fundTransferHandler.GetFundTransfersByToID)
}
