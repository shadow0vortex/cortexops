# Soak Testing & Stability Validation

CortexOps must maintain stability under sustained telemetry pressure.

## Methodology
The soak test runs for 24 hours in a simulated cluster environment (Kind).
- **Ingestion Rate**: 500 events/second (simulated).
- **Topology Mutation**: 10 pod evictions/hour.
- **RCA Load**: Continuous generation of AI reports.

## Metrics to Monitor
- `go_goroutines`: Should remain stable (bounded).
- `process_resident_memory_bytes`: Should not show linear growth (no leaks).
- `cortexops_telemetry_published_total`: Total volume processed.
- `nats_stream_lag`: Backpressure monitoring.

## Success Criteria
- [ ] Resident memory remains below 1GiB per service.
- [ ] Goroutine count remains below 500 per service.
- [ ] Topology graph remains consistent with cluster state.
- [ ] Zero unhandled panics or crashes.
