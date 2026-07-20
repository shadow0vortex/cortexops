# Architecture Overview

CortexOps is built as a set of decoupled microservices, each with a single responsibility, communicating over NATS JetStream — a high-performance, persistent message bus.

## Distributed Components

### 1. Collector Service (`cmd/collector`)
- **Role**: Telemetry ingestion from the Kubernetes API Server.
- **Tech**: `client-go` Informers with authenticated NATS publishing.
- **Output**: Strongly-typed `TelemetryEnvelope` (Protobuf) published to JetStream.
- **Identity**: Dedicated `ServiceAccount` with read-only K8s RBAC bindings.

### 2. Topology Service (`cmd/topology`)
- **Role**: Real-time dependency graph of cluster resources.
- **Tech**: In-memory Directed Graph with asynchronous PostgreSQL persistence via `pgx/v5` connection pooling. Schema-isolated under the `topology` database schema.
- **Capability**: BFS traversal for blast-radius analysis with resilient crash-recovery via `golang-migrate` versioned migrations.
- **API**: Versioned REST API (`/v1/topology/nodes`, `/v1/topology/blast-radius/{id}`) protected by Bearer token authentication (`DIAG_API_TOKEN`).

### 3. Correlation Engine (`cmd/correlator`)
- **Role**: Causal incident detection from raw telemetry streams.
- **Tech**: Deterministic heuristic scoring using TraceIDs, temporal windowing, and topology proximity.
- **Output**: `CorrelatedIncident` published to NATS.
- **Guarantee**: Fully deterministic — identical inputs always produce identical outputs.

### 4. RCA Service (`cmd/rca`)
- **Role**: Root cause analysis and advisory report generation.
- **Tech**: RAG pipeline using Qdrant Vector DB for historical similarity search + LLM integration for context-grounded advisory reports.
- **Safety**: Advisory-only — the RCA service cannot trigger infrastructure mutations.

### 5. Remediation Service (`cmd/remediation`)
- **Role**: Policy-governed workflow orchestration and Kubernetes mutation.
- **Tech**: Temporal durable workflows with OPA policy evaluation.
- **Identity**: Dedicated `ServiceAccount` with scoped K8s RBAC bindings (pod delete, deployment patch).
- **Safety**: Dry-Run → Policy Gate → Execute → Verify → Rollback.
- **Governance**: OPA Rego policies enforce namespace protection, action type allow-listing, and maintenance window blocking.

## Infrastructure Dependencies

| Component | Version | Purpose |
|-----------|---------|---------|
| NATS JetStream | `2.10-alpine` | Event bus with exactly-once delivery |
| PostgreSQL | `17-alpine` | Topology persistence, Temporal state |
| Qdrant | `v1.8.0` | Vector similarity search for RCA |
| Temporal | `1.23.0` | Durable workflow orchestration |
| Prometheus | `v2.51.0` | Metrics collection and alerting |
| Grafana | `10.4.0` | Dashboards and log exploration |
| Loki | `2.9.2` | Centralized log aggregation |
| Promtail | `2.9.2` | Container log shipping |

## Communication Matrix

| Source | Destination | Protocol | Port | Purpose |
| :--- | :--- | :--- | :--- | :--- |
| Collector | NATS | JetStream | 4222 | Publish telemetry events |
| Correlator | NATS | JetStream | 4222 | Consume telemetry, publish incidents |
| Correlator | Topology | HTTP | 9091 | Query dependency graph |
| RCA | NATS | JetStream | 4222 | Consume incidents, publish reports |
| RCA | Qdrant | HTTP | 6333 | Vector similarity search |
| Remediation | NATS | JetStream | 4222 | Consume RCA reports |
| Remediation | Temporal | gRPC | 7233 | Workflow persistence and scheduling |
| Remediation | K8s API | HTTPS | 443 | Infrastructure mutation |
| All Services | Prometheus | HTTP | 9091 | `/metrics` scraping |

## Network Security

All inter-service communication is governed by Kubernetes `NetworkPolicy` resources:
- **Default deny** ingress on all CortexOps pods.
- Explicit ingress allow from `app: cortexops` pods and the `monitoring` namespace (Prometheus).
- Explicit egress allow to DNS (53), NATS (4222), PostgreSQL (5432), Temporal (7233), Qdrant (6333).

## Authentication

| Layer | Mechanism | Scope |
|-------|-----------|-------|
| NATS (dev) | User/password | All services |
| NATS (prod) | NKey per-service identity | Per-service pub/sub permissions |
| Diagnostics API | Bearer token (`DIAG_API_TOKEN`) | Data endpoints only |
| Temporal UI | BasicAuth via Traefik middleware | Dashboard access |
| K8s API | Per-service `ServiceAccount` RBAC | `collector` and `remediation` only |
