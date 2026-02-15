package routes

import (
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/handlers"
	"github.com/levion-studio/paybazaar/internal/middlewares"
	"github.com/levion-studio/paybazaar/internal/repositories"
	"github.com/levion-studio/paybazaar/pkg"
)

func (r *routes) LimitRoutes(db *database.Database, jwtUtils *pkg.JwtUtils) {
	limitRepo := repositories.NewLimitRepository(db)
	limitHandler := handlers.NewLimitHandler(limitRepo)

	lrg := r.Router.Group("/limit", middlewares.AuthorizationMiddleware(jwtUtils))
	lrg.POST("/create", limitHandler.CreateLimitRequest, middlewares.RequireRoles("admin"))
	lrg.PUT("/update", limitHandler.UpdateLimitRequest, middlewares.RequireRoles("admin"))
	lrg.DELETE("/delete/:limit_id", limitHandler.DeleteLimitRequest, middlewares.RequireRoles("admin"))
	lrg.GET("/get/all", limitHandler.GetAllLimitsRequest, middlewares.RequireRoles("admin"))
	lrg.GET("/get/:retailer_id/:service", limitHandler.GetAllLimitsRequest, middlewares.RequireRoles("admin"))
}
