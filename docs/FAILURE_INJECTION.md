# Failure Injection Framework

CortexOps includes a deterministic failure injection framework to validate the control plane's responsiveness and safety.

## Supported Scenarios

### 1. `rollout-fail`
Patches the `demo-frontend` deployment with a non-existent image tag.
- **Observed Behavior**: `ImagePullBackOff` events.
- **Correlation**: Grouped by the Frontend's topology nodes.
- **Remediation**: Triggers a rollout restart or rollback depending on policy.

### 2. `crashloop`
Injects a failing command into the `demo-api` container.
- **Observed Behavior**: `BackOff` events and rapid container restarts.
- **Correlation**: Causal chain linking the API failure to potential frontend 5xx errors.

### 3. `scaling`
Scales the frontend to 10 replicas instantly.
- **Observed Behavior**: Resource pressure and pod scheduling events.
- **CortexOps Logic**: Evaluates blast radius to see if the nodes can handle the burst.

## Implementation Details
The failure injection is orchestrated via a Go CLI in `cmd/demo/main.go` which uses `client-go` for atomic patches to the cluster state.
