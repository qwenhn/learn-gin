include .env
export

.PHONY: server 
		migrate-create migrate-up migrate-down migrate-down-n migrate-force migrate-drop migrate-goto
		docker-up docker-down 
		go-tidy

CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATION_DIRS = internal/db/migrations

# Run server
server:
	go run .

# CLI documentation: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
# Create a new migration (make migrate-create NAME=)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)

# Run all pending migration (make migrate-up)
migrate-up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up

# Rollback the last migration
migrate-down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1

# Start development environment
docker-up:
	docker compose up --build

# Stop development environment
docker-down:
	docker compose down

# Rollback N migrations
migrate-down-n:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(N)

# Force migration version (use with caution example: make migrate-force VERSION=1) 
migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)

# Drop everything (include schema migration)
migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop

# Apply specific migration version (make migrate-goto VERSION=1)
migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)

# Update Go dependencies
go-tidy:
	find . -name go.mod -execdir go mod tidy \;
