.PHONY: all build test lint proto clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get

# Paths
API_DIR=api/v1
PKG_DIR=pkg

all: proto build test

# Protobuf generation (Requires protoc and protoc-gen-go installed)
proto:
	@echo "Generating Go files from Protobuf definitions..."
	protoc --proto_path=. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(API_DIR)/*.proto

# Fetch dependencies
deps:
	@echo "Tidying and downloading dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Build all binaries
build: deps
	@echo "Building services..."
	# $(GOBUILD) -o bin/collector ./cmd/collector
	# Add actual builds here as services are created

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# Clean build artifacts
clean:
	@echo "Cleaning artifacts..."
	rm -rf bin/
	rm -f $(API_DIR)/*.pb.go

# Initialize local dev environment
dev-up:
	@echo "Starting local dev dependencies (NATS, Postgres, Qdrant)..."
	docker-compose -f deploy/compose/docker-compose.yaml up -d

dev-down:
	@echo "Stopping local dev dependencies..."
	docker-compose -f deploy/compose/docker-compose.yaml down
