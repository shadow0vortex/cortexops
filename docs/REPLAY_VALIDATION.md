# Replay Safety Validation

CortexOps is designed for exactly-once delivery and idempotent remediation. This document outlines the validation procedure and success criteria.

---

## Design Guarantees

1. **Broker-Level Deduplication**: NATS JetStream uses `Nats-Msg-Id` to suppress duplicate messages at the broker layer before they reach consumers.
2. **Temporal Workflow Idempotency**: Each incident maps to a single Temporal workflow execution. Duplicate triggers are rejected by the workflow ID uniqueness constraint.
3. **OPA Policy Consistency**: The same policy input always produces the same allow/deny decision — Rego evaluation is purely functional.

---

## Validation Procedure

### Step 1: Broker Deduplication
1. Publish a burst of 10,000 identical K8s events to NATS with the same `Nats-Msg-Id`.
2. Verify that JetStream stores exactly 1 message.
3. The Correlator must receive exactly 1 event.

### Step 2: Incident Correlation
1. The Correlation Engine evaluates the telemetry.
2. Only one `CorrelatedIncident` should be generated for the identical burst.
3. No duplicate incidents should appear in the Topology API.

### Step 3: Workflow Idempotency
1. Simulate a duplicate event bypassing the broker (direct Temporal trigger).
2. Temporal must reject the duplicate workflow start (same workflow ID).
3. Only one remediation action is executed.

### Step 4: Crash Recovery
1. Kill the Correlator pod during event processing.
2. NATS JetStream redelivers unacknowledged messages after reconnection.
3. The Correlator deduplicates against its state store.
4. No duplicate incidents are created.

---

## Automated Validation
Run the full replay validation suite:
```bash
# Duplicate event storm
go run ./cmd/chaos -- --test duplicate-storm --count 10000

# Workflow idempotency
go run ./cmd/chaos -- --test workflow-idempotency

# Database partition + recovery
go run ./cmd/chaos -- --test db-partition
```

Or via Make target:
```bash
make chaos-test
```

---

## Success Criteria

- [x] NATS JetStream deduplicates identical messages (verified via `Nats-Msg-Id`).
- [x] Only one incident is visible in the Diagnostics API after a 10k duplicate storm.
- [x] No duplicate remediation actions are recorded in the Temporal execution history.
- [x] No duplicate alerts fire in Prometheus after replay events.
- [x] Crash recovery produces zero duplicate incidents after pod restart.
