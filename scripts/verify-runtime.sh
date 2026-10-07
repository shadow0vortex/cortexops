#!/bin/sh
set -eu

check_http() {
  name=$1
  url=$2
  attempt=1

  printf 'Checking %s...\n' "$name"
  while [ "$attempt" -le 30 ]; do
    if curl --fail --silent --show-error --connect-timeout 3 --max-time 5 "$url" >/dev/null; then
      printf '%s is healthy.\n' "$name"
      return 0
    fi
    attempt=$((attempt + 1))
    sleep 2
  done

  printf 'ERROR: %s did not return a successful HTTP response: %s\n' "$name" "$url" >&2
  return 1
}

echo "Starting CortexOps Runtime Verification..."

check_http collector http://collector:9091/debug/healthz
check_http correlator http://correlator:9091/debug/healthz
check_http topology http://topology:9091/debug/healthz
check_http rca http://rca:9091/debug/healthz
check_http remediation http://remediation:9091/debug/healthz
check_http nats http://nats:8222/varz
check_http prometheus http://prometheus:9090/-/healthy
check_http grafana http://grafana:3000/api/health
check_http temporal-ui http://temporal-ui:8080/
check_http qdrant http://qdrant:6333/healthz

echo "All HTTP runtime checks passed. PostgreSQL and Temporal server health are covered by their Compose healthchecks."
