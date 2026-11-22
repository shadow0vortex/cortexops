.PHONY: all build test lint proto clean dev-up dev-down dev-reset race chaos chaos-test diagnostics deps tidy proto-lint proto-breaking verify-runtime validate-pipeline validate-temporal help

# Configuration
COMPOSE_FILE := deploy/compose/docker-compose.dev.yaml
BIN_DIR      := bin
SERVICES     := collector correlator topology rca remediation

ifeq ($(OS),Windows_NT)
BASH ?= "C:/Program Files/Git/bin/bash.exe"
else
BASH ?= bash
endif

all: proto lint test build

help:
	@echo "CortexOps Development Makefile"
	@echo "------------------------------"
	@echo "dev-up         : Start ALL services (infrastructure + runtime)"
	@echo "dev-down       : Stop all services"
	@echo "dev-reset      : Stop services and delete local volumes"
	@echo "proto          : Generate Go code from Protobuf"
	@echo "lint           : Run golangci-lint"
	@echo "test           : Run all tests"
	@echo "build          : Build all service binaries"
	@echo "tidy           : Update Go module files"
	@echo "clean          : Remove build artifacts"
	@echo "bootstrap      : Deploy demo workloads to Kind/Minikube"
	@echo "demo-failure   : Inject deterministic failure (SCENARIO=rollout-fail|crashloop|scaling)"
	@echo "demo-recovery  : Restore demo environment from failure"
	@echo "diagnostics    : Run live platform diagnostics"
	@echo "validate-pipeline: Run golden path operational validation"
	@echo "chaos-test     : Run integration-tagged chaos tests (requires a reachable Kubernetes cluster)"
	@echo "dashboards     : Show dashboard access information"

# Initialize local dev environment
dev-up devup:
	@echo "Starting CortexOps platform..."
	docker compose -f $(COMPOSE_FILE) --profile full up -d --build

dev-down devdown:
	@echo "Stopping CortexOps platform..."
	docker compose -f $(COMPOSE_FILE) --profile full down

dev-reset:
	@echo "Stopping CortexOps platform and deleting local volumes..."
	docker compose -f $(COMPOSE_FILE) --profile full down -v

# Verification
verify-runtime:
	@echo "Verifying runtime health..."
	docker compose -f $(COMPOSE_FILE) --profile full --profile tools run --rm runtime-check

validate-pipeline:
	@echo "Running Golden Path Validation..."
	$(BASH) scripts/validate-pipeline.sh

validate-temporal:
	@echo "Validating Temporal Workflows..."
	$(BASH) scripts/validate-temporal.sh

chaos-test:
	@echo "Running integration-tagged chaos tests..."
	go test -v -tags=integration ./test/e2e/suites -run '^TestChaos_'

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
	@echo "Downloading dependencies..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go mod download

tidy:
	@echo "Tidying Go module files..."
	docker compose -f $(COMPOSE_FILE) run --rm go-builder go mod tidy

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

# Run cluster-dependent chaos tests
chaos: chaos-test

# Run diagnostics (Dockerized)
diagnostics:
	@echo "Running CortexOps diagnostics..."
	@echo "--- Platform Health ---"
	@curl -fsS http://localhost:9091/debug/healthz
	@echo --- Topology Nodes ---
	@curl -fsS http://localhost:9091/v1/topology/nodes
	@echo --- Topology Blast Radius (Example) ---
	@curl -fsS "http://localhost:9091/v1/topology/blast-radius/pod%2Fcortexops-demo%2Fdemo-api"

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
else
    RM := rm -rf
endif

# Clean build artifacts
clean:
	@echo "Cleaning artifacts..."
	$(RM) $(BIN_DIR)
