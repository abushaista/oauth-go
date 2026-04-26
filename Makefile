.PHONY: help build run test clean docker-up docker-down docker-build dev

help:
	@echo "OAuth Server - Available commands"
	@echo ""
	@echo "  make build          Build the OAuth server binary"
	@echo "  make run            Run the OAuth server"
	@echo "  make test           Run tests"
	@echo "  make fmt            Format code"
	@echo "  make lint           Run linter"
	@echo "  make clean          Clean build artifacts"
	@echo ""
	@echo "  make docker-up      Start all containers (Postgres + OAuth Server)"
	@echo "  make docker-down    Stop all containers"
	@echo "  make docker-build   Rebuild the OAuth server Docker image"
	@echo "  make docker-db      Start only Postgres + Adminer"
	@echo ""
	@echo "  make dev            Start development server (local Go + Dockerized DB)"

build:
	go build -o bin/oauth-server cmd/api/main.go

run: build
	./bin/oauth-server

test:
	go test -v -cover ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

clean:
	rm -f bin/oauth-server
	rm -rf bin/
	rm -f *.out

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose up -d --build

docker-db:
	docker compose up -d postgres adminer

dev: docker-db
	@echo "Waiting for Postgres to be ready..."
	@sleep 3
	DATABASE_URL="postgres://oauth_user:oauth_password@localhost:5432/oauth_db?sslmode=disable" \
	JWT_SECRET="dev-secret-key" \
	go run cmd/api/main.go

.DEFAULT_GOAL := help
