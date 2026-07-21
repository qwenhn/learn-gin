include .env
export

include .env.dev
export

.PHONY: migrate-create migrate-up migrate-down migrate-down-n migrate-force migrate-drop migrate-goto
		docker-up docker-down 
		go-tidy
		sqlc-gen
		noapp stop-noapp dev stop-dev prod stop-prod logs-prod bash

ENV_FILE=.env
DEV_ENV_FILE=.env.dev
PROD_ENV_FILE=.env.prod
CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATION_DIRS = internal/db/migrations

NOAPP_COMPOSE=docker/docker-compose.noapp.yml
DEV_COMPOSE=docker/docker-compose.dev.yml
PROD_COMPOSE=docker/docker-compose.prod.yml

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
	docker compose -f $(DEV_COMPOSE) down

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

# Build a binary file
build:
	go build -o bin/shopping-api ./cmd/api

run-binary:
	./bin/shopping-api

# https://docs.sqlc.dev/en/stable/overview/install.html
sqlc-gen:
	sqlc generate

# Build container for noapp
noapp:
	docker compose -f $(NOAPP_COMPOSE) down
	docker compose -f $(NOAPP_COMPOSE) --env-file $(ENV_FILE) up -d --build

# Stop all container noapp
stop-noapp:
	docker compose -f $(NOAPP_COMPOSE) down

# Build container for dev
dev:
	docker compose -f $(DEV_COMPOSE) down
	docker compose -f $(DEV_COMPOSE) --env-file $(DEV_ENV_FILE) up --build

# Stop all container dev
stop-dev:
	docker compose -f $(DEV_COMPOSE) down

# Build container for production
prod:
	docker compose -f $(PROD_COMPOSE) down
	docker compose -f $(PROD_COMPOSE) --env-file $(ENV_FILE) up -d --build

# Stop all container production
stop-prod:
	docker compose -f $(PROD_COMPOSE) down

# View logs container production
logs-prod:
	docker compose -f $(PROD_COMPOSE) logs -f --tail=100

# Go to container golang-api
bash:
	docker exec -it golang-api /bin/sh
