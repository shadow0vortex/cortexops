# Debugging Guide

Use this guide to troubleshoot the CortexOps platform in development and production environments.

## Prerequisites
- Docker Compose stack running (`docker compose --profile full up -d`)
- `curl` or equivalent HTTP client available

---

## Diagnostics API

The Topology service exposes a versioned REST API on port `9091`. Data endpoints require a Bearer token when `DIAG_API_TOKEN` is configured.

### Health Check (Unauthenticated)
Verify the service is running and responsive to Kubernetes probes:
```bash
curl http://localhost:9091/health
# Expected: {"status":"healthy"}

curl http://localhost:9091/debug/healthz
# Used by K8s liveness/readiness probes
```

### List Topology Nodes
Inspect all nodes in the dependency graph:
```bash
curl -H "Authorization: Bearer $DIAG_API_TOKEN" \
  http://localhost:9091/v1/topology/nodes
# Returns: {"count": N, "nodes": [...]}
```

### Calculate Blast Radius
Determine the downstream impact of a specific node failure:
```bash
curl -H "Authorization: Bearer $DIAG_API_TOKEN" \
  http://localhost:9091/v1/topology/blast-radius/pod-cortexops-demo-api-123
# Returns: {"node": "...", "count": N, "impacted": [...]}
```

---

## Observability Stack

### Grafana (Dashboards & Logs)
- **URL**: `http://localhost:3000`
- Pre-provisioned with Prometheus and Loki datasources.
- The CortexOps dashboard shows ingestion rate, incident depth, and remediation efficacy.

### Prometheus (Metrics)
- **URL**: `http://localhost:9090`
- All 5 services expose `/metrics` on port 9091.
- Alerting rules fire on: service down, high error rate, latency threshold breach.

```bash
# Verify targets are being scraped
curl http://localhost:9090/api/v1/targets | jq '.data.activeTargets[].health'
```

### Loki (Centralized Logs)
- **URL**: `http://localhost:3100`
- Promtail automatically ships container logs from the Docker daemon.
- Query logs in Grafana's Explore view using LogQL:

```
{container="topology"} |= "error"
{container="correlator"} | json | level="ERROR"
```

### NATS Monitoring
- **URL**: `http://localhost:8222`
- Check connection count and message throughput:
```bash
curl http://localhost:8222/varz | jq '{connections, in_msgs, out_msgs}'
curl http://localhost:8222/jsz | jq '.streams'
```

---

## Common Troubleshooting Playbooks

### No incidents appearing in Grafana
1. **Check Collector health**:
   ```bash
   docker compose -f deploy/compose/docker-compose.dev.yaml ps collector
   curl http://localhost:9091/health  # Topology health
   ```
2. **Verify NATS connectivity** (check authenticated connection):
   ```bash
   curl http://localhost:8222/connz | jq '.connections | length'
   ```
3. **Ensure K8s events are being generated**:
   ```bash
   kubectl get events --all-namespaces --sort-by='.lastTimestamp' | tail -20
   ```
4. **Check Correlator logs**:
   ```bash
   docker compose -f deploy/compose/docker-compose.dev.yaml logs correlator --tail=50
   ```

### Remediation workflows failing
1. **Open Temporal UI**: `http://localhost:8233` (credentials: `admin / password`)
2. Inspect the `History` tab of the failed workflow for activity errors.
3. **Check OPA policy evaluation**:
   - Verify the action type is in the allowed set (`POD_RESTART`, `DEPLOYMENT_ROLLOUT_RESTART`, `HORIZONTAL_SCALE`).
   - Verify the target namespace is not protected (`kube-system`, `kube-public`, `cortexops`).
   - Verify no maintenance window is active.
4. **Check Remediation service logs**:
   ```bash
   docker compose -f deploy/compose/docker-compose.dev.yaml logs remediation --tail=50
   ```

### Services failing to start (container exit)
1. **Check resource limits**: Services are capped at `0.5 CPU / 512M RAM`. Check if the host has sufficient resources:
   ```bash
   docker stats --no-stream
   ```
2. **Check non-root user permissions**: Services run as UID 10001. Verify the binary is executable:
   ```bash
   docker compose -f deploy/compose/docker-compose.dev.yaml exec topology id
   # Expected: uid=10001(cortexops)
   ```

### Database migration errors
1. **Check migration status**:
   ```bash
   docker compose -f deploy/compose/docker-compose.dev.yaml exec postgres \
     psql -U cortex -d cortexops -c "SELECT * FROM schema_migrations;"
   ```
2. **Verify per-service schemas exist**:
   ```bash
   docker compose -f deploy/compose/docker-compose.dev.yaml exec postgres \
     psql -U cortex -d cortexops -c "\dn"
   # Expected: topology, correlator, rca, remediation
   ```

### NATS authentication failures
1. **Verify credentials in compose**: Services must use `nats://admin:cortexpassword@nats:4222`.
2. **Check NATS logs**:
   ```bash
   docker compose -f deploy/compose/docker-compose.dev.yaml logs nats --tail=20
   ```
3. Look for `Authorization Violation` in the output.

---

## Chaos Testing

The `cmd/chaos` CLI provides automated failure injection:

```bash
# Duplicate event storm (tests NATS deduplication)
go run ./cmd/chaos -- --test duplicate-storm --count 1000

# Workflow idempotency validation
go run ./cmd/chaos -- --test workflow-idempotency

# Database partition simulation
go run ./cmd/chaos -- --test db-partition
```
