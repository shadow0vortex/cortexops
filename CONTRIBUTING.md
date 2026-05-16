# Contributing to CortexOps

We welcome contributions from SREs, Kubernetes engineers, and distributed systems practitioners.

## Core Engineering Principles

If you are opening a PR, you MUST adhere to the following principles:

1. **Determinism over Magic:** If you add a feature to the Correlation Engine, it must be mathematically deterministic. We do not accept PRs that introduce opaque ML models into the critical path.
2. **Zero Global Mutable State:** Ensure your code is dependency-injected. Pass `ctx` explicitly. Use interfaces in `pkg/core`.
3. **Replay Safety:** Your code must behave identically if the exact same event stream is replayed 5 times from NATS.
4. **Idempotent Actions:** If you introduce a new Remediation type in `api/v1/remediation.proto`, the corresponding K8s implementation MUST be idempotent.

## Development Workflow

1. Fork the repo and create your branch from `main`.
2. Run `make proto` if you modify gRPC contracts.
3. Ensure `go test -race ./...` passes.
4. Run the validation harness: `go test ./test/e2e/...`

## Architectural Reviews

Major changes require an RFC document submitted to `docs/RFCs/`. State your Failure Assumptions clearly.
