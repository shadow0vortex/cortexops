# CortexOps

**Kubernetes telemetry in. Causal incidents out. Remediation under governance.**

<div align="center">

[![Backend CI](https://github.com/shadow0vortex/cortexops/actions/workflows/backend.yml/badge.svg)](https://github.com/shadow0vortex/cortexops/actions/workflows/backend.yml)
[![Frontend CI](https://github.com/shadow0vortex/cortexops/actions/workflows/frontend.yml/badge.svg)](https://github.com/shadow0vortex/cortexops/actions/workflows/frontend.yml)
[![Docker Images](https://github.com/shadow0vortex/cortexops/actions/workflows/docker.yml/badge.svg)](https://github.com/shadow0vortex/cortexops/actions/workflows/docker.yml)
[![Helm Validation](https://github.com/shadow0vortex/cortexops/actions/workflows/helm.yml/badge.svg)](https://github.com/shadow0vortex/cortexops/actions/workflows/helm.yml)
[![Security Pipeline](https://github.com/shadow0vortex/cortexops/actions/workflows/security.yml/badge.svg)](https://github.com/shadow0vortex/cortexops/actions/workflows/security.yml)

</div>

CortexOps is an event-driven Kubernetes operations control plane. It turns real cluster events into correlated incidents, advisory RCA, and durable, policy-gated remediation workflows.

> [!IMPORTANT]
> CortexOps is an active-development engineering platform, not a hosted production SaaS. The included local stack is verified end to end. Hosted user identity, tenant isolation, durable customer records, and live authenticated dashboards are still roadmap work.

## What happens when Kubernetes has a bad day?

```text
Kubernetes event
       |
       v
  [ Collector ]  --> typed telemetry
       |
       v
 [ NATS JetStream ] --> durable, replayable transport
       |
       v
 [ Correlator ]  --> incident + blast radius
       |
       v
 [ RCA Engine ]  --> evidence and advisory analysis
       |
       v
 [ Remediation ] --> OPA policy gate
                            |
                            v
                      [ Temporal ]
                            |
                            v
                    verified Kubernetes action
```

## Built for the control plane

| Capability | Implementation |
|---|---|
| Kubernetes telemetry | Real event and metric collection |
| Typed transport | Protobuf envelopes over NATS JetStream |
| Incident correlation | Time, trace, and topology-aware heuristics |
| Dependency intelligence | In-memory graph with PostgreSQL-backed persistence |
| RCA | RAG context, LLM advisory analysis, and degraded-mode fallback |
| Remediation | Temporal workflows with deterministic verification and rollback |
| Policy enforcement | Open Policy Agent, fail-closed by default |
| Observability | Prometheus, Grafana, Loki, Promtail, and OpenTelemetry |
| Delivery | Docker, Docker Compose, Helm, and GitHub Actions |

AI proposes context. OPA authorizes. Temporal executes. CortexOps never treats an LLM output as an infrastructure mutation.

## Run the real local demo

### Prerequisites

- Docker Desktop or Docker Engine with Compose v2
- Go 1.25+
- `kubectl`
- `kind`
- Helm 3.14+ for chart validation

### 1. Create a disposable cluster

```bash
kind create cluster --name cortexops-dev
```

### 2. Start the platform

```bash
make dev-up
make bootstrap
make verify-runtime
```

### 3. Open the local control surfaces

| Surface | URL |
|---|---|
| Grafana | <http://localhost:3000> |
| Temporal | <http://localhost:8233> |
| Prometheus | <http://localhost:9090> |
| NATS monitor | <http://localhost:8222> |

Grafana, Temporal UI, and NATS monitoring use development-only defaults. Never expose this stack without authentication.

### 4. Trigger and recover a failure

```bash
make demo-failure SCENARIO=rollout-fail
make validate-pipeline
make demo-recovery
```

`validate-pipeline` exercises the live path from Kubernetes event to Temporal workflow and restores the demo workload afterward.

## Engineering confidence

Run the standard validation suite:

```bash
go test ./...
go vet ./...
make verify-runtime
make validate-pipeline
make chaos-test
make validate-temporal
```

Integration-tagged chaos tests keep external infrastructure and cluster dependencies out of ordinary unit runs:

```bash
go test -v -tags=integration ./test/e2e/suites -run '^TestChaos_'
```

Frontend:

```bash
cd frontend
npm ci
npm run lint
npm run build
```

## Repository map

```text
cortexops/
├── cmd/                    Service entrypoints
├── internal/              Correlation, RCA, remediation, topology, diagnostics
├── pkg/                   Reusable broker, telemetry, validation, and topology code
├── api/v1/                 Protocol and API contracts
├── build/docker/           Shared backend container build
├── deploy/
│   ├── compose/            Local Docker development environment
│   ├── helm/cortexops/     Kubernetes production chart
│   ├── argocd/             Existing GitOps application manifest
│   ├── grafana/            Dashboards and datasources
│   ├── loki/              Log collection
│   ├── prometheus/        Metrics and alert rules
│   └── nats/              Broker configuration
├── sandbox/workloads/      Disposable demo Kubernetes workloads
├── test/e2e/               Golden-path and chaos validation
└── docs/                   Architecture, operations, security, and readiness evidence
```

## CI and release pipeline

GitHub Actions validates:

- Go formatting, vet, build, race tests, and golangci-lint
- Next.js lint and production build
- Helm dependency resolution, lint, and rendering
- Docker images for pull requests and published images for configured pushes
- Dependency, container, secret, and static security scanning

Docker publishing requires repository secrets:

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

Pull-request Docker builds remain build-only even when publishing secrets are unavailable.

## Deploy

### Helm

```bash
helm dependency build deploy/helm/cortexops
helm lint deploy/helm/cortexops

helm upgrade --install cortexops deploy/helm/cortexops \
  --namespace cortexops --create-namespace \
  --set ingress.enabled=true \
  --set ingress.host=api.example.com \
  --set image.tag=v1.0.0
```

For prerequisites and secret setup, use [the self-hosted deployment runbook](docs/SELF_HOSTED_DEPLOYMENT.md).

### Argo CD

The repository currently contains an Argo CD Application manifest. Full GitOps implementation, environment promotion, image update automation, and delivery verification are planned as the next infrastructure milestone rather than implied by the current manifest.

## Security model

- Services run as UID 10001
- Containers use read-only root filesystems where chart settings allow
- Least-privilege service accounts and Kubernetes RBAC
- Per-service network policies
- Diagnostics APIs use bearer authentication
- NATS supports development credentials and production NKeys
- Every remediation remains dry-run and policy-gated before mutation
- Maintenance windows fail closed
- Kubernetes mutation is logged and should be approved in production workflows

Read [the production readiness assessment](docs/PRODUCTION_READINESS_FINAL.md) before using the project beyond development.

## Roadmap

- Authenticated hosted tenant model with organization and role boundaries
- Tenant-scoped database schemas and durable incident/audit history
- Live frontend backed by authenticated APIs
- Argo CD application promotion and image update automation
- Helm-backed Kubernetes restore and disaster-recovery testing
- Complete OIDC, secrets rotation, and multi-cluster credential management

## Documentation

- [Architecture overview](docs/ARCHITECTURE_OVERVIEW.md)
- [Feature matrix](docs/FEATURE_MATRIX.md)
- [Engineering decisions](docs/ENGINEERING_DECISIONS.md)
- [Runtime operations](docs/RUNTIME_OPERATIONS.md)
- [Debugging guide](docs/DEBUGGING_GUIDE.md)
- [Production readiness](docs/PRODUCTION_READINESS_FINAL.md)
- [Security hardening report](docs/SECURITY_HARDENING_REPORT.md)
- [Replay validation](docs/REPLAY_VALIDATION.md)

## Contributing

Contributions should keep the project honest: add tests, preserve deterministic behavior, and document unsupported capabilities rather than implying production readiness.

See [CONTRIBUTING.md](CONTRIBUTING.md), [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md), and [SECURITY.md](SECURITY.md).

<div align="center">

CortexOps is MIT licensed.

**Measure the cluster. Understand the incident. Gate the action.**

</div>
