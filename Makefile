.PHONY: all build test lint proto clean dev-up dev-down race chaos diagnostics deps proto-lint proto-breaking help

# Configuration
COMPOSE_FILE := deploy/compose/docker-compose.dev.yaml
BIN_DIR      := bin
SERVICES     := collector correlator topology rca remediation

all: proto lint test build

help:
	@echo "CortexOps Development Makefile"
	@echo "------------------------------"
	@echo "dev-up         : Start ALL services (infrastructure + runtime)"
	@echo "dev-down       : Stop all services"
	@echo "proto          : Generate Go code from Protobuf"
	@echo "lint           : Run golangci-lint"
	@echo "test           : Run all tests"
	@echo "build          : Build all service binaries"
	@echo "clean          : Remove build artifacts"
	@echo "bootstrap      : Deploy demo workloads to Kind/Minikube"
	@echo "demo-failure   : Inject deterministic failure (SCENARIO=rollout-fail|crashloop|scaling)"
	@echo "demo-recovery  : Restore demo environment from failure"
	@echo "diagnostics    : Run live platform diagnostics"
	@echo "validate-pipeline: Run golden path operational validation"
	@echo "chaos-test     : Run operational failure validation"
	@echo "dashboards     : Show dashboard access information"

# Initialize local dev environment
dev-up devup:
	@echo "Starting CortexOps platform..."
	docker compose -f $(COMPOSE_FILE) --profile full up -d --build

dev-down devdown:
	@echo "Stopping CortexOps platform..."
	docker compose -f $(COMPOSE_FILE) --profile full down -v

# Verification
verify-runtime:
	@echo "Verifying runtime health..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder bash scripts/verify-runtime.sh

validate-pipeline:
	@echo "Running Golden Path Validation..."
	bash scripts/validate-pipeline.sh

validate-temporal:
	@echo "Validating Temporal Workflows..."
	bash scripts/validate-temporal.sh

chaos-test:
	@echo "Running chaos validation (Duplicate Storm & Idempotency)..."
	go run ./cmd/chaos/main.go duplicate-storm
	go run ./cmd/chaos/main.go workflow-idempotency

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
	docker compose -f $(COMPOSE_FILE) run --rm golangci-lint golangci-lint run -v --timeout 10m

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
	@echo "Running CortexOps diagnostics..."
	@echo "--- Platform Health ---"
	curl -s http://localhost:9091/debug/healthz || echo "Platform offline"
	@echo "\n--- Active Incidents ---"
	curl -s http://localhost:9091/debug/incidents/active || echo "Platform offline"
	@echo "\n--- Topology Blast Radius (Example) ---"
	curl -s "http://localhost:9091/debug/graph/blast-radius?id=pod/cortexops-demo/demo-api" || echo "Platform offline"

# Demo Automation
SCENARIO ?= rollout-fail

bootstrap:
	@echo "Bootstrapping demo environment..."
	go run ./cmd/demo/main.go bootstrap

demo-failure:
	@echo "Injecting demo failure ($(SCENARIO))..."
	go run ./cmd/demo/main.go inject -scenario=$(SCENARIO)

demo-recovery:
	@echo "Recovering demo environment..."
	go run ./cmd/demo/main.go recover

dashboards:
	@echo "CortexOps Dashboards"
	@echo "--------------------"
	@echo "Grafana URL: http://localhost:3000"
	@echo "Dashboards source: deploy/grafana/dashboards/"
	@echo "Available: cortexops-health.json, cortexops-demo.json"

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
