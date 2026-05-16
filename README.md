# CortexOps

**Autonomous Kubernetes Operations Platform**

CortexOps is an event-driven, topology-aware control plane for Kubernetes. It ingests infrastructure telemetry, builds deterministic causal chains, provides RAG-grounded advisory Root Cause Analysis (RCA), and safely orchestrates policy-governed remediations via Temporal.

Designed for Site Reliability Engineering (SRE) teams, CortexOps addresses the MTTR (Mean Time To Resolution) crisis in microservice architectures by bridging the gap between passive observability and autonomous self-healing, without compromising operational safety.

---

## 1. Project Overview

Modern cloud-native systems generate overwhelming telemetry. When a node fails, it cascades into hundreds of disparate alerts across pods, services, and ingresses. Traditional monitoring tools surface symptoms; CortexOps determines causation.

CortexOps exists to:
1. **Normalize** disjointed Kubernetes events and metrics into a unified, deterministic telemetry stream.
2. **Correlate** symptoms using a real-time, in-memory topology graph to calculate exact blast radii.
3. **Analyze** root causes using bounded, read-only AI enrichment grounded strictly in historical telemetry.
4. **Remediate** safely using Temporal workflows governed by OPA (Open Policy Agent) and human-in-the-loop Slack approvals.

**Architecture Philosophy:**
- **Determinism Over Magic:** The correlation engine is a math-driven state machine. AI is strictly confined to an advisory enrichment layer. AI cannot execute commands.
- **Replay Safety:** The entire pipeline relies on NATS JetStream and Temporal. If CortexOps crashes mid-remediation, it recovers exactly where it left off without duplicating operations.
- **Fail-Closed Governance:** Every remediation is dry-run against the K8s API. If an anomaly is detected, or if stabilization verification fails, the orchestrator automatically rolls back the mutation.

---

## 2. Core Features

- **Telemetry Ingestion Engine:** Connects to K8s Informers, flattening events into strongly-typed Protobuf envelopes.
- **Topology Intelligence Layer:** Maintains a real-time directed graph of cluster dependencies, enabling blast-radius traversal (e.g., Node -> ReplicaSet -> Pod -> PVC).
- **Deterministic Event Correlation:** Groups temporally proximate events based on topological affinity into causal chains.
- **AI-Assisted RCA (Advisory):** Retrieves historically similar incidents via Qdrant Vector DB to generate context-grounded analysis without hallucination.
- **Policy-Governed Remediation:** Executes `POD_RESTART`, `DEPLOYMENT_ROLLOUT`, or `HORIZONTAL_SCALE` via `go.temporal.io/sdk`, strictly bound by OPA allowlists.
- **Runtime Diagnostics:** Continuous goroutine/heap leak detection with automated `pprof` stack dumping.
- **Chaos Validation:** First-class testing harness for asserting system invariants during network partitions and pod evictions.

---

## 3. Architecture Overview

CortexOps utilizes a microservice architecture built strictly in Go.

```text
[ K8s API ] ---> ( Informers ) ---> [ Collector ]
                                        | (Protobuf / OTel)
                                        v
                                [ NATS JetStream ]
                                        | (Consumer Group)
[ OPA Engine ] <----+                   v
                    |           [ Correlation Engine ] <---> [ Topology Graph ]
                    |                   |
                    |                   v
                    +-------- [ Remediation Engine (Temporal) ] ---> [ K8s Mutator ]
                                        |
                                        v
                            [ Qdrant Vector DB (RCA) ]
```

### Component Roles
- **NATS JetStream:** Provides durability, backpressure, and replay semantics for all telemetry.
- **Topology Graph:** In-memory `client-go` informer cache providing sub-millisecond dependency resolution.
- **Temporal:** Orchestrates long-running remediation lifecycles, ensuring idempotency, retries, and rollback safety.
- **Qdrant:** Stores incident embeddings to power historical similarity search for the RCA engine.

---

## 4. Repository Structure

CortexOps is a Go monorepo utilizing Go Workspaces and Protobuf definitions.

```text
├── api/v1/                 # Protobuf definitions (events.proto, remediation.proto)
├── cmd/cortexops/          # Main binary entrypoint
├── internal/
│   ├── collector/          # K8s event watchers and metric scrapers
│   ├── correlation/        # Temporal windowing and causal scoring engine
│   ├── topology/           # Directed graph of K8s dependencies
│   ├── rca/                # RAG AI engine, Token truncation, Qdrant client
│   ├── remediation/        # OPA policy eval, K8s mutator, slack approvals
│   ├── orchestration/      # Temporal workflows and activities
│   └── diagnostics/        # Leak detection, pprof dumps, introspection API
├── pkg/
│   └── core/               # Shared Dependency Injected Interfaces (LLM, VectorStore)
├── deploy/                 # Helm charts, Grafana dashboards, Dockerfiles
└── test/e2e/               # Isolated Chaos and Replay validation harness
```

---

## 5. Operational Safety Model

CortexOps is designed for hostile environments. Safety is enforced deterministically.

1. **AI Governance:** The AI layer (`LLMClient`) is strictly decoupled from the `RemediationExecutor`. The AI output is an immutable `RCAReport`. It cannot trigger infrastructure actions.
2. **Blast-Radius Enforcement:** The Topology graph calculates impact depth. If an incident's blast radius impacts >3 hierarchical tiers, the OPA risk score dictates a mandatory human-in-the-loop approval.
3. **Immutable Audit Trails:** Every transition (e.g., `DRY_RUNNING -> EXECUTING`) emits an `AuditRecord` mapped to a TraceID, stored permanently.
4. **Rollback Safety:** If a deployment rollout fails readiness probes within 5 minutes, Temporal executes a rollback activity to revert the patch.

