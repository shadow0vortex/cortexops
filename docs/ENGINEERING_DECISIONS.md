# Engineering Decisions

This document tracks the foundational design choices made during the development of CortexOps.

## 1. Decoupled Microservices vs. Monolith
- **Decision**: Separate binaries for each service (Collector, Correlator, Topology, etc.).
- **Rationale**: Enables independent scaling, failure isolation, and enforces clear interface boundaries. It also reflects real-world Kubernetes-native patterns.

## 2. In-Memory Graph for Topology
- **Decision**: Thread-safe in-memory graph store backed by asynchronous PostgreSQL persistence and SHA-256 state hashing.
- **Rationale**: Prioritizes sub-millisecond query performance for real-time correlation while ensuring resilient, crash-safe state recovery. Asynchronous snapshotting prevents database IO overhead from blocking the fast-path correlation pipeline.

## 3. Heuristic Scoring over Probabilistic Models
- **Decision**: Weighted scoring based on TraceIDs, time, and topology.
- **Rationale**: Determinism is non-negotiable for autonomous control planes. Heuristics are 100% auditable and reproducible, unlike probabilistic black-box models.

## 4. Advisory-Only AI
- **Decision**: AI produces `RCAReports` but cannot call `Execute()`.
- **Rationale**: Safety. The LLM acts as an "expert advisor" for the human SRE, while the "autopilot" (Correlation + OPA) handles the safe, deterministic actions.

## 5. Exactly-Once Delivery (NATS + SeqID)
- **Decision**: Utilizing NATS JetStream with mandatory `Nats-Msg-Id`.
- **Rationale**: Infrastructure events can be noisy and duplicated. Deduplication at the broker layer is essential for maintaining the integrity of the causal chain.
