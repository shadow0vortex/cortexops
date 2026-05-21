.PHONY: all build test lint proto clean dev-up dev-down race chaos diagnostics deps proto-lint proto-breaking help

# Configuration
COMPOSE_FILE := deploy/compose/docker-compose.dev.yaml
BIN_DIR      := bin
SERVICES     := collector

all: proto lint test build

help:
	@echo "CortexOps Development Makefile"
	@echo "------------------------------"
	@echo "dev-up         : Start infrastructure dependencies"
	@echo "dev-down       : Stop infrastructure dependencies"
	@echo "proto          : Generate Go code from Protobuf"
	@echo "lint           : Run golangci-lint"
	@echo "test           : Run all tests"
	@echo "build          : Build all service binaries"
	@echo "clean          : Remove build artifacts"

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
	docker compose -f $(COMPOSE_FILE) run --rm golangci-lint golangci-lint run -v --timeout 5m

# Build all binaries (Dockerized)
build: deps $(SERVICES)

$(SERVICES):
	@echo "Building $@..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder sh -c "mkdir -p $(BIN_DIR) && go build -o $(BIN_DIR)/$@ ./cmd/$@"

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
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go run ./cmd/collector/main.go -diagnostics || true

# Detect OS
ifeq ($(OS),Windows_NT)
    RM := powershell.exe -NoProfile -Command Remove-Item -Recurse -Force
    RM_F := powershell.exe -NoProfile -Command "Get-ChildItem -Path api/v1/*.pb.go | Remove-Item -Force"
else
    RM := rm -rf
    RM_F := rm -f api/v1/*.pb.go
endif

# Clean build artifacts
clean:
	@echo "Cleaning artifacts..."
	$(RM) $(BIN_DIR)
	$(RM_F)
