# CI/CD Architecture

This document describes the continuous integration and delivery architecture for the CortexOps platform.

## Pipeline Architecture

CortexOps utilizes a modular CI/CD architecture based on GitHub Actions. Instead of a monolithic workflow, responsibilities are split across domain-specific pipelines:

1. **Frontend Pipeline** (`frontend.yml`): Focuses exclusively on the Next.js frontend monorepo.
2. **Backend Pipeline** (`backend.yml`): Focuses on Go microservices, Protobufs, and API layers.
3. **Security Pipeline** (`security.yml`): Runs SAST, filesystem scanners, and credential leak detection.
4. **Docker Validation** (`docker.yml`): Validates container builds across all components.
5. **Helm Validation** (`helm.yml`): Ensures infrastructure-as-code manifests are compliant.
6. **Release Automation** (`release.yml`): Handles versioning, binary cross-compilation, and release artifacts.

## Branch Strategy

- **`main`**: The primary branch. It is protected and requires passing status checks from all modular workflows before a PR can be merged.
- **Feature Branches**: Development occurs on short-lived feature branches, submitting PRs to `main`.

## Trigger Matrix

| Workflow | Path Triggers | Events |
|---|---|---|
| **Frontend** | `frontend/**` | Push, Pull Request |
| **Backend** | `**/*.go`, `go.mod`, `Makefile` | Push, Pull Request |
| **Security** | All files | Push, Pull Request, Cron (Weekly) |
| **Docker** | All files | Push, Pull Request |
| **Helm** | `deploy/helm/**` | Push, Pull Request |
| **Release** | Tags matching `v*` | Push (Tags) |

## Dependency Management & Caching

To optimize CI execution times, aggressive caching strategies are employed:
- **Go Modules**: Cached via `actions/setup-go`.
- **Node Modules**: Cached via `actions/setup-node`.
- **Docker Layers**: Cached via `docker/setup-buildx-action` (BuildKit local caching).

## Quality Gates

Before merging to `main`, a Pull Request must satisfy:
1. `npm run lint` and `npx tsc --noEmit` (Frontend).
2. `golangci-lint` and `go vet` (Backend).
3. `go test -v -race` (Unit tests).
4. Trivy and Gosec scans (no new critical/high vulnerabilities).
5. Successful Docker builds for all updated services.
6. `helm lint` and `helm template` (if Helm charts are modified).
