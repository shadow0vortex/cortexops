# CortexOps Demo Guide

This guide walkthrough the steps to observe CortexOps in action using the provided demonstration environment.

## Prerequisites
- Docker & Docker Compose
- A local Kubernetes cluster (Kind or Minikube)
- Go 1.22+

## 1. Environment Setup

Start the infrastructure dependencies (NATS, Temporal, Qdrant, Postgres):
```bash
make dev-up
```

Verify that all services are healthy:
```bash
docker compose -f deploy/compose/docker-compose.dev.yaml ps
```

## 2. Bootstrapping the Demo

Deploy the synthetic microservice topology into your cluster:
```bash
make bootstrap
```
This creates the `cortexops-demo` namespace and deploys a multi-tier application (Frontend -> API -> Worker/DB/Cache).

## 3. Running CortexOps

In a separate terminal, start the CortexOps Collector:
```bash
$env:NATS_URL="nats://localhost:4222"; go run ./cmd/collector/main.go
```

## 4. Triggering a Failure

Inject a deterministic failure (e.g., a rollout failure in the frontend):
```bash
make demo-failure SCENARIO=rollout-fail
```

## 5. Observing the Platform

### Visualization
Open Grafana to see the telemetry ingestion and incident creation:
- **URL**: `http://localhost:3000`
- **Dashboard**: `CortexOps Demo Overview`

### Remediation
Open the Temporal UI to watch the remediation workflow:
- **URL**: `http://localhost:8233`
- **Workflow**: `RemediationWorkflow`

### Diagnostics
Inspect the internal state via the Diagnostics API:
```bash
make diagnostics
```

## 6. Recovery

Revert the failure and restore the environment:
```bash
make demo-recovery
```
