# Architecture Overview

CortexOps is built as a set of decoupled microservices, each with a single responsibility, communicating over a high-performance message bus.

## Distributed Components

### 1. Collector Service
- **Role**: Ingestion.
- **Tech**: `client-go` Informers.
- **Output**: Strongly-typed `TelemetryEnvelope` (Protobuf).

### 2. Topology Service
- **Role**: Dependency Graph.
- **Tech**: In-memory Directed Graph + Diagnostics API.
- **Capability**: BFS traversal for blast-radius analysis.

### 3. Correlation Engine
- **Role**: Causality Detection.
- **Tech**: Heuristic Scoring + Temporal Windowing.
- **Output**: `CorrelatedIncident`.

### 4. RCA Service
- **Role**: AI Analysis.
- **Tech**: RAG + Qdrant Vector DB + LLM.
- **Capability**: Context-grounded advisory reports.

### 5. Remediation Service
- **Role**: Orchestration & Mutator.
- **Tech**: Temporal Workflows + OPA Engine.
- **Safety**: Dry-Run -> Execute -> Verify -> Rollback.

## Communication Matrix

| Source | Destination | Protocol | Purpose |
| :--- | :--- | :--- | :--- |
| Collector | NATS | JetStream | Publish Telemetry |
| Correlator | Topology | HTTP | Query Dependencies |
| Correlator | NATS | JetStream | Consume Telemetry / Publish Incidents |
| RCA | Qdrant | HTTP/gRPC | Similarity Search |
| Remediation | Temporal | gRPC | Workflow Persistence |
| Remediation | K8s API | HTTPS | Infrastructure Mutation |
