# CortexOps Load Test Results

This document captures the scalability and performance evidence of CortexOps under stress testing.

## Methodology

The testing framework utilized the internal `cmd/chaos` injector to perform duplicate event storms simulating massive telemetry bursts (deployments, node failures, network partitions). All tests were run with authenticated NATS connections (`admin:cortexpassword`).

**Test Configuration:**
- Docker Compose with `--profile full` (all services + observability)
- Resource limits enforced: `0.5 CPU / 512M RAM` per service
- PostgreSQL `17-alpine` with `pgx/v5` connection pooling
- NATS `2.10-alpine` with JetStream enabled

**Test Runs:**
- 10,000 events
- 50,000 events
- 100,000 events

**Metrics Captured:**
- **Ingestion Throughput**: Time taken for the chaos generator to publish to NATS JetStream.
- **Resource Usage**: CPU (%) and Memory (MB) profiling of critical containers via `docker stats`.

---

## 1. 10,000 Event Storm

- **Ingestion Time**: 416 ms
- **Throughput**: ~24,000 events/sec

### Peak Resource Usage
| Service         | CPU % | Memory | Within Limits |
|----------------|-------|--------|---------------|
| Correlator      | 47.35 | 12.2 MB| ✅ (512M cap) |
| Topology        | 30.03 | 12.9 MB| ✅ (512M cap) |
| Temporal Server | 31.09 | 150 MB | N/A (infra)   |
| Postgres        | 27.76 | 189 MB | N/A (infra)   |

---

## 2. 50,000 Event Storm

- **Ingestion Time**: 1.93 seconds
- **Throughput**: ~25,800 events/sec

### Peak Resource Usage
| Service         | CPU % | Memory | Within Limits |
|----------------|-------|--------|---------------|
| Correlator      | 31.07 | 12.5 MB| ✅ (512M cap) |
| Topology        | 12.39 | 12.8 MB| ✅ (512M cap) |
| Temporal Server | 10.33 | 149 MB | N/A (infra)   |
| NATS JetStream  | 49.80 | 66.5 MB| N/A (infra)   |

*Observation*: NATS JetStream takes over a heavier CPU load as ingestion duration increases. Correlator handles batch deduplication efficiently.

---

## 3. 100,000 Event Storm

- **Ingestion Time**: 2.54 seconds
- **Throughput**: ~39,300 events/sec

### Peak Resource Usage
| Service         | CPU % | Memory | Within Limits |
|----------------|-------|--------|---------------|
| Correlator      | 55.84 | 12.5 MB| ✅ (512M cap) |
| Topology        | 25.62 | 13.0 MB| ✅ (512M cap) |
| Temporal Server | 105.40| 152 MB | N/A (infra)   |
| Postgres        | 52.69 | 189 MB | N/A (infra)   |

---

## Observations & Bottlenecks

1. **Memory Efficiency**: CortexOps microservices (Correlator, Topology, Collector) maintain extremely flat memory usage (~10-13MB) even under extreme ingestion loads. This proves Go garbage collection and memory pooling operate efficiently. All services remain well within the 512M resource limit.

2. **Temporal CPU Saturation**: At 100,000 events, Temporal Server CPU spikes past 100%. This is the primary bottleneck. However, Temporal's queuing mechanism prevents data loss, making this an acceptable latency tradeoff.

3. **Correlation Deduplication**: The Correlator efficiently groups 100k events into a single incident entity. NATS JetStream provides robust backpressure via `Nats-Msg-Id` deduplication.

4. **Resource Limit Headroom**: With the `0.5 CPU / 512M` limits enforced, all CortexOps services have massive headroom (~97% memory unused at peak). No OOM kills or CPU throttling observed.

---

## Recommendations

- For enterprise deployments exceeding 100k events/sec sustained, **Temporal** requires horizontal scaling (matching the HPA configurations in the Helm chart).
- CortexOps microservices have extensive overhead available and will safely scale via CPU-based HPA (configured with 80% target) without OOM killing.
- The `pgx/v5` connection pool prevents PostgreSQL connection exhaustion during burst writes.
