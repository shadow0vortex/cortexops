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

    subgraph Observability Stack
        M[Prometheus]
        N[Grafana]
        O[Loki]
        P[Promtail]
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

    C -.->|/metrics| M
    E -.->|/metrics| M
    F -.->|/metrics| M
    G -.->|/metrics| M
    J -.->|/metrics| M
    M -->|Datasource| N
    O -->|Datasource| N
    P -->|Ship Logs| O
```

## 2. Remediation Lifecycle (Temporal Workflow)

```mermaid
stateDiagram-v2
    [*] --> PROPOSED
    PROPOSED --> POLICY_EVALUATING : Trigger Temporal Workflow
    POLICY_EVALUATING --> MAINTENANCE_BLOCKED : Maintenance Window Active
    POLICY_EVALUATING --> APPROVAL_PENDING : Risk Score > 0.5
    POLICY_EVALUATING --> REJECTED : Blocked by OPA
    
    MAINTENANCE_BLOCKED --> [*] : Retry After Window

    APPROVAL_PENDING --> APPROVED : Authorized
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
    participant DB as PostgreSQL

    Note over NATS,DB: Standard Execution
    NATS->>CE: Event A (id: 123)
    CE->>DB: Save Incident (id: 123)

    Note over NATS,DB: Network Partition / Pod Crash
    NATS--xCE: Connection Dropped
    CE--xCE: Pod Restarts

    Note over NATS,DB: Recovery & Replay Phase
    CE->>NATS: Reconnect (Authenticated)
    NATS->>CE: Event A (id: 123) [Redelivered]
    CE->>DB: Query Exists? (id: 123)
    DB-->>CE: True
    CE->>CE: Drop Duplicate Event
```

## 4. Network Security Architecture

```mermaid
graph LR
    subgraph cortexops-namespace
        subgraph Per-Service NetworkPolicy
            C[Collector SA]
            CO[Correlator SA]
            T[Topology SA]
            R[RCA SA]
            RE[Remediation SA]
        end
    end

    subgraph infrastructure
        NATS[NATS :4222]
        PG[PostgreSQL :5432]
        TMP[Temporal :7233]
        QD[Qdrant :6333]
    end

    subgraph monitoring
        PROM[Prometheus :9091]
    end

    C -->|publish| NATS
    CO -->|consume/publish| NATS
    CO -->|query| T
    R -->|consume| NATS
    R -->|search| QD
    RE -->|consume| NATS
    RE -->|workflows| TMP
    T -->|persist| PG

    PROM -.->|scrape /metrics| C
    PROM -.->|scrape /metrics| CO
    PROM -.->|scrape /metrics| T
    PROM -.->|scrape /metrics| R
    PROM -.->|scrape /metrics| RE
```

## 5. Database Schema Isolation

```mermaid
graph TD
    subgraph PostgreSQL cortexops
        PUB[public schema]
        TS[topology schema]
        CS[correlator schema]
        RS[rca schema]
        RMS[remediation schema]
    end

    TS -->|topology_snapshots| T[Topology Service]
    CS -->|reserved| CO[Correlator]
    RS -->|reserved| R[RCA Service]
    RMS -->|reserved| RE[Remediation]
```

## 6. Backup & Recovery Flow

```mermaid
sequenceDiagram
    participant Cron as backup-cron
    participant PG as PostgreSQL
    participant QD as Qdrant
    participant Vol as /backups Volume

    loop Every 1 Hour
        Cron->>PG: pg_dump --format=custom
        PG-->>Cron: cortexops.dump
        Cron->>Vol: Save timestamped backup
        Cron->>QD: POST /collections/*/snapshots
        QD-->>Cron: Snapshot ID
        Cron->>Vol: Prune backups > 7 days
    end
```