---

## 6. Local Development Setup

### Prerequisites
- Go 1.22+
- Docker & Docker Compose
- Kind or Minikube
- `protoc` (Protocol Buffers Compiler)

### Bootstrap Environment

1. **Clone & Install Dependencies**
   ```bash
   git clone https://github.com/cortexops/cortexops.git
   cd cortexops
   go mod tidy
   ```

2. **Generate Protobufs**
   ```bash
   make proto
   ```

3. **Start Infrastructure Services (Temporal, NATS, Qdrant, Postgres)**
   ```bash
   docker-compose -f deploy/docker-compose.yaml up -d
   ```

4. **Start Kind Cluster**
   ```bash
   kind create cluster --name cortexops-dev
   ```

---

## 7. Running CortexOps Locally

To run the CortexOps engine against your local Kind cluster:

```bash
# Export configuration
export KUBECONFIG=~/.kube/config
export NATS_URL="nats://localhost:4222"
export TEMPORAL_URL="localhost:7233"
export QDRANT_URL="http://localhost:6333"

# Run the core engine
make run
```

### Accessing Observability
- **Grafana:** `http://localhost:3000` (View Platform Health & Goroutine Profiler)
- **Temporal UI:** `http://localhost:8233` (View Remediation Workflow state)

---

## 8. End-to-End Demo Scenarios

### Scenario: Failed Deployment Rollout
Trigger a rollout utilizing a broken image tag:
```bash
kubectl create deployment nginx --image=nginx:1.24.0
# Intentionally break it
kubectl set image deployment/nginx nginx=nginx:nonexistent
```

**Expected CortexOps Behavior:**
1. **Telemetry:** Ingests `BackOff` and `FailedPullImage` K8s events.
2. **Correlation:** Groups events by the `nginx` ReplicaSet topology.
3. **RCA:** Qdrant flags similar historical pull errors.
4. **Remediation:** Temporal evaluates OPA. Action triggers a `DEPLOYMENT_ROLLOUT_RESTART`.
5. **Rollback:** Because the image is still broken, verification fails. Temporal reverts the annotation patch.

---

## 9. Runtime Diagnostics & Debugging

CortexOps exposes a read-only introspection API on port `9091` for SREs.

**Inspect Topology State:**
```bash
curl -X GET "http://localhost:9091/debug/graph/node?id=pod/nginx-1234"
```

**Memory & Goroutine Leaks:**
The background `Profiler` monitors stability. If goroutines exceed 5000, it automatically generates a snapshot in `/tmp/cortexops/dumps/`.
To manually capture a profile:
```bash
curl -sK -v http://localhost:6060/debug/pprof/heap > heap.out
go tool pprof heap.out
```

---

## 10. Chaos Engineering

We validate resilience using the `test/e2e/chaos` harness. 

**Run Replay Validation:**
```bash
go test ./test/e2e/suites -run TestReplayIdempotency -v
```
**Expectation:** The test simulates a NATS network partition. Upon reconnect, CortexOps processes the burst. The harness asserts that exactly **one** incident is generated due to event deduplication.

---

## 11. Observability

All services expose Prometheus `/metrics`. 
Key SLIs:
- `cortexops_telemetry_normalization_duration_seconds`
- `cortexops_ai_inference_latency_seconds`
- `cortexops_remediation_rollback_total`
- `cortexops_platform_health_score` (0.0 - 1.0)

*(Grafana JSON definitions are located in `deploy/grafana/dashboards/`)*

---

## 12. Production Deployment

Deploy CortexOps into your K8s cluster via the hardened Helm chart.

```bash
helm install cortexops ./deploy/helm/cortexops \
  --namespace cortexops-system \
  --create-namespace \
  --set replicaCount=3
```

**Security Assumptions:**
- Runs as non-root.
- Read-only root filesystem.
- `PodDisruptionBudget` configured (`minAvailable: 2`).
- RBAC is strictly scoped. No cluster-admin rights are required.

---

## 13. Reliability Guarantees

- **Idempotent Executions:** K8s API patches executed by Temporal activities are fundamentally idempotent. 
- **Degraded Modes:** If Qdrant or the LLM is down, CortexOps sets `IsDegraded=true` and passes raw telemetry to the remediation engine. The incident pipeline **does not halt**.
- **Replay Semantics:** Correlation buckets act on Event Timestamps, not ingestion timestamps. Replaying a 7-day old NATS stream will yield the exact same historical graph.

---

## 14. Development Workflow

```bash
make proto      # Recompile gRPC/Protobuf files
make lint       # Run staticcheck and golangci-lint
make test       # Run unit and race-detector tests
make e2e        # Run isolated Kubernetes integration harness
```

---

## 15. Roadmap

- [ ] **Phase 8:** Multi-cluster federation via NATS Leaf Nodes.
- [ ] **Phase 9:** Integration with Istio/Envoy metrics for network-layer blast radius analysis.
- [ ] **Phase 10:** Configurable SLO-driven remediation policies.

---

## 16. Contributing Guide

1. All new remediation activities MUST be integrated with Temporal `DryRun` steps.
2. Global mutable state is strictly prohibited. Use Dependency Injection (`pkg/core`).
3. PRs must pass `go test -race` and include an architecture update if state machines are altered.

---

## 17. License & Disclaimer

**License:** MIT

**Operational Safety Disclaimer:**
CortexOps modifies Kubernetes state. While governed by OPA and Temporal safeties, autonomous remediation carries inherent risks. Ensure rigorous dry-run testing in staging environments before authorizing `EXECUTING` states in production clusters.
