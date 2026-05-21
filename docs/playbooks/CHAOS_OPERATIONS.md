# Playbook: Chaos Operations

Guidelines for validating system resilience via intentional failure injection.

## 1. Validating Replay Safety
Execute a telemetry burst storm:
```bash
make chaos-test
```
Verify that only one incident is created per storm in the Grafana dashboard.

## 2. Simulating NATS Outage
1. Pause the NATS container: `docker pause compose-nats-1`.
2. Generate cluster events.
3. Unpause NATS: `docker unpause compose-nats-1`.
4. Verify that the Collector resumes and pushes all buffered events successfully.

## 3. Simulating Remediation Rollback
1. Inject a failure with an inherently unfixable state (e.g., non-existent image).
2. Wait for the `VERIFYING` state in Temporal UI.
3. Confirm that after 5 minutes of instability, the workflow executes the `Rollback` activity.
