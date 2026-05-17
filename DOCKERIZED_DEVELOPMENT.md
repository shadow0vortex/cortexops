# CortexOps Dockerized Development Workflow

This document outlines the fully containerized and reproducible development workflow for CortexOps.

## Philosophy

To eliminate platform-specific setup issues, Go toolchain drift, and Makefile incompatibilities, the entire CortexOps developer workflow runs through Docker.

**You do NOT need to install locally:**
- `protoc` (Protobuf Compiler)
- `protoc-gen-go` / `protoc-gen-go-grpc`
- Go linters (`golangci-lint`)
- `staticcheck`

**Required Local Dependencies:**
- Docker Desktop
- Docker Compose
- Git

---

## Local Sandbox Startup

A complete developer platform is provided out-of-the-box via `docker-compose.dev.yaml`.

To initialize the environment (NATS, PostgreSQL, Qdrant, Temporal, Prometheus, Grafana), run:

```bash
make dev-up
```

To tear down the local dependencies:

```bash
make dev-down
```

---

## Make Targets (Containerized Tooling)

All Makefile targets have been refactored to use Docker compose profiles and images to ensure execution is deterministic across macOS, Linux, and Windows.

### 1. Protobuf Toolchain

We use a pinned, hermetic protobuf builder container.

```bash
make proto         # Generates Go files from Protobuf definitions
make proto-lint    # Lints Protobuf definitions
make proto-breaking # Checks for breaking Protobuf changes
```

### 2. Dependency Management

Downloads and tidies Go modules in a reproducible environment:

```bash
make deps
```

### 3. Linting and Testing

All validation happens in standard container environments:

```bash
make lint          # Lints Go code via golangci-lint
make test          # Runs unit tests
make race          # Runs tests with the race detector enabled
make chaos         # Runs e2e chaos tests
make diagnostics   # Runs diagnostics tooling
```

### 4. Building Binaries

Build your binaries without needing Go installed locally:

```bash
make build         # Builds services into bin/
```

### 5. Cleaning Artifacts

```bash
make clean         # Cleans compiled binaries and generated protobuf files
```

---

## CI/CD Reproducibility

The GitHub Actions workflows (`.github/workflows/ci.yaml`) have been updated to reuse the Dockerized tooling. This means CI acts as a perfect mirror of your local development environment:
- Identical `protoc` and plugin versions.
- Identical Go compiler and linter versions.
- Replay-safe builds that generate the exact same protobufs regardless of the runtime environment.

## DevContainer Onboarding

If you prefer VS Code, the provided `.devcontainer/devcontainer.json` supports full integration. Opening the project in the DevContainer automatically triggers the `make deps` and `make proto` bootstrapping so you can immediately begin coding with working Intellisense.
