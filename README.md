# CortexOps

**Autonomous Kubernetes Operations & Remediation Control Plane**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/shadow0vortex/cortexops)](https://goreportcard.com/report/github.com/shadow0vortex/cortexops)
[![Tests](https://github.com/shadow0vortex/cortexops/actions/workflows/ci.yaml/badge.svg)](https://github.com/shadow0vortex/cortexops/actions/workflows/ci.yaml)

CortexOps is an event-driven, topology-aware control plane extension for Kubernetes. It ingests infrastructure telemetry, builds deterministic causal chains, provides RAG-grounded advisory Root Cause Analysis (RCA), and safely orchestrates policy-governed remediations via Temporal.

Designed for high-scale Site Reliability Engineering (SRE) teams, CortexOps addresses the MTTR crisis by bridging the gap between passive monitoring and autonomous self-healing—without sacrificing operational safety.

---

## 🚀 Key Features

*   **Telemetry Normalization**: Ingests disparate K8s events and metrics into strongly-typed Protobuf envelopes.
*   **Topology Intelligence**: Real-time directed graph maintaining cluster dependencies for sub-millisecond blast-radius traversal.
*   **Causal Correlation**: Heuristic-based scoring engine that groups symptoms into incidents using TraceIDs, time, and topology.
*   **Advisory RAG RCA**: Context-grounded AI analysis retrieving historical patterns from Qdrant Vector DB.
*   **Durable Remediation**: Temporal-orchestrated workflows for `POD_RESTART`, `ROLLOUT_RESTART`, and `SCALING`, with automatic verification and rollback.
*   **Governance by OPA**: Every action is validated against Open Policy Agent allowlists before execution.

---

## 🏗 Architecture

CortexOps is built as a set of decoupled Go microservices.

```mermaid
graph LR
    K8s[Kubernetes API] --> Collector
    Collector --> NATS[NATS JetStream]
    NATS --> Correlator
    Correlator --> Topology[Topology Service]
    Correlator --> RCA[RCA Service]
    RCA --> Qdrant[(Qdrant DB)]
    Correlator --> Rem[Remediation Service]
    Rem --> Temporal[[Temporal Orchestrator]]
    Rem --> K8s
```

*For more details, see our [Architecture Diagram Suite](docs/architecture/telemetry_ingestion.md).*

---

## 🚦 Quick Start

### 1. Start Infrastructure & Services
Spin up the full stack (NATS, Temporal, Qdrant, Postgres, Grafana + CortexOps Microservices) in one command:
```bash
make dev-up
```

### 2. Verify System Health
```bash
make verify-runtime
```

### 3. Bootstrap Demo Environment
Deploy the synthetic microservice topology into your cluster:
```bash
make bootstrap
```

### 4. Inject a Failure Scenario
```bash
make demo-failure SCENARIO=rollout-fail
```

### 5. Observe Results
- **Grafana**: `http://localhost:3000` (Incident Tracking)
- **Temporal UI**: `http://localhost:8233` (Remediation Lifecycle)
- **Diagnostics API**: `make diagnostics`

---

## 🛡 Operational Guarantees

### Determinism Over "Magic"
The correlation engine is a math-driven state machine. AI is strictly confined to an advisory enrichment layer and cannot trigger infrastructure mutations.

### Replay Safety
Utilizing NATS JetStream deduplication and event-timestamp-driven windowing, replaying identical telemetry bursts always yields the same causal results.

### Fail-Closed Governance
Every remediation is dry-run against the K8s API. If an anomaly is detected, or if stabilization verification fails, the orchestrator automatically rolls back the mutation.

---

## 📂 Documentation

*   [**Demo Guide**](docs/DEMO_GUIDE.md): Detailed walkthrough of the platform in action.
*   [**Architecture Overview**](docs/ARCHITECTURE_OVERVIEW.md): Component roles and communication matrix.
*   [**Engineering Decisions**](docs/ENGINEERING_DECISIONS.md): Why we built it this way.
*   [**Interview Talking Points**](docs/INTERVIEW_TALKING_POINTS.md): Technical deep-dive for recruiters and engineers.
*   [**Operational Playbooks**](docs/playbooks/INCIDENT_RESPONSE.md): Runbooks for platform maintenance.

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md).

### Security
To report a vulnerability, please refer to our [Security Policy](SECURITY.md).

---

## 📜 License

MIT License. See [LICENSE](LICENSE) for details.

---

**Operational Safety Disclaimer**: CortexOps modifies Kubernetes state. While governed by OPA and Temporal safeties, autonomous remediation carries inherent risks. Ensure rigorous dry-run testing before authorizing `EXECUTING` states in production.
