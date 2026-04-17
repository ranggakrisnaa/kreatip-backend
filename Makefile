APP_NAME    = kreatip-backend
BUILD_DIR   = bin
MAIN_WEB    = cmd/web/main.go
MAIN_WORKER = cmd/worker/main.go
MIGRATE_URL = postgres://postgres:postgres@localhost:5432/kreatip?sslmode=disable
MIGRATIONS  = db/migrations

.PHONY: all run run-worker build build-worker deps install-tools lint test migrate-up migrate-down migrate-create clean docker-up docker-down

all: build

## Run HTTP server
run:
	go run $(MAIN_WEB)

## Run async worker
run-worker:
	go run $(MAIN_WORKER)

## Build HTTP server binary
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_WEB)

## Build worker binary
build-worker:
	go build -o $(BUILD_DIR)/$(APP_NAME)-worker $(MAIN_WORKER)

## Download and tidy dependencies
deps:
	go mod download
	go mod tidy

## Install dev CLI tools (golang-migrate, golangci-lint)
install-tools:
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

## Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

## Run all tests
test:
	go test ./... -v -race -cover

## Run tests with coverage report
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

## Apply all pending migrations
migrate-up:
	migrate -path $(MIGRATIONS) -database "$(MIGRATE_URL)" up

## Rollback last migration
migrate-down:
	migrate -path $(MIGRATIONS) -database "$(MIGRATE_URL)" down 1

## Create new migration: make migrate-create name=create_users
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS) -seq $(name)

## Start Docker services (postgres, redis, mailhog)
docker-up:
	docker compose up -d

## Stop Docker services
docker-down:
	docker compose down

## Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)/ coverage.out coverage.html
