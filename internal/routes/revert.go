package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) RevertRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	revertRepo := repositories.NewRevertRepository(db)
	revertHandler := handlers.NewRevertHandler(revertRepo)

	rrg := r.Router.Group("/revert", middlewares.AuthorizationMiddleware(jwtUtils))
	rrg.POST("/create", revertHandler.CreateRevertRequest)
}
