package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) DistributorRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	disRepo := repositories.NewDistributorRepository(db, jwtUtils)
	disHandler := handlers.NewDistributorHandler(disRepo)

	r.Router.POST("/distributor/login", disHandler.LoginDistributorRequest)
	drg := r.Router.Group("/distributor", middlewares.AuthorizationMiddleware(jwtUtils))
	drg.POST("/create", disHandler.CreateDistributorRequest, middlewares.RequireRoles("admin", "master_distributor"))
	drg.PUT("/update", disHandler.UpdateDistributorRequest, middlewares.RequireRoles("admin", "distributor"))
	drg.DELETE("/delete/:distributor_id", disHandler.DeleteDistributorRequest, middlewares.RequireRoles("admin"))
	drg.GET("/get/all", disHandler.ListDistributorsRequest, middlewares.RequireRoles("admin"))
	drg.GET("/get/:distributor_id", disHandler.GetDistributorByIDRequest, middlewares.RequireRoles("admin", "master_distributor", "distributor"))
	drg.GET("/get/:master_distributor_id", disHandler.ListDistributorsByMasterDistributorIDRequest, middlewares.RequireRoles("admin", "master_distributor"))
	drg.GET("/get/dropdown/:master_distributor_id", disHandler.GetDistributorsByMasterDistributorIDForDropdownRequest, middlewares.RequireRoles("admin", "master_distributor"))
}
