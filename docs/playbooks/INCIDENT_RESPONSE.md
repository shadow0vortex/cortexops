# Playbook: Incident Response

This guide provides steps for responding to platform-level incidents within CortexOps.

## 1. High Ingestion Lag
**Symptom**: `nats_stream_lag` is increasing; incidents appear late in Grafana.
**Action**:
1. Check Correlator logs: `docker compose logs correlator`.
2. Verify NATS health: `curl http://localhost:8222/varz`.
3. Scale the correlator service if CPU/Memory limits are breached.

## 2. Failed AI RCA
**Symptom**: RCA reports show `IsDegraded=true`.
**Action**:
1. Check Qdrant availability: `curl http://localhost:6333/healthz`.
2. Verify LLM API connectivity.
3. The system will continue to remediate using raw telemetry; no immediate action is required unless the degraded mode persists.

## 3. Stalled Remediations
**Symptom**: Workflows in Temporal UI are stuck in `RUNNING` for >10 minutes.
**Action**:
1. Check Remediation Worker logs: `docker compose logs remediation`.
2. Verify the task queue is not backed up in the Diagnostics API.
3. If the worker crashed, simply restart the container; Temporal will resume the workflow.
