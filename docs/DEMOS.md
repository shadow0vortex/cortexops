# CortexOps Demo Scenarios

This document contains step-by-step reproducible operational demos.

## 1. Failed Deployment Rollout (End-to-End RCA & Rollback)
**Goal:** Prove CortexOps can correlate symptoms, consult AI, attempt a remediation, and execute a deterministic rollback when verification fails.

### Steps:
1. Apply the broken backend from the sandbox:
   ```bash
   kubectl apply -f sandbox/workloads/nginx-deployment.yaml
   ```
2. **Observe Telemetry:** The `collector` watches `FailedCreate` and `ImagePullBackOff` events stream into NATS.
3. **Correlation Engine:** Aggregates these events into a `CorrelatedIncident` tied to the `demo-broken-backend` ReplicaSet.
4. **RCA Generation:** Qdrant returns past image pull errors. The LLM grounds the context, stating "Image tag nonexistent-tag cannot be resolved."
5. **Remediation & Rollback:** 
   - Policy allows a `DEPLOYMENT_ROLLOUT_RESTART`. 
   - Temporal executes the K8s patch.
   - Verification waits 5 minutes. The pods are still failing (because the image tag is inherently broken).
   - Temporal triggers `action.Rollback()`. State reverts safely.

## 2. Replay Recovery Validation (Idempotency)
**Goal:** Prove NATS JetStream event replay doesn't spawn duplicate executions.

### Steps:
1. Run a normal incident sequence (e.g., Delete a pod manually `kubectl delete pod <name>`).
2. Verify the incident was generated exactly once.
3. Use the NATS CLI to forcefully rewind the consumer sequence to replay the last 1000 events:
   ```bash
   nats consumer update cortexops-group --deliver all
   ```
4. **Observe:** The Correlation Engine receives the events again. Because events are hashed by `event_id`, they are deterministically dropped. 
5. **Verification:** Query the `AuditStore`. Exactly `0` new remediations were executed during the replay storm.

## 3. AI Degraded Mode
**Goal:** Prove the platform fails-closed safely if third-party LLM providers go down.

### Steps:
1. Inject a network partition to the AI provider using the Chaos framework:
   ```bash
   go test ./test/e2e/suites -run TestDegradedAIBehavior
   ```
2. Delete a pod to trigger an incident.
3. **Observe:** The RAG pipeline hits a context timeout (3 seconds). 
4. **Verification:** The `RCAReport` is published with `IsDegraded = true`. The Slack approval request contains raw telemetry instead of an AI summary. Remediation continues seamlessly.
