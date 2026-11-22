# CortexOps Production Readiness Assessment

## Current verdict

**Not production ready as a hosted, multi-tenant product.** The repository contains a useful Kubernetes operations prototype and a working local integration demo, but it does not yet provide the identity, tenant isolation, durable operational data, or live user dashboard required for real customer deployments.

The previous version of this report claimed production readiness and included performance/security ratings that were not independently reproducible. Those claims have been removed; do not use them as deployment or resume evidence.

## What was verified locally

Using Docker Compose and a local `kind` cluster:

| Check | Result |
|---|---|
| `docker compose -f deploy/compose/docker-compose.dev.yaml --profile full config --quiet` | Passed |
| `go test ./...` | Passed with the local Kubernetes context available |
| `go vet ./...` | Passed |
| `npm run lint` | Passed with 18 existing warnings |
| `npm run build` | Passed; Next.js warned about an additional lockfile outside the repository |
| `make verify-runtime` | Passed for all Compose runtime services |
| `make diagnostics` | Passed; returned live topology data |
| `make validate-pipeline` | Passed; demo Kubernetes events reached collector, correlator, RCA, remediation, and a Temporal workflow |
| `helm lint deploy/helm/cortexops` | Not run; Helm was not installed in the validation environment |

These are local development checks, not a production deployment, load test, security certification, or proof of service-level objectives. The e2e test package also contains mocked assertions and should not be treated as evidence of full system behavior by itself.

## Production blockers

### Identity and tenant isolation

- The diagnostics API uses one optional shared bearer token; if it is unset, protected routes are unauthenticated.
- There is no OIDC login, user/account model, organization membership, role-based authorization, or tenant-scoped API/data access.
- The frontend has no authenticated connection to the backend.

### Durable product data

- Remediation audit records are held in an in-memory `MemoryAuditStore`.
- There is no durable, tenant-scoped incident/action history or retention model.
- The topology graph is in memory in the local Compose setup; persistence and restore behavior have not been verified as a recovery procedure.

### Real user-facing monitoring

- The observability page displays hard-coded sample values and simulated charts.
- The Kanban board is not backed by live incident/workflow APIs.
- Topology discovery currently watches Pods and records a limited scheduled-on relationship; it does not yet build a complete service/deployment dependency graph.

### Safe production remediation

- Remediation logic contains demo-specific target selection, including the `cortexops-demo/demo-frontend` target.
- Automated execution, approval, rollback, and tenant authorization have not been validated against production-grade multi-tenant policies.
- The Compose configuration is deliberately development-only. It uses local kubeconfig access and a TLS-verification bypass; do not deploy it or reuse those settings in production.

### Operational evidence

- Helm validation, backup restoration, disaster recovery, sustained-load tests, and an external identity-provider flow were not run in this assessment.
- The current test suite does not establish tenant-isolation, authorization, or recovery guarantees.
- The repository's earlier scorecards and benchmark numbers should be considered unverified until rerun from reproducible scripts with raw results.

## Recommended delivery gates

1. **Identity and authorization:** add OIDC/JWT validation, organization membership, roles, explicit tenant context, and deny-by-default API authorization. Missing production identity configuration must fail startup.
2. **Tenant-safe persistence:** add migrations for organizations, users/memberships, clusters, incidents, remediation actions, and append-only audit records. Enforce tenant scoping in every query and test cross-tenant access denial.
3. **Live observability:** replace demo numbers and board state with authenticated APIs backed by actual telemetry, incident, topology, and Temporal data. Show loading, empty, stale, and error states.
4. **Governed remediation:** remove fixed demo targets; require validated target scope, policy decision, approval where needed, durable audit, idempotency, and tested rollback before enabling execution.
5. **Release proof:** split unit/integration/cluster tests, run CI against a disposable cluster, lint Helm, test migrations and restore, scan immutable images, and publish measured SLO/load/security evidence.

Until these gates pass, describe CortexOps as an **active-development Kubernetes operations prototype with a locally verified integration demo**, not an enterprise-ready hosted service.
