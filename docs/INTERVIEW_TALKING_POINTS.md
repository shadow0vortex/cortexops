# Interview Talking Points: CortexOps

This document provides high-signal talking points for technical interviews, focusing on the engineering rigor and distributed systems design behind CortexOps.

## 1. Why Temporal?
"I used **Temporal** to orchestrate remediation because infrastructure mutations are fundamentally unreliable. A simple script might fail mid-execution, leaving the system in an inconsistent state. Temporal provides **durable execution**, ensuring that even if the CortexOps platform crashes, the remediation workflow resumes exactly where it left off, with built-in support for idempotent retries and automatic rollbacks."

## 2. Determinism vs. AI "Magic"
"In CortexOps, I enforced a strict boundary: **AI is advisory-only**. While I used RAG (Retrieval-Augmented Generation) to enrich incidents with historical context, the actual decision to remediate is governed by **deterministic heuristics** and **OPA policies**. This prevents 'hallucinated' infrastructure changes and ensures that the system's behavior remains predictable and auditable."

## 3. Topology-Aware Correlation
"Standard monitoring tools often group events solely by time, leading to alert fatigue. CortexOps implements a **real-time topology graph**. This allows us to correlate symptoms based on their physical and logical relationships in the cluster (e.g., matching a pod failure to its host node's pressure events), resulting in high-confidence causal chains rather than just temporal clusters."

## 4. Replay Safety
"The system is designed with **NATS JetStream** as the backbone. By using sequence-based deduplication and event-timestamp-driven windowing, we ensure that replaying a historical telemetry stream results in the exact same incident generation. This is critical for post-mortems and validating system behavior during development."

## 5. Blast Radius Calculation
"Before any action is proposed, the platform performs a **BFS traversal** on the topology graph to calculate the blast radius of a failure. This depth-of-impact score is fed into our risk assessment engine, enabling the system to automatically escalate high-impact failures to human-in-the-loop approval gates while autonomously fixing low-risk issues."
