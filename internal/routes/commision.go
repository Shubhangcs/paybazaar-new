package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) CommisionRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	commisionRepo := repositories.NewCommisionRepository(db)
	commisionHandler := handlers.NewCommisionHandler(commisionRepo)

	crg := r.Router.Group("/commision", middlewares.AuthorizationMiddleware(jwtUtils))
	crg.POST("/create", commisionHandler.CreateCommision, middlewares.RequireRoles("admin"))
	crg.DELETE("/delete/:commision_id", commisionHandler.DeleteCommision, middlewares.RequireRoles("admin"))
	crg.PUT("/update/:commision_id", commisionHandler.UpdateCommision, middlewares.RequireRoles("admin"))
	crg.GET("/get/all", commisionHandler.GetAllCommisions, middlewares.RequireRoles("admin"))
	crg.GET("/get/user/:user_id", commisionHandler.GetCommisionByUserID, middlewares.RequireRoles("admin"))
	crg.GET("/get/:commision_id", commisionHandler.GetCommisionByID, middlewares.RequireRoles("admin"))
}
