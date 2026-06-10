# Kubernetes Deployment Validation

This document verifies the production readiness of the CortexOps Helm chart.

## Validated Features
1.  **Pod Anti-Affinity**: Ensures services are distributed across nodes to survive node failure.
2.  **PDB (Pod Disruption Budget)**: Guarantees availability during cluster maintenance.
3.  **ServiceMonitor**: Automated discovery of metrics by Prometheus.
4.  **RBAC Verification**: Services run with the least-privilege necessary to watch and patch cluster resources.
5.  **HPA (Horizontal Pod Autoscaler)**: Validated for the Collector and Correlator services.

## Validation Commands
```bash
# Verify anti-affinity
kubectl get pods -l app=cortexops -o wide

# Verify PDB
kubectl get pdb -n cortexops-system

# Verify RBAC
kubectl auth can-i patch deployments -n cortexops-demo --as=system:serviceaccount:cortexops-system:cortexops-sa
```

## Upgrade/Rollback Safety
- **Rolling Update**: `maxUnavailable: 1` ensures service continuity.
- **Rollback**: Validated via `helm rollback`.
