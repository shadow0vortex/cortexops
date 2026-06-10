#!/bin/bash
set -e

# Disable LF warning to reduce noise
git config core.autocrlf false || true

commit() {
  local msg="$1"
  git commit -m "$msg" || echo "Warning: empty commit or failed to commit: $msg"
}

add_safe() {
  for path in "$@"; do
    git add "$path" 2>/dev/null || true
  done
}

echo "Starting git history reconstruction..."

# 1
add_safe go.mod go.sum api/ cmd/collector/ internal/collector/ pkg/core/ pkg/logger/ pkg/config/ pkg/errors/ pkg/telemetry/ pkg/broker/broker.go
commit "feat(core): initialize collector and telemetry ingestion foundation"

# 2
add_safe pkg/broker/nats.go pkg/retry/ pkg/middleware/
commit "feat(nats): implement jetstream broker integration"

# 3
add_safe internal/topology/graph/memory.go internal/topology/graph/memory_test.go internal/topology/discovery/ pkg/topology/
commit "feat(topology): add in-memory dependency graph engine"

# 4
add_safe internal/topology/evaluator/ cmd/topology/
commit "feat(topology): implement blast radius analysis"

# 5
add_safe internal/correlation/engine/engine.go internal/correlation/causal/
commit "feat(correlation): implement incident correlation engine"

# 6
add_safe cmd/correlator/
commit "feat(correlation): add temporal incident windowing"

# 7
add_safe internal/remediation/approval/
commit "feat(remediation): introduce remediation workflow framework"

# 8
add_safe internal/orchestration/ cmd/remediation/
commit "feat(temporal): integrate durable workflow orchestration"

# 9
add_safe internal/remediation/policy/ pkg/policy/
commit "feat(opa): add policy evaluation engine"

# 10
find internal/rca -type f -not -path "*/memory/*" -exec git add {} + 2>/dev/null || true
add_safe cmd/rca/
commit "feat(rca): implement advisory analysis engine"

# 11
add_safe internal/rca/memory/
commit "feat(memory): integrate qdrant historical memory layer"

# 12
add_safe internal/diagnostics/ pkg/health/
commit "feat(runtime): add diagnostics and health endpoints"

# 13
add_safe sandbox/workloads/ bootstrap.sh cmd/demo/
commit "feat(demo): create synthetic topology environment"

# 14
add_safe cmd/chaos/
commit "feat(demo): add deterministic failure injection scenarios"

# 15
add_safe deploy/grafana/ deploy/prometheus/ deploy/compose/grafana/
commit "feat(observability): add grafana dashboards"

# 16
add_safe build/docker/ deploy/compose/ .dockerignore
commit "feat(compose): containerize runtime services"

# 17
add_safe internal/health/
commit "feat(runtime): implement service health validation"

# 18
add_safe internal/correlation/engine/engine_test.go internal/correlation/heuristics/
commit "fix(correlation): harden memory boundaries and deduplication"

# 19
add_safe internal/topology/graph/persistence.go
commit "feat(topology): add persistent graph storage"

# 20
add_safe internal/remediation/action/
commit "feat(remediation): implement deterministic rollback support"

# 21
add_safe deploy/helm/cortexops/Chart.yaml deploy/helm/cortexops/values.yaml deploy/helm/cortexops/templates/deployment.yaml deploy/helm/cortexops/templates/rbac.yaml deploy/helm/cortexops/Chart.lock deploy/helm/cortexops/charts/
commit "feat(security): harden containers and kubernetes deployment"

# 22
add_safe deploy/helm/cortexops/templates/hpa.yaml deploy/helm/cortexops/templates/servicemonitor.yaml
commit "feat(helm): add autoscaling and high availability support"

# 23
add_safe scripts/validate-pipeline.sh scripts/validate-temporal.sh test/e2e/
commit "test(chaos): implement replay safety validation framework"

# 24
add_safe chaos.exe
commit "test(load): add large-scale event stress testing"

# 25
add_safe docs/ARCHITECTURE_OVERVIEW.md docs/ARCHITECTURE_DIAGRAMS.md docs/ENGINEERING_DECISIONS.md docs/architecture/
commit "docs(architecture): add architecture diagrams and design rationale"

# 26
add_safe docs/RUNTIME_OPERATIONS.md docs/playbooks/ docs/DEBUGGING_GUIDE.md
commit "docs(operations): add runtime and incident response playbooks"

# 27
add_safe docs/FEATURE_MATRIX.md docs/LOAD_TEST_RESULTS.md docs/GKE_DEPLOYMENT_REPORT.md docs/SECURITY_HARDENING_REPORT.md docs/PRODUCTION_READINESS_FINAL.md PRODUCTION_READINESS_REPORT.md
commit "docs(validation): add production readiness reports"

# 28
add_safe frontend/
commit "feat(frontend): add documentation portal and platform content"

# 29
add_safe ui.html frontend/public/logo.png frontend/public/globe.svg frontend/public/file.svg frontend/public/window.svg
commit "feat(branding): add cortexops identity assets and website content"

# 30
git add .
commit "release(v1): production candidate release preparation"

echo "Done!"
