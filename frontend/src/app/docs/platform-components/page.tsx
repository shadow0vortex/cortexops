export default function PlatformComponentsPage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Platform Components</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        The independent microservices that make up the CortexOps control plane.
      </p>

      <div className="space-y-8">
        
        {/* Collector Service */}
        <div className="bg-zinc-900/30 border border-zinc-800 rounded-2xl p-8">
          <h2 className="text-2xl font-semibold text-white mb-2">Collector Service</h2>
          <p className="text-cortex-400 text-sm font-mono mb-6">Purpose: Telemetry Ingestion</p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
            <div>
              <span className="text-zinc-500 block mb-1">Inputs</span>
              <p className="text-zinc-300">Kubernetes API Watch Streams, Prometheus metrics.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Outputs</span>
              <p className="text-zinc-300">Normalized protobuf events to NATS JetStream.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Dependencies</span>
              <p className="text-zinc-300">Kubernetes API, NATS JetStream.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Failure Modes</span>
              <p className="text-zinc-300">If K8s API is unreachable, local caching handles retries. If NATS is down, backpressure is applied to prevent OOM.</p>
            </div>
          </div>
        </div>

        {/* Correlation Engine */}
        <div className="bg-zinc-900/30 border border-zinc-800 rounded-2xl p-8">
          <h2 className="text-2xl font-semibold text-white mb-2">Correlation Engine</h2>
          <p className="text-cortex-400 text-sm font-mono mb-6">Purpose: Incident Grouping</p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
            <div>
              <span className="text-zinc-500 block mb-1">Inputs</span>
              <p className="text-zinc-300">Event streams from NATS JetStream.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Outputs</span>
              <p className="text-zinc-300">Correlated incident objects to RCA Service.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Dependencies</span>
              <p className="text-zinc-300">NATS, Topology Service, PostgreSQL.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Failure Modes</span>
              <p className="text-zinc-300">Uses deterministic heuristic scoring. If Topology Service is slow, correlation degrades gracefully using cached relationships.</p>
            </div>
          </div>
        </div>

        {/* Topology Service */}
        <div className="bg-zinc-900/30 border border-zinc-800 rounded-2xl p-8">
          <h2 className="text-2xl font-semibold text-white mb-2">Topology Service</h2>
          <p className="text-cortex-400 text-sm font-mono mb-6">Purpose: Blast Radius Analysis</p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
            <div>
              <span className="text-zinc-500 block mb-1">Inputs</span>
              <p className="text-zinc-300">State change events from K8s.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Outputs</span>
              <p className="text-zinc-300">Graph queries (ancestors, descendants).</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Dependencies</span>
              <p className="text-zinc-300">PostgreSQL (for async persistence).</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Failure Modes</span>
              <p className="text-zinc-300">Runs in-memory. If the pod crashes, it restores the graph from Postgres snapshots upon restart.</p>
            </div>
          </div>
        </div>

        {/* RCA Service */}
        <div className="bg-zinc-900/30 border border-zinc-800 rounded-2xl p-8">
          <h2 className="text-2xl font-semibold text-white mb-2">RCA Service</h2>
          <p className="text-cortex-400 text-sm font-mono mb-6">Purpose: Advisory Root Cause Generation</p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
            <div>
              <span className="text-zinc-500 block mb-1">Inputs</span>
              <p className="text-zinc-300">Incidents, Qdrant vectors (historical docs).</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Outputs</span>
              <p className="text-zinc-300">Advisory RCAReport object.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Dependencies</span>
              <p className="text-zinc-300">Qdrant, External LLM Provider.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Failure Modes</span>
              <p className="text-zinc-300">If the LLM is unavailable, it fails closed to deterministic rules. The AI is advisory only and never mutates state.</p>
            </div>
          </div>
        </div>

        {/* Remediation Service */}
        <div className="bg-zinc-900/30 border border-zinc-800 rounded-2xl p-8">
          <h2 className="text-2xl font-semibold text-white mb-2">Remediation Service</h2>
          <p className="text-cortex-400 text-sm font-mono mb-6">Purpose: Workflow Execution & Governance</p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
            <div>
              <span className="text-zinc-500 block mb-1">Inputs</span>
              <p className="text-zinc-300">RCAReport, Human Approval.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Outputs</span>
              <p className="text-zinc-300">Infrastructure mutations via K8s API.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Dependencies</span>
              <p className="text-zinc-300">Temporal, OPA, Kubernetes API.</p>
            </div>
            <div>
              <span className="text-zinc-500 block mb-1">Failure Modes</span>
              <p className="text-zinc-300">If OPA denies the action, execution is aborted. If mutation fails, Temporal triggers the compensation (rollback) transaction.</p>
            </div>
          </div>
        </div>

      </div>
    </>
  );
}
