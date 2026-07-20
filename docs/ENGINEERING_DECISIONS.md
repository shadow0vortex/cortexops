# Engineering Decisions

This document tracks the foundational design choices made during the development of CortexOps. Each decision is recorded as an Architecture Decision Record (ADR).

## ADR-001: Decoupled Microservices vs. Monolith
- **Decision**: Separate binaries for each service (Collector, Correlator, Topology, RCA, Remediation).
- **Rationale**: Enables independent scaling, failure isolation, and enforces clear interface boundaries. It also reflects real-world Kubernetes-native patterns.
- **Trade-off**: Increased operational complexity managed through Docker Compose profiles and Helm chart templating.

## ADR-002: In-Memory Graph for Topology
- **Decision**: Thread-safe in-memory graph store backed by asynchronous PostgreSQL persistence via `pgx/v5` and SHA-256 state hashing.
- **Rationale**: Prioritizes sub-millisecond query performance for real-time correlation while ensuring resilient, crash-safe state recovery. Asynchronous snapshotting prevents database IO overhead from blocking the fast-path correlation pipeline.
- **Migration**: Schema managed by `golang-migrate/migrate/v4` with versioned up/down SQL scripts.

## ADR-003: Heuristic Scoring over Probabilistic Models
- **Decision**: Weighted scoring based on TraceIDs, time, and topology.
- **Rationale**: Determinism is non-negotiable for autonomous control planes. Heuristics are 100% auditable and reproducible, unlike probabilistic black-box models.

## ADR-004: Advisory-Only AI
- **Decision**: AI produces `RCAReports` but cannot call `Execute()`.
- **Rationale**: Safety. The LLM acts as an "expert advisor" for the human SRE, while the "autopilot" (Correlation + OPA) handles the safe, deterministic actions.

## ADR-005: Exactly-Once Delivery (NATS + SeqID)
- **Decision**: Utilizing NATS JetStream with mandatory `Nats-Msg-Id`.
- **Rationale**: Infrastructure events can be noisy and duplicated. Deduplication at the broker layer is essential for maintaining the integrity of the causal chain.

## ADR-006: pgx/v5 over lib/pq
- **Decision**: Replaced `lib/pq` with `jackc/pgx/v5` as the PostgreSQL driver.
- **Rationale**: `pgx/v5` provides native connection pooling, binary protocol support, and superior performance. `lib/pq` is in maintenance mode and lacks built-in pool management.
- **Impact**: Connection exhaustion under load is eliminated; the pool automatically manages lifecycle.

## ADR-007: Per-Service ServiceAccounts and RBAC
- **Decision**: Each microservice receives its own Kubernetes `ServiceAccount`. Only `collector` and `remediation` receive RBAC bindings.
- **Rationale**: Least-privilege principle. Services like `topology`, `correlator`, and `rca` have zero Kubernetes API access, reducing the blast radius of a compromised pod.
- **Implementation**: Helm `rbac.yaml` template conditionally binds roles only to the two services that need K8s API access.

## ADR-008: OPA Maintenance Window Enforcement
- **Decision**: The OPA Rego policy engine evaluates `input.maintenance_window` and blocks all remediation actions when `true`.
- **Rationale**: Automated remediation during scheduled maintenance is dangerous and can conflict with human operators performing manual changes. A hard policy gate prevents this.
- **Trade-off**: Requires the calling workflow to correctly inject the `maintenance_window` field. If the field is omitted, it defaults to `false` (remediation allowed).

## ADR-009: Non-Root Container Runtime
- **Decision**: All containers run as UID 10001 (`cortexops` user) with `readOnlyRootFilesystem`, `drop: ALL` capabilities, and `seccomp: RuntimeDefault`.
- **Rationale**: Prevents privilege escalation, filesystem tampering, and restricts syscall surface. The `/tmp` directory is provided via `emptyDir` for write operations.

## ADR-010: API Versioning (`/v1/` Prefix)
- **Decision**: All data-serving API endpoints are prefixed with `/v1/`. Health and probe endpoints remain unversioned.
- **Rationale**: Enables backward-compatible API evolution. New versions (`/v2/`) can coexist during migration windows without breaking existing integrations.
- **Convention**: Health → `/health`, `/debug/healthz`. Data → `/v1/topology/nodes`, `/v1/topology/blast-radius/{id}`.

## ADR-011: Bearer Token Authentication on Diagnostics API
- **Decision**: Data endpoints require a `Bearer` token (`DIAG_API_TOKEN`). Health endpoints remain unauthenticated for K8s probe compatibility.
- **Rationale**: The diagnostics API exposes sensitive topology and blast-radius data. Without authentication, any pod in the cluster could exfiltrate infrastructure intelligence.
- **Dev Mode**: If `DIAG_API_TOKEN` is unset, authentication is bypassed for local development ergonomics.

## ADR-012: ArgoCD GitOps with Self-Healing
- **Decision**: An ArgoCD `Application` manifest with `selfHeal: true` and `prune: true` is provided as the recommended deployment method.
- **Rationale**: GitOps ensures that the cluster state is always consistent with the repository. Self-healing automatically corrects manual drift, and pruning removes orphaned resources.

## ADR-013: Per-Service PostgreSQL Schemas
- **Decision**: Each microservice owns its own PostgreSQL schema (`topology`, `correlator`, `rca`, `remediation`) within the shared `cortexops` database.
- **Rationale**: Schema isolation prevents data collision between services, enables independent schema evolution, and provides a natural boundary for future database-per-service migration if needed.
- **Migration**: Managed via `000002_per_service_schemas.up.sql` / `000002_per_service_schemas.down.sql`.

## ADR-014: NATS NKey Authentication for Production
- **Decision**: Production deployments use NATS NKey cryptographic identity with per-service publish/subscribe permissions. Development environments use simple user/password authentication.
- **Rationale**: NKeys provide non-revokable, password-less authentication tied to ed25519 key pairs. Per-service topic permissions enforce least-privilege messaging (e.g., the RCA service cannot publish to remediation topics).
