#!/usr/bin/env bash
set -eo pipefail

echo "========================================="
echo " CortexOps Golden Path Validation"
echo "========================================="

# Detect OS to handle Docker logs and paths properly
if [[ "$OS" == "Windows_NT" ]]; then
    DOCKER_CMD="docker.exe"
    MAKE_CMD="make"
else
    DOCKER_CMD="docker"
    MAKE_CMD="make"
fi

echo "[1/4] Bootstrapping Demo Workload..."
$MAKE_CMD bootstrap

# Wait for workload
echo "Waiting for demo-frontend pods to become ready..."
kubectl wait --for=condition=Available deployment/demo-frontend -n cortexops-demo --timeout=60s || true

# Save the current time so we only scan new logs
START_TIME=$($DOCKER_CMD exec compose-remediation-1 date -u +"%Y-%m-%dT%H:%M:%SZ" || date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "[2/4] Injecting Rollout Failure..."
$MAKE_CMD demo-failure SCENARIO=rollout-fail

echo "Waiting for failure to propagate (15s)..."
sleep 15

echo "[3/4] Scanning Kubernetes Events..."
EVENT_COUNT=$(kubectl get events -n cortexops-demo --field-selector type=Warning --since=1m -o json | grep -c '"reason": "Failed"' || echo "0")
if [ "$EVENT_COUNT" -gt 0 ]; then
    echo "✅ Kubernetes recorded the failure event."
else
    echo "⚠️  No warning events found in Kubernetes recently."
fi

echo "[4/4] Verifying End-to-End Pipeline Observability..."
FAILED=0

check_log() {
    local service=$1
    local search_string=$2
    local result=$($DOCKER_CMD logs --since="$START_TIME" compose-${service}-1 2>&1 | grep "$search_string" | head -n 1)
    
    if [[ -n "$result" ]]; then
        echo "✅ ${service} processing confirmed!"
        echo "   -> $result"
    else
        echo "❌ ${service} processing failed! Could not find '$search_string' in logs."
        FAILED=1
    fi
}

check_log "collector" "Validation: Collector published telemetry"
check_log "correlator" "Validation: Correlator processing telemetry"
check_log "rca" "Validation: RCA processing incident"
check_log "remediation" "Validation: Remediation proposing action from RCA"
check_log "remediation" "Validation: Temporal workflow created"

echo "========================================="
if [ $FAILED -eq 0 ]; then
    echo "🎉 VALIDATION PASSED: The entire CortexOps pipeline is operational."
    exit 0
else
    echo "💥 VALIDATION FAILED: Pipeline broke down. See logs above."
    exit 1
fi
