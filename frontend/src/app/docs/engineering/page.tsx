export default function EngineeringPage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Engineering Decisions</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        The design maturity and tradeoffs behind CortexOps.
      </p>

      <div className="space-y-8">
        
        <div className="bg-zinc-900/30 border border-zinc-800 p-8 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-4">Why Microservices?</h2>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Failure isolation.</strong> If the RCA Engine runs out of memory, the Correlation Engine still processes events and drops generic alerts.</li>
            <li><strong>Independent scaling.</strong> The Collector scales linearly with K8s cluster size, while the Topology Service scales by cluster complexity.</li>
            <li><strong>Clear service ownership.</strong> Forces well-defined API boundaries via Protobufs.</li>
          </ul>
        </div>

        <div className="bg-zinc-900/30 border border-zinc-800 p-8 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-4">Why NATS JetStream?</h2>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Replay safety.</strong> The Correlation Engine can crash, reboot, and replay from the exact sequence ID it left off.</li>
            <li><strong>Event durability.</strong> JetStream persists telemetry to disk synchronously, ensuring no data loss even during network partitions.</li>
            <li><strong>Decoupled services.</strong> Collectors don&apos;t know about Correlators. The broker acts as the ultimate shock absorber.</li>
            <li><strong>Guarantees:</strong> Exactly-once delivery semantics using message IDs.</li>
          </ul>
        </div>

        <div className="bg-zinc-900/30 border border-zinc-800 p-8 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-4">Why Temporal?</h2>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Durable execution.</strong> Workflows are effectively immortal. If a pod dies during a 10-minute wait for human approval, it resumes immediately upon scheduling.</li>
            <li><strong>Workflow replay.</strong> When code changes, Temporal handles history replay to ensure the workflow doesn&apos;t diverge unsafely.</li>
            <li><strong>State recovery.</strong> Built-in retry loops and compensation blocks replace thousands of lines of bespoke distributed systems code.</li>
          </ul>
        </div>

        <div className="bg-zinc-900/30 border border-zinc-800 p-8 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-4">Why Heuristics over ML for Correlation?</h2>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Deterministic behavior.</strong> In an operational control plane, you must be able to trace exactly *why* two events were correlated.</li>
            <li><strong>Auditability.</strong> A weighted sum based on TraceID matches, time-proximity, and topology depth is 100% explainable in an incident review.</li>
            <li><strong>Predictable outcomes.</strong> Probabilistic models hallucinate. Heuristics do not.</li>
          </ul>
        </div>

        <div className="bg-cortex-900/10 border border-cortex-500/30 p-8 rounded-2xl">
          <h2 className="text-2xl font-semibold text-white mb-4">Why Advisory AI?</h2>
          <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
            <li><strong>Safety.</strong> AI models are not deterministic enough to hold the keys to production infrastructure.</li>
            <li><strong>Governance.</strong> By decoupling the recommendation (LLM) from the action (Temporal/OPA), we enforce safety programmatically.</li>
            <li><strong>Human accountability.</strong> The AI acts as the ultimate Staff Engineer standing over your shoulder, but a human or deterministic policy clicks the button.</li>
          </ul>
        </div>

      </div>
    </>
  );
}
