# CortexOps Engineering Showcase

This document provides SREs and Principal Engineers with structured talking points and STAR-format stories for architecture reviews or interviews.

## 1. Systems Engineering Highlight: Idempotent Workflows
**The Problem:** Distributed systems fail unpredictably. If CortexOps orchestrates a Kubernetes patch but the pod crashes before saving the state, recovering naively would apply the patch twice.
**The Solution:** We transitioned the remediation state machine from an in-memory Go routine to **Temporal**.
**The Impact:** Temporal provides durable execution. Activities (like mutating K8s) are bound by exponential retries. Because K8s is declarative, repeating a patch during a retry is mathematically idempotent. Operations became 100% replay-safe.

## 2. AI Governance: The Advisory-Only Boundary
**The Problem:** LLMs hallucinate. Allowing an LLM to directly issue `kubectl` commands against production clusters is an unacceptable operational risk.
**The Solution:** CortexOps enforces an architectural firewall. The AI pipeline is read-only, producing an `RCAReport`. Remediation is triggered exclusively by deterministic OPA policies mapping cluster symptoms to predefined execution profiles (`POD_RESTART`).
**The Impact:** Achieved human-in-the-loop autonomous operations without risking catastrophic cluster deletion.

## 3. Dealing with Event Storms (Backpressure)
**The Problem:** During a node failure, thousands of cascading pod events flood the system, risking OOM-kills on the correlation engine.
**The Solution:** We implemented **NATS JetStream** as a shock absorber. The ingestion pipeline pushes events to disk-backed streams.
**The Impact:** If correlation slows down, consumer lag increases, but zero events are dropped. When the storm subsides, the engine drains the queue deterministically.

## 4. Discussion Topic: Why Not ML-based Correlation?
Instead of using opaque ML models to cluster alerts, we built a real-time topology graph using `client-go` informers. ML models drift and require retraining. Our graph mathematically guarantees that a failing Pod belongs to a specific ReplicaSet, creating deterministic, explainable blast-radii.
