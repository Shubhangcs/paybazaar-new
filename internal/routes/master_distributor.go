package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) MasterDistributorRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	mdRepo := repositories.NewMasterDistributorRepository(db, jwtUtils)
	mdHandler := handlers.NewMasterDistributorHandler(mdRepo)

	r.Router.POST("/md/login", mdHandler.LoginMasterDistributorRequest)
	mdrg := r.Router.Group("/md", middlewares.AuthorizationMiddleware(jwtUtils))
	mdrg.POST("/create", mdHandler.CreateMasterDistributorRequest, middlewares.RequireRoles("admin"))
	mdrg.PUT("/update", mdHandler.UpdateMasterDistributorRequest, middlewares.RequireRoles("admin", "master_distributor"))
	mdrg.DELETE("/delete/:master_distributor_id", mdHandler.DeleteMasterDistributorRequest, middlewares.RequireRoles("admin"))
	mdrg.GET("/get/all", mdHandler.ListMasterDistributorsRequest, middlewares.RequireRoles("admin"))
	mdrg.GET("/get/:master_distributor_id", mdHandler.GetMasterDistributorByIDRequest, middlewares.RequireRoles("master_distributor", "admin"))
	mdrg.GET("/get/:admin_id", mdHandler.ListMasterDistributorsByAdminIDRequest, middlewares.RequireRoles("admin"))
	mdrg.GET("/get/dropdown/:admin_id", mdHandler.GetMasterDistributorsByAdminIDForDropdownRequest, middlewares.RequireRoles("admin"))
}
