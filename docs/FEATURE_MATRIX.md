# CortexOps Feature Matrix

This document provides an authoritative source of truth for platform capabilities, documenting implementation status and operational evidence.

| Feature               | Status      | Evidence            |
| --------------------- | ----------- | ------------------- |
| Telemetry Ingestion   | IMPLEMENTED | Collector + NATS JetStream integration processes incoming events at scale. |
| Topology Intelligence | IMPLEMENTED | Topology Service + Network Graph API accurately resolves asset dependencies. |
| Blast Radius Analysis | IMPLEMENTED | Topology Traversal algorithm isolates impacted downstream microservices. |
| Event Correlation     | IMPLEMENTED | Correlator Service successfully matches metrics, logs, and events into unified Incidents. |
| RCA Engine            | IMPLEMENTED | LLM + Qdrant pipeline generates deterministic root cause analysis. |
| OPA Governance        | IMPLEMENTED | Rego evaluation gates automated remediation against compliance policies. |
| Temporal Workflows    | IMPLEMENTED | Validated via execution history: Remediation service orchestrates actions via Temporal workers. |
| Replay Safety         | IMPLEMENTED | Duplicate event storms handled idempotently, verified via chaos injection. |
| Rollback Support      | IMPLEMENTED | Executor rollback routines successfully reverse state on workflow failure. |
