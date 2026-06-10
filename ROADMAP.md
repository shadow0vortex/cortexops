# CortexOps Roadmap

This document outlines the strategic direction for CortexOps as it moves toward a 1.0 release.

## Q3 2026: Multi-Cluster Federation
- **Goal:** Manage fleet-wide incidents from a single management plane.
- **Architecture:** Utilize NATS Leaf Nodes to bridge telemetry from edge clusters to the central management plane. All local remediations will execute on the edge.

## Q4 2026: Service Mesh Integration
- **Goal:** Advanced network-level causal analysis.
- **Architecture:** Ingest Envoy/Istio sidecar metrics directly into the Temporal Correlation buckets to determine exact latency hop degradation.

## Q1 2027: Predictive Queue Scaling
- **Goal:** Prevent Temporal workflow starvation automatically.
- **Architecture:** Dynamically autoscale Temporal worker fleets based on NATS Consumer Lag metrics.

*(Note: We will not introduce autonomous policy generation or self-modifying LLM agents. All policy must remain declaratively defined by human operators via GitOps.)*
