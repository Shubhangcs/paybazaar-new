package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DatabaseURL string
}

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabaseConnection(cfg Config) (*Database, error) {
	// Creating connection config
	conf, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to configure database connection: %w", err)
	}

	// Configuring connection
	conf.MaxConns = 10
	conf.MinConns = 2
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute
	conf.ConnConfig.ConnectTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Creating a pool connection
	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Returning the connection
	return &Database{
		pool: pool,
	}, nil
}

func (db *Database) Ping(ctx context.Context) error {
	// Pinging the database for availablity
	return db.pool.Ping(ctx)
}

func (db *Database) Close() {
	// Closing database connection
	db.pool.Close()
}
