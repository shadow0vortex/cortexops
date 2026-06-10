export default function InfrastructurePage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Infrastructure</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        The core distributed systems that back the CortexOps platform.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        
        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-3">NATS JetStream</h2>
          <p className="text-sm text-zinc-400 mb-4">
            Acts as the distributed event log for all incoming telemetry. 
          </p>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Why?</strong> Provides replayable, ordered streams. Essential for crash recovery without data loss.</li>
            <li><strong>Guarantees:</strong> Exactly-once delivery semantics using message IDs.</li>
          </ul>
        </div>

        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-3">Temporal Workflows</h2>
          <p className="text-sm text-zinc-400 mb-4">
            The orchestration engine responsible for executing remediations.
          </p>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Why?</strong> Provides durable execution. Workflows can sleep for days awaiting human approval or safely retry upon network partition.</li>
            <li><strong>Guarantees:</strong> Replay safety. If a worker dies mid-execution, a new worker resumes exact state without duplicating prior activities.</li>
          </ul>
        </div>

        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-3">PostgreSQL</h2>
          <p className="text-sm text-zinc-400 mb-4">
            The persistent storage layer for the Topology Graph, Audit logs, and Temporal state.
          </p>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Why?</strong> Relational guarantees are required for the audit trail of autonomous actions.</li>
            <li><strong>Failure Mode:</strong> If PG goes down, CortexOps halts autonomous execution safely (fails closed).</li>
          </ul>
        </div>

        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-3">Qdrant Memory Layer</h2>
          <p className="text-sm text-zinc-400 mb-4">
            A vector database storing historical incidents, runbooks, and past remediation outcomes.
          </p>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Why?</strong> Allows the RCA engine to perform semantic search to find similar past outages.</li>
            <li><strong>Failure Mode:</strong> If unavailable, the RCA engine degrades to deterministic heuristics instead of Retrieval-Augmented Generation.</li>
          </ul>
        </div>

        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-2xl md:col-span-2">
          <h2 className="text-2xl font-semibold text-white mb-3">OPA Governance</h2>
          <p className="text-sm text-zinc-400 mb-4">
            The Open Policy Agent (OPA) enforces security constraints before any infrastructure mutation occurs.
          </p>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Why?</strong> Decouples security logic from workflow code. Ensures that even if the AI or Correlator hallucinations an incorrect fix, OPA will block it if it violates policy (e.g., &quot;Never restart kube-system pods&quot;).</li>
            <li><strong>Guarantees:</strong> Fail-closed execution. No action is taken unless explicitly allowed by a loaded rego policy.</li>
          </ul>
        </div>

      </div>
    </>
  );
}
