# Pipeline Reference

This document provides a technical reference for each stage in the CortexOps CI/CD pipeline.

## 1. Frontend Validation (`frontend.yml`)

Runs exclusively in the `frontend/` directory.

- **Dependencies**: Uses `npm ci` with `~/.npm` caching based on `frontend/package-lock.json`.
- **Linting**: `npm run lint` leverages ESLint configured for Next.js and TypeScript.
- **Type Checking**: `npx tsc --noEmit` ensures strict type safety without producing build artifacts.
- **Build**: `npm run build` compiles the Next.js production bundle.

## 2. Backend Validation (`backend.yml`)

Runs from the repository root.

- **Dependencies**: `go mod tidy` and verifies that `go.mod` and `go.sum` are synchronized.
- **Static Analysis**: `go vet ./...` and `golangci-lint` (v1.55.2).
- **Unit Tests**: `go test -v -race ./...` with `CGO_ENABLED=1` for race detection.
- **Compilation**: Builds `collector`, `correlator`, `topology`, `rca`, and `remediation` binaries into `bin/`.

## 3. Docker Validation (`docker.yml`)

Ensures all services are container-ready.

- **Matrix Execution**: Evaluates `collector`, `correlator`, `topology`, `rca`, `remediation`, and `frontend` concurrently.
- **Caching**: Utilizes `docker/setup-buildx-action` and `actions/cache` to store BuildKit layers in `/tmp/.buildx-cache`.
- **Validation**: Performs `docker build` without pushing to a registry (`push: false`).

## 4. Helm Validation (`helm.yml`)

Ensures Infrastructure as Code compliance.

- **Linting**: `helm lint deploy/helm/cortexops` checks chart formatting and structure.
- **Templating**: `helm template` tests rendering of all sub-charts with default `values.yaml`.

## 5. Security Scanning (`security.yml`)

See [SECURITY_PIPELINE.md](./SECURITY_PIPELINE.md) for detailed configuration.

## 6. Release Automation (`release.yml`)

See [RELEASE_PROCESS.md](./RELEASE_PROCESS.md) for detailed release mechanics.
