import { SystemArchitecture } from "@/components/diagrams/SystemArchitecture";
import { RemediationLifecycle } from "@/components/diagrams/RemediationLifecycle";
import { ReplayRecovery } from "@/components/diagrams/ReplayRecovery";

export default function ArchitecturePage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Architecture</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        A resilient, event-driven control plane built for operational correctness.
      </p>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">System Overview</h2>
        <SystemArchitecture />
        <div className="mt-8 space-y-4 text-zinc-300">
          <p>
            CortexOps operates strictly outside the critical path of your application traffic. It functions as an out-of-band control plane listening to operational telemetry.
          </p>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Telemetry Ingestion & Topology</h2>
        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl mb-6">
          <h3 className="text-lg font-medium text-white mb-3">How does it work?</h3>
          <p className="text-sm text-zinc-400">
            The Collector service watches the Kubernetes API for events, metrics, and state changes. These are normalized into a standard protobuf schema and published to NATS JetStream. Concurrently, the Topology Engine maintains an in-memory graph of all resources to instantly calculate blast radius.
          </p>
        </div>
        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl mb-6">
          <h3 className="text-lg font-medium text-white mb-3">Why was it designed this way?</h3>
          <p className="text-sm text-zinc-400">
            An in-memory graph ensures sub-millisecond queries during a cascading failure. Normalizing via NATS ensures that the correlation engine receives ordered, replayable event streams, decoupling ingestion from processing.
          </p>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Remediation Orchestration</h2>
        <p className="text-zinc-300 mb-6">
          When an incident is correlated and a root cause is proposed, the system orchestrates a deterministic remediation workflow using Temporal.
        </p>
        <RemediationLifecycle />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
          <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl">
            <h3 className="text-lg font-medium text-white mb-3">Operational Guarantees</h3>
            <p className="text-sm text-zinc-400">
              Every workflow state is durably persisted. If the remediation worker crashes during execution, Temporal resumes the exact step upon restart. OPA policies guarantee that no forbidden mutations occur.
            </p>
          </div>
          <div className="bg-cortex-500/10 border border-cortex-500/30 p-6 rounded-xl">
            <h3 className="text-lg font-medium text-cortex-300 mb-3">What happens when it fails?</h3>
            <p className="text-sm text-zinc-300">
              If a patch fails or post-execution telemetry degrades, the workflow automatically transitions to the <code>ROLLING_BACK</code> state. If rollback fails, the incident escalates to PagerDuty.
            </p>
          </div>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Replay Safety & Chaos Degradation</h2>
        <p className="text-zinc-300 mb-6">
          CortexOps guarantees exactly-once processing semantics during network partitions or pod failures.
        </p>
        <ReplayRecovery />
        <div className="mt-8">
          <h3 className="text-lg font-medium text-white mb-3">Why was it designed this way?</h3>
          <p className="text-sm text-zinc-400">
            In distributed systems, failures are inevitable. If NATS redelivers an event due to a timeout, the Correlation Engine must deterministically drop the duplicate using the Audit DB to prevent triggering duplicate remediation workflows.
          </p>
        </div>
      </section>
    </>
  );
}
