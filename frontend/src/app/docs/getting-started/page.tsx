import { Terminal } from "lucide-react";

export default function GettingStartedPage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Getting Started</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        Deploy CortexOps locally and explore a deterministic incident orchestration platform.
      </p>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Overview</h2>
        <p className="text-zinc-300 mb-4">
          CortexOps is an operational intelligence and orchestration platform. It is not &quot;autonomous AGI&quot; or &quot;self-driving infrastructure.&quot; It is a highly deterministic, governance-first control plane for automating incident remediation safely.
        </p>
        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl mb-6">
          <h3 className="text-lg font-medium text-white mb-3">What problem does this solve?</h3>
          <p className="text-sm text-zinc-400">
            Modern distributed systems emit massive amounts of telemetry during an incident. SREs face &quot;alert fatigue&quot; and cognitive overload when trying to correlate failures across microservices. CortexOps solves this by ingesting telemetry, calculating blast radius, identifying root causes, and executing pre-approved, replay-safe remediation workflows.
          </p>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Quick Start</h2>
        <p className="text-zinc-300 mb-6">
          Deploy the core control plane. This initializes NATS JetStream, PostgreSQL, Temporal, Prometheus, Grafana, and Qdrant.
        </p>
        <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300 mb-4">
          <div className="flex items-center gap-3">
            <Terminal className="w-4 h-4 text-zinc-500" />
            <span>make dev-up</span>
          </div>
        </div>
        <p className="text-zinc-300 mb-6 mt-8">
          Verify that the runtime is healthy and all distributed components are communicating.
        </p>
        <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300 mb-4">
          <div className="flex items-center gap-3">
            <Terminal className="w-4 h-4 text-zinc-500" />
            <span>make verify-runtime</span>
          </div>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Demo Environment</h2>
        <p className="text-zinc-300 mb-6">
          Bootstrap a sample microservices topology (Frontend, API, Worker, Redis, PostgreSQL) to see correlation in action.
        </p>
        <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300 mb-4">
          <div className="flex items-center gap-3">
            <Terminal className="w-4 h-4 text-zinc-500" />
            <span>make bootstrap</span>
          </div>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Failure Injection</h2>
        <p className="text-zinc-300 mb-6">
          Inject chaos into the demo environment to trigger CortexOps&apos; correlation and remediation lifecycle.
        </p>
        <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300 mb-4">
          <div className="flex items-center gap-3">
            <Terminal className="w-4 h-4 text-zinc-500" />
            <span>make demo-failure SCENARIO=rollout-fail</span>
          </div>
        </div>
        <ul className="list-disc pl-5 text-sm text-zinc-400 mb-6">
          <li><strong>rollout-fail</strong>: Simulates a bad Kubernetes deployment.</li>
          <li><strong>crashloop</strong>: Forces a worker into CrashLoopBackOff.</li>
          <li><strong>scale-pressure</strong>: Triggers CPU starvation.</li>
        </ul>
        <div className="bg-cortex-500/10 border border-cortex-500/30 p-6 rounded-xl">
          <h3 className="text-lg font-medium text-cortex-300 mb-3">What happens when it fails?</h3>
          <p className="text-sm text-zinc-300">
            During these failure injections, CortexOps evaluates OPA policies. If a proposed remediation fails validation or lacks human approval, the workflow <strong>fails closed</strong>. The infrastructure is never mutated without verifiable safety guarantees.
          </p>
        </div>
      </section>
    </>
  );
}
