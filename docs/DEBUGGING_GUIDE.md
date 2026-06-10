# Debugging Guide

Use this guide to troubleshoot the CortexOps platform or the demo environment.

## Diagnostics API
The platform exposes a read-only introspection API on port `9091`.

### Inspecting Topology
Find out what CortexOps knows about a specific node:
```bash
curl "http://localhost:9091/debug/graph/node?id=pod/cortexops-demo/demo-api-123"
```

Calculate the blast radius of a potential failure:
```bash
curl "http://localhost:9091/debug/graph/blast-radius?id=node/kind-control-plane"
```

### Inspecting Incidents
List all currently active correlation windows:
```bash
curl "http://localhost:9091/debug/incidents/active"
```

## Common Playbooks

### Issue: No incidents appearing in Grafana
1. Check if the Collector is running: `ps aux | grep collector`
2. Verify NATS connectivity: `curl http://localhost:8222/varz`
3. Ensure K8s events are being generated: `kubectl get events -n cortexops-demo`

### Issue: Remediation workflows failing
1. Open Temporal UI at `http://localhost:8233`
2. Inspect the `History` tab of the failed workflow.
3. Check the `TaskQueue` status in the Diagnostics API: `curl http://localhost:9091/debug/temporal/status`
