.PHONY: all build test lint proto clean dev-up dev-down race chaos diagnostics deps proto-lint proto-breaking

# Docker Compose File
COMPOSE_FILE := deploy/compose/docker-compose.dev.yaml

all: proto lint test build

# Initialize local dev environment
dev-up devup:
	@echo "Starting local dev dependencies..."
	docker compose -f $(COMPOSE_FILE) up -d

dev-down devdown:
	@echo "Stopping local dev dependencies..."
	docker compose -f $(COMPOSE_FILE) down -v

# Protobuf generation (Dockerized)
proto:
	@echo "Generating Go files from Protobuf definitions..."
	docker compose -f $(COMPOSE_FILE) run --rm -w /workspace/api/v1 proto-builder generate

proto-lint:
	@echo "Linting Protobuf definitions..."
	docker compose -f $(COMPOSE_FILE) run --rm -w /workspace/api/v1 proto-builder lint

proto-breaking:
	@echo "Checking for breaking Protobuf changes..."
	docker compose -f $(COMPOSE_FILE) run --rm -w /workspace/api/v1 proto-builder breaking --against '.git#branch=main'

# Fetch dependencies (Dockerized)
deps:
	@echo "Tidying and downloading dependencies..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go mod tidy
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go mod download

# Lint Go code (Dockerized)
lint:
	@echo "Linting Go code..."
	docker compose -f $(COMPOSE_FILE) run --rm golangci-lint golangci-lint run -v

# Build all binaries (Dockerized)
build: deps
	@echo "Building services..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go build -o bin/collector ./cmd/placeholder
	# Add actual builds here as services are created

# Run tests (Dockerized)
test: deps
	@echo "Running tests..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go test -v ./...

# Run race tests (Dockerized)
race: deps
	@echo "Running tests with race detector..."
	docker compose -f $(COMPOSE_FILE) run --rm -e CGO_ENABLED=1 go-builder go test -v -race ./...

# Run chaos tests (Dockerized)
chaos: deps
	@echo "Running chaos tests..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go test -v -tags=chaos ./test/e2e/chaos/...

# Run diagnostics (Dockerized)
diagnostics:
	@echo "Running diagnostics..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go run ./cmd/placeholder/main.go -diagnostics || true

# Clean build artifacts
clean:
	@echo "Cleaning artifacts..."
	rm -rf bin/
	rm -f api/v1/*.pb.go
