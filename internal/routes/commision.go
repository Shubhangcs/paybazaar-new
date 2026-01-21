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
	crg.POST("/create", commisionHandler.CreateCommisionRequest, middlewares.RequireRoles("admin"))
	crg.DELETE("/delete/commision", commisionHandler.DeleteCommisionRequest, middlewares.RequireRoles("admin"))
	crg.PUT("/update/commision", commisionHandler.UpdateCommisionDetailsRequest, middlewares.RequireRoles("admin"))
	crg.GET("/get/commision/:user_id/:service", commisionHandler.GetCommisionByUserIDAndServiceRequest, middlewares.RequireRoles("admin"))
	crg.GET("/get/commision/:user_id", commisionHandler.GetCommisionByUserIDAndServiceRequest, middlewares.RequireRoles("admin"))
	crg.GET("/get/commisions/:commision_id", commisionHandler.GetCommisionDetailsByCommisionIDRequest, middlewares.RequireRoles("admin"))
	crg.GET("/get/tds/:user_id", commisionHandler.GetTDSCommisionByUserIDRequest, middlewares.RequireRoles("admin", "retailer"))
	crg.GET("/get/tds", commisionHandler.GetAllTDSCommisionRequest, middlewares.RequireRoles("admin"))
}
