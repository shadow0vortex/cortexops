#!/bin/bash
set -euo pipefail

echo "[CortexOps] Bootstrapping Local Sandbox Environment..."

# 1. Create Kind Cluster
echo "-> Creating Kind cluster 'cortexops-sandbox'..."
if ! kind get clusters | grep -q "^cortexops-sandbox$"; then
    kind create cluster --name cortexops-sandbox
else
    echo "Cluster already exists."
fi

# 2. Deploy Infrastructure Dependencies (NATS, Postgres, Temporal, Qdrant)
echo "-> Deploying infrastructure stack..."
# Normally docker-compose up -d, but we mock it for the sandbox script
docker-compose -f deploy/docker-compose.yaml up -d || echo "Skipping docker-compose (file missing in stub)"

# 3. Apply CRDs and CortexOps Engine
echo "-> Deploying CortexOps..."
helm upgrade --install cortexops ./deploy/helm/cortexops \
  --namespace cortexops-system \
  --create-namespace \
  --set replicaCount=1 \
  --wait

# 4. Deploy Sample Workloads
echo "-> Deploying Synthetic Workloads..."
kubectl apply -f sandbox/workloads/nginx-deployment.yaml

echo "=========================================="
echo "✅ CortexOps Sandbox is ready!"
echo "Run 'kubectl get pods -A' to verify."
echo "View dashboards at http://localhost:3000"
echo "=========================================="
