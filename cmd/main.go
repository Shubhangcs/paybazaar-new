package main

import (
	"log"

	"github.com/levion-studio/paybazaar/internal/config"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/routes"
	"github.com/levion-studio/paybazaar/pkg"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	cfg := config.Load()

	db, err := database.NewDatabaseConnection(database.Config{
		DatabaseURL: cfg.DatabaseURL,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	jwtUtils := pkg.NewJwtUtils(pkg.JwtConfig{
		SecretKey: cfg.SecretKey,
		Expiry:    cfg.Expiry,
	})

	router := routes.NewRoutes(routes.Config{
		ServerENV: cfg.ServerEnv,
		JWTUtils:  jwtUtils,
		Database:  db,
	})

	return router.Router.Start(cfg.ServerPort)
}
