import { ShieldAlert, GitCommit, FileLock, HardDrive, BrainCircuit } from "lucide-react";

export default function SecurityGovernancePage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Security & Governance</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        Strict boundaries and safety guarantees for autonomous infrastructure operations.
      </p>

      <div className="bg-cortex-900/20 border border-cortex-500/30 rounded-2xl p-8 mb-12">
        <h2 className="text-2xl font-semibold text-white mb-4">The AI Safety Model</h2>
        <p className="text-zinc-300 mb-6">
          CortexOps operates under a strict separation of concerns regarding Artificial Intelligence.
        </p>
        <div className="flex flex-col md:flex-row items-center gap-6">
          <div className="flex-1 bg-black/40 border border-zinc-800 rounded-xl p-6 text-center">
            <BrainCircuit className="w-8 h-8 text-zinc-500 mx-auto mb-3" />
            <h3 className="text-white font-medium mb-2">AI = Recommendation</h3>
            <p className="text-xs text-zinc-400">The LLM analyzes telemetry and proposes a root cause and a suggested fix. It has <strong>zero</strong> direct access to the Kubernetes API.</p>
          </div>
          <div className="text-cortex-500 text-2xl font-bold">+</div>
          <div className="flex-1 bg-black/40 border border-cortex-500/50 rounded-xl p-6 text-center shadow-[0_0_20px_rgba(168,85,247,0.1)]">
            <ShieldAlert className="w-8 h-8 text-cortex-400 mx-auto mb-3" />
            <h3 className="text-white font-medium mb-2">OPA + Temporal = Action</h3>
            <p className="text-xs text-zinc-400">The workflow engine receives the proposal, subjects it to OPA policy evaluation, requests human approval, and executes safely.</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        
        <div>
          <div className="flex items-center gap-3 mb-3">
            <GitCommit className="w-6 h-6 text-cortex-400" />
            <h3 className="text-xl font-semibold text-white">Deterministic Execution</h3>
          </div>
          <p className="text-sm text-zinc-400 leading-relaxed">
            Correlation heuristics and remediation actions are 100% deterministic. Given the same set of telemetry events, CortexOps will always produce the exact same incident grouping and trigger the same workflow path. There are no probabilistic models in the critical execution path.
          </p>
        </div>

        <div>
          <div className="flex items-center gap-3 mb-3">
            <HardDrive className="w-6 h-6 text-cortex-400" />
            <h3 className="text-xl font-semibold text-white">Replay Safety</h3>
          </div>
          <p className="text-sm text-zinc-400 leading-relaxed">
            If a network partition causes NATS JetStream to redeliver an event, or if a Temporal worker crashes and is restarted, the system guarantees that mutations are not applied twice. Idempotency keys and deterministic state machines prevent unintended side effects.
          </p>
        </div>

        <div>
          <div className="flex items-center gap-3 mb-3">
            <FileLock className="w-6 h-6 text-cortex-400" />
            <h3 className="text-xl font-semibold text-white">Policy Enforcement</h3>
          </div>
          <p className="text-sm text-zinc-400 leading-relaxed">
            Every `Execute()` call must pass Open Policy Agent (OPA) validation. Policies are written in Rego and distributed across the cluster. Examples include blocking operations on the `kube-system` namespace or requiring 2-person approvals for statefulset scale-downs.
          </p>
        </div>

        <div>
          <div className="flex items-center gap-3 mb-3">
            <ShieldAlert className="w-6 h-6 text-cortex-400" />
            <h3 className="text-xl font-semibold text-white">Rollback Guarantees</h3>
          </div>
          <p className="text-sm text-zinc-400 leading-relaxed">
            Remediation is not fire-and-forget. The workflow pauses in a `VERIFYING` state to analyze post-patch telemetry. If the system does not stabilize, a pre-calculated compensation transaction is executed to roll back the cluster to its previous state.
          </p>
        </div>

      </div>
    </>
  );
}
