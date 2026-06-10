# Replay Safety Validation

CortexOps is designed for exactly-once delivery and idempotent remediation. This document outlines the validation procedure for replay safety.

## Validation Procedure
1.  **Ingestion**: A burst of identical K8s events is published to NATS with the same `Nats-Msg-Id`.
2.  **Deduplication**: The NATS JetStream broker must suppress duplicates based on the `Nats-Msg-Id`.
3.  **Incident Correlation**: The Correlation Engine evaluates the telemetry. Only one `CorrelatedIncident` should be generated for the identical burst.
4.  **Workflow Idempotency**: If a duplicate event bypasses the broker (simulated), the Temporal workflow must be idempotent, ensuring only one remediation action is executed.

## Automated Validation
Run the following command to execute the automated replay burst:
```bash
make chaos-test
```

## Success Criteria
- [ ] Only one incident is visible in the Diagnostics API.
- [ ] No duplicate remediation actions are recorded in the `AuditStore`.
- [ ] No duplicate alerts fire in Prometheus/Grafana.
