package routes

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/pkg"
)

type routes struct {
	Router *echo.Echo
}

type Config struct {
	ServerENV string
	JWTUtils  *pkg.JwtUtils
	Database  *database.Database
}

func NewRoutes(cfg Config) *routes {
	router := echo.New()
	log.Printf("server is running in %s mode", cfg.ServerENV)

	router.Validator = NewValidator()
	// Common Middlewares
	router.Use(middleware.CORS())

	routes := &routes{
		router,
	}

	// Routes Functions
	routes.AdminRoutes(cfg.Database, cfg.JWTUtils)
	routes.DistributorRoutes(cfg.Database, cfg.JWTUtils)
	routes.FundRequestRoutes(cfg.Database, cfg.JWTUtils)
	routes.MasterDistributorRoutes(cfg.Database, cfg.JWTUtils)
	routes.RetailerRoutes(cfg.Database, cfg.JWTUtils)
	routes.WalletTransactionRoutes(cfg.Database, cfg.JWTUtils)
	routes.RevertRoutes(cfg.Database, cfg.JWTUtils)
	routes.BankRouter(cfg.Database, cfg.JWTUtils)

	return routes
}
