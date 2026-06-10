# Runtime Operations

CortexOps provides deep visibility into the runtime state of your Kubernetes cluster.

## Telemetry Flow
1. **Collector** tails K8s Informers.
2. **Normalizer** flattens events into Protobuf.
3. **NATS** buffers the telemetry stream.
4. **Correlation Engine** groups events by topology and time.

## Operational Dashboards
Our Grafana dashboards are pre-configured to surface key platform SLIs:
- **Ingestion Rate**: Volume of K8s events being processed.
- **Incident Depth**: Average number of events per correlated incident.
- **Remediation Efficacy**: Ratio of successful remediations to rollbacks.

## Temporal Workflow Navigation
Every remediation is a durable Temporal workflow.
- **PROPOSED**: Action is pending policy check.
- **DRY_RUNNING**: Validating against K8s API (Read-Only).
- **EXECUTING**: Mutation in progress.
- **VERIFYING**: Waiting for system stabilization.
- **SUCCESS/ROLLBACK**: Final state.

You can inspect the full execution history and activity heartbeats in the Temporal UI at `http://localhost:8233`.
