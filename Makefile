.PHONY: help dev dev-up dev-down dev-logs db-shell api-shell build test clean

# Default target
help:
	@echo "PetTime Development Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  dev        Start development environment (docker compose up)"
	@echo "  dev-up     Start development environment in background"
	@echo "  dev-down   Stop development environment"
	@echo "  dev-logs   Show logs from all services"
	@echo "  db-shell   Open PostgreSQL shell"
	@echo "  api-shell  Open shell in API container"
	@echo "  build      Build the Go binary locally"
	@echo "  test       Run tests"
	@echo "  clean      Clean up build artifacts"
	@echo ""

# Development
dev:
	docker compose up --build

dev-up:
	docker compose up -d --build

dev-down:
	docker compose down

dev-logs:
	docker compose logs -f

# Database
db-shell:
	docker compose exec db psql -U pettime -d pettime

# API
api-shell:
	docker compose exec api sh

# Build
build:
	cd backend && go build -o pettime-api ./cmd/api

# Test
test:
	cd backend && go test -v ./...

# Clean
clean:
	rm -f backend/pettime-api
	docker compose down -v
