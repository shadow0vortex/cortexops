# CortexOps Architecture Diagrams

These diagrams provide operational clarity into the distributed nature of CortexOps.

## 1. High-Level System Architecture

```mermaid
graph TD
    subgraph Kubernetes Cluster
        A[K8s API Server]
        B[Pods / Deployments]
    end

    subgraph CortexOps Control Plane
        C[Collector / Informers]
        D[(NATS JetStream)]
        E[Correlation Engine]
        F[Topology Graph]
        G[RCA Engine]
        H[(Qdrant Vector DB)]
        I[LLM Provider]
        J[Remediation Engine]
        K[(Temporal Orchestrator)]
        L[OPA Policy Engine]
    end

    A -->|Watch Events| C
    C -->|Normalize| D
    D -->|Consume| E
    A -->|Watch State| F
    E <-->|Query Blast Radius| F
    E -->|Correlated Incident| G
    G <-->|Fetch Historical| H
    G <-->|Generate Advisory| I
    G -->|RCA Report| J
    J <-->|Evaluate Rules| L
    J <-->|Durability| K
    K -->|Execute Patch| A
```

## 2. Remediation Lifecycle (Temporal Workflow)

```mermaid
stateDiagram-v2
    [*] --> PROPOSED
    PROPOSED --> POLICY_EVALUATING : Trigger Temporal Workflow
    POLICY_EVALUATING --> APPROVAL_PENDING : Risk Score > 0.7
    POLICY_EVALUATING --> REJECTED : Blocked by OPA (e.g. kube-system)
    
    APPROVAL_PENDING --> APPROVED : Slack Authorized
    APPROVAL_PENDING --> REJECTED : Timeout / Denied
    
    APPROVED --> DRY_RUNNING
    DRY_RUNNING --> EXECUTING : Validation Passed
    DRY_RUNNING --> FAILED : RBAC / Syntax Error
    
    EXECUTING --> VERIFYING
    VERIFYING --> SUCCESS : Telemetry Stabilized
    VERIFYING --> ROLLING_BACK : Telemetry Fails
    
    ROLLING_BACK --> ROLLBACK_COMPLETE
    
    SUCCESS --> [*]
    ROLLBACK_COMPLETE --> [*]
    REJECTED --> [*]
    FAILED --> [*]
```

## 3. Replay & Recovery Flow (Chaos Degradation)

```mermaid
sequenceDiagram
    participant NATS as JetStream
    participant CE as Correlation Engine
    participant DB as Audit DB

    Note over NATS,DB: Standard Execution
    NATS->>CE: Event A (id: 123)
    CE->>DB: Save Incident (id: 123)

    Note over NATS,DB: Chaos Partition Occurs
    NATS--xCE: Connection Dropped
    CE--xCE: Pod Crashes

    Note over NATS,DB: Recovery & Replay Phase
    CE->>NATS: Reconnect & Re-subscribe
    NATS->>CE: Event A (id: 123) [Redelivered]
    CE->>DB: Query Exists? (id: 123)
    DB-->>CE: True
    CE->>CE: Drop Duplicate Event deterministically
```
