# Runtime Operations

CortexOps provides deep visibility into the runtime state of your Kubernetes cluster through a fully integrated observability and operations stack.

---

## Telemetry Flow

1. **Collector** tails Kubernetes API Server events via `client-go` Informers.
2. **Normalizer** serializes events into strongly-typed Protobuf `TelemetryEnvelope` messages.
3. **NATS JetStream** buffers the telemetry stream with exactly-once delivery guarantees (via `Nats-Msg-Id` deduplication).
4. **Correlation Engine** groups events by topology proximity, temporal windowing, and TraceID into `CorrelatedIncident` objects.
5. **RCA Engine** enriches incidents with vector-similar historical context (Qdrant) and LLM-generated advisory reports.
6. **Remediation Engine** evaluates OPA policies and orchestrates Temporal workflows for safe Kubernetes mutation.

---

## Service Health Monitoring

All 5 CortexOps services expose health and metrics endpoints:

| Endpoint | Port | Purpose | Authentication |
|----------|------|---------|----------------|
| `/health` | 9091 | Service health (JSON) | None |
| `/debug/healthz` | 9091 | K8s liveness/readiness probe | None |
| `/metrics` | 9091 | Prometheus scraping | None |
| `/v1/topology/*` | 9091 | Topology data API | Bearer token |

### Kubernetes Probes (Helm Deployment)
Every service has both probes configured:
- **Liveness**: `httpGet /debug/healthz:9091` — restarts unresponsive pods.
- **Readiness**: `httpGet /debug/healthz:9091` — holds traffic from unready pods.
- Tuning: `initialDelaySeconds: 10`, `periodSeconds: 15`, `failureThreshold: 3`.

---

## Operational Dashboards

### Grafana (`http://localhost:3000`)
Pre-provisioned dashboards surface key platform SLIs:
- **Ingestion Rate**: Volume of K8s events processed per second by the Collector.
- **Incident Depth**: Average number of correlated events per incident.
- **Remediation Efficacy**: Ratio of successful remediations to rollbacks.
- **Service Latency**: P50/P95/P99 request duration across services.
- **Error Rate**: Per-service error rate with alerting thresholds.

Datasources are auto-provisioned:
- **Prometheus**: Metrics from all services via `/metrics`.
- **Loki**: Centralized logs shipped by Promtail.

### Prometheus Alerting (`deploy/prometheus/rules.yml`)
Active alert rules include:
| Alert | Condition | Severity |
|-------|-----------|----------|
| `CortexOpsServiceDown` | Service health endpoint unreachable for > 2m | Critical |
| `CortexOpsHighErrorRate` | Error rate > 5% over 5m window | Warning |
| `CortexOpsHighLatency` | P95 latency > 500ms over 5m window | Warning |
| `CortexOpsNATSBackpressure` | JetStream pending message count > 10k | Warning |

---

## Temporal Workflow Navigation

Every remediation is a durable Temporal workflow, accessible via the Temporal UI at `http://localhost:8233`.

### Workflow States
| State | Description |
|-------|-------------|
| `PROPOSED` | Action pending policy evaluation |
| `POLICY_EVALUATING` | OPA Rego rules being evaluated |
| `MAINTENANCE_BLOCKED` | Blocked by active maintenance window |
| `APPROVAL_PENDING` | Risk score > 0.5; awaiting human authorization |
| `REJECTED` | Blocked by OPA (namespace, action type, or maintenance) |
| `DRY_RUNNING` | Validating against K8s API (read-only) |
| `EXECUTING` | Kubernetes mutation in progress |
| `VERIFYING` | Waiting for telemetry stabilization |
| `SUCCESS` | Remediation completed, system stabilized |
| `ROLLING_BACK` | Stabilization failed, reversing mutation |
| `ROLLBACK_COMPLETE` | Rollback finished |

### OPA Policy Rules
Remediation actions are gated by the following policy layers:
1. **Action Allow-List**: Only `POD_RESTART`, `DEPLOYMENT_ROLLOUT_RESTART`, `HORIZONTAL_SCALE` are permitted.
2. **Namespace Protection**: `kube-system`, `kube-public`, `kube-node-lease`, `cortexops` are immutable.
3. **Maintenance Window**: All actions blocked when `input.maintenance_window == true`.
4. **Risk Threshold**: Actions with `risk_score > 0.5` require explicit human approval.

---

## Backup Operations

Automated backups run hourly via the `backup-cron` container:

| Target | Method | Retention |
|--------|--------|-----------|
| PostgreSQL | `pg_dump --format=custom` | 7 days |
| Qdrant | `POST /collections/*/snapshots` | 7 days |

Backups are stored in the `backupdata` Docker volume. Pruning of artifacts older than 7 days is automatic.

### Manual Backup
```bash
docker compose -f deploy/compose/docker-compose.dev.yaml exec backup-cron \
  sh /scripts/backup.sh
```

### Verify Backup
```bash
docker compose -f deploy/compose/docker-compose.dev.yaml exec backup-cron \
  ls -lah /backups/
```

---

## Resource Management

### Docker Compose (Development)
All runtime services are resource-bounded:
- **CPU**: 0.5 cores per service
- **Memory**: 512MB per service
- **Network**: Isolated `cortexops-dev` bridge network

### Kubernetes (Production)
Resource requests and limits are defined per-service in `values.yaml`:

| Service | CPU Request | CPU Limit | Memory Request | Memory Limit |
|---------|-------------|-----------|----------------|--------------|
| Collector | 100m | 200m | 128Mi | 256Mi |
| Correlator | 200m | 500m | 256Mi | 512Mi |
| Topology | 100m | 200m | 512Mi | 1Gi |
| RCA | 500m | 1000m | 1Gi | 2Gi |
| Remediation | 200m | 500m | 256Mi | 512Mi |

All services support HPA (Horizontal Pod Autoscaler) with CPU and memory targets at 80%.
