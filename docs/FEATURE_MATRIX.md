# CortexOps Feature Matrix

This document provides an authoritative source of truth for platform capabilities, documenting implementation status and operational evidence.

## Core Pipeline

| Feature | Status | Evidence |
| :--- | :--- | :--- |
| Telemetry Ingestion | ✅ IMPLEMENTED | Collector + NATS JetStream processes K8s events at scale via authenticated connections. |
| Topology Intelligence | ✅ IMPLEMENTED | In-memory directed graph with `pgx/v5` PostgreSQL persistence. Schema-isolated under `topology` schema. |
| Blast Radius Analysis | ✅ IMPLEMENTED | BFS traversal isolates downstream impacted services. Exposed via `/v1/topology/blast-radius/{id}`. |
| Event Correlation | ✅ IMPLEMENTED | Deterministic heuristic scoring groups events into `CorrelatedIncident` objects by TraceID, time, and topology. |
| RCA Engine | ✅ IMPLEMENTED | RAG pipeline (Qdrant + LLM) generates context-grounded root cause advisories. Advisory-only — cannot trigger mutations. |
| OPA Governance | ✅ IMPLEMENTED | Rego evaluation gates remediation against action type, namespace, and maintenance window policies. |
| Temporal Workflows | ✅ IMPLEMENTED | Durable workflow orchestration: Propose → Policy → Dry-Run → Execute → Verify → Rollback. |
| Replay Safety | ✅ IMPLEMENTED | NATS `Nats-Msg-Id` deduplication ensures exactly-once delivery. Validated via chaos injection. |
| Rollback Support | ✅ IMPLEMENTED | Executor rollback routines reverse state on workflow failure. Automated stabilization checks. |

## Security & Governance

| Feature | Status | Evidence |
| :--- | :--- | :--- |
| Non-Root Containers | ✅ IMPLEMENTED | UID 10001, `readOnlyRootFilesystem`, `drop: ALL` capabilities, `seccomp: RuntimeDefault`. |
| API Authentication | ✅ IMPLEMENTED | Bearer token middleware on `/v1/` endpoints. Health endpoints remain unauthenticated for probes. |
| NATS Authentication | ✅ IMPLEMENTED | User/password (dev), NKey with per-service pub/sub permissions (production). |
| Per-Service RBAC | ✅ IMPLEMENTED | Per-service `ServiceAccount`s. Only `collector` and `remediation` have K8s API bindings. |
| NetworkPolicies | ✅ IMPLEMENTED | Default-deny ingress per service. Explicit egress to infrastructure ports. |
| HTTP Security Headers | ✅ IMPLEMENTED | CSP, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy in Next.js. |
| Maintenance Windows | ✅ IMPLEMENTED | OPA `maintenance_active` rule blocks all automated remediation during maintenance. |

## Observability

| Feature | Status | Evidence |
| :--- | :--- | :--- |
| Prometheus Metrics | ✅ IMPLEMENTED | All 5 services expose `/metrics` on port 9091. `ServiceMonitor` for K8s-native scraping. |
| OpenTelemetry Traces | ✅ IMPLEMENTED | OTLP gRPC trace exporter configured in each service via `pkg/telemetry`. |
| Grafana Dashboards | ✅ IMPLEMENTED | Auto-provisioned Prometheus + Loki datasources and CortexOps operations dashboard. |
| Log Aggregation | ✅ IMPLEMENTED | Loki + Promtail ship and index container logs. Queryable via Grafana Explore. |
| Alerting Rules | ✅ IMPLEMENTED | Prometheus alerts for service down, high error rate, high latency, NATS backpressure. |
| Health Endpoints | ✅ IMPLEMENTED | `/health` and `/debug/healthz` on all services for K8s liveness/readiness probes. |

## Infrastructure & Deployment

| Feature | Status | Evidence |
| :--- | :--- | :--- |
| Helm Chart | ✅ IMPLEMENTED | Production-grade chart with HPA, PDB, anti-affinity, and `ServiceMonitor`. |
| ArgoCD GitOps | ✅ IMPLEMENTED | `Application` manifest with `selfHeal: true`, `prune: true`, retry backoff. |
| Database Migrations | ✅ IMPLEMENTED | `golang-migrate/v4` with versioned up/down SQL scripts. Per-service schema isolation. |
| Automated Backups | ✅ IMPLEMENTED | Hourly `pg_dump` + Qdrant snapshots via `backup-cron` container. 7-day retention. |
| Resource Limits | ✅ IMPLEMENTED | CPU/memory bounds on all services (Compose and Helm). |
| Network Isolation | ✅ IMPLEMENTED | Named Docker bridge (dev), Kubernetes `NetworkPolicy` (production). |
| Platform Integration | ✅ IMPLEMENTED | Traefik labels, Homepage auto-discovery, Uptime Kuma monitoring. |
