#!/bin/bash
set -e

echo "Starting CortexOps Runtime Verification..."

# Check HTTP endpoints
SERVICES=("topology:9091/debug/healthz" "nats:8222/varz" "prometheus:9090/-/healthy" "grafana:3000/api/health")

for svc in "${SERVICES[@]}"; do
  echo "Checking $svc..."
  if ! curl -s "http://$svc" > /dev/null; then
    echo "ERROR: Service $svc is unreachable"
    exit 1
  fi
done

# Check Temporal
echo "Checking Temporal status..."
if ! curl -s "http://temporal:8233" > /dev/null; then
  echo "ERROR: Temporal UI is unreachable"
  exit 1
fi

# Check Qdrant
echo "Checking Qdrant status..."
if ! curl -s "http://qdrant:6333/healthz" > /dev/null; then
  echo "ERROR: Qdrant is unreachable"
  exit 1
fi

echo "All services are up and healthy!"
