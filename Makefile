# Load environment variables from .env if present
ifneq (,$(wildcard .env))
	include .env
	export
endif

MIGRATIONS_DIR=internal/database/migrations

.PHONY: migrate-up migrate-down migrate-force migrate-status migrate-create guard-db

# Ensure DATABASE_URL is set
guard-db:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif

# Apply all pending migrations
migrate-up: guard-db
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# Rollback last migration
migrate-down: guard-db
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

# Force migration version (use carefully)
# Example: make migrate-force v=4
migrate-force: guard-db
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(v)

# Show current migration version
migrate-status: guard-db
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

# Create a new migration file
# Example: make migrate-create name=add_indexes
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name)

run:
	@go run cmd/*
