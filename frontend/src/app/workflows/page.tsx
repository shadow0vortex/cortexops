import { RemediationLifecycle } from "@/components/diagrams/RemediationLifecycle";
import { Workflow, ShieldAlert, GitBranch, History } from "lucide-react";
import Link from "next/link";

export default function WorkflowsPage() {
  return (
    <div className="flex flex-col relative overflow-hidden pt-32 pb-24 px-6">
      
      {/* Background Gradients */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[500px] bg-[radial-gradient(ellipse_at_top,rgba(168,85,247,0.15),transparent_50%)] pointer-events-none"></div>

      <div className="max-w-7xl mx-auto w-full relative z-10">
        <div className="flex flex-col items-center text-center mb-20">
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-full border border-cortex-500/30 bg-cortex-500/10 text-cortex-300 text-sm font-medium mb-6">
            <Workflow className="w-4 h-4" />
            Deterministic Execution
          </div>
          <h1 className="text-5xl md:text-6xl font-bold tracking-tight text-white mb-6">
            Remediation <span className="text-transparent bg-clip-text bg-gradient-to-r from-cortex-300 to-cortex-600">Workflows</span>
          </h1>
          <p className="text-xl text-zinc-400 max-w-3xl font-light">
            CortexOps leverages Temporal to provide durable, replay-safe, and highly-available remediation pipelines. Never lose state during an infrastructure failure.
          </p>
        </div>

        {/* Interactive Diagram */}
        <div className="mb-24">
          <RemediationLifecycle />
        </div>

        {/* Feature Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-24">
          <div className="glass-panel border border-zinc-800 rounded-3xl p-8 bg-zinc-900/50 hover:bg-zinc-900/80 transition-colors">
            <div className="w-12 h-12 rounded-2xl bg-cortex-500/20 text-cortex-400 flex items-center justify-center mb-6">
              <History className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-white mb-3">Durable Execution</h3>
            <p className="text-zinc-400 text-sm leading-relaxed">
              Every step of a remediation workflow is persisted. If the orchestrator crashes during a node restart, the workflow resumes exactly where it left off upon recovery.
            </p>
          </div>

          <div className="glass-panel border border-zinc-800 rounded-3xl p-8 bg-zinc-900/50 hover:bg-zinc-900/80 transition-colors">
            <div className="w-12 h-12 rounded-2xl bg-cortex-500/20 text-cortex-400 flex items-center justify-center mb-6">
              <ShieldAlert className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-white mb-3">Policy-Governed</h3>
            <p className="text-zinc-400 text-sm leading-relaxed">
              Workflows must pass Open Policy Agent (OPA) checks before modifying cluster state. High-risk actions automatically pause and request human-in-the-loop approval.
            </p>
          </div>

          <div className="glass-panel border border-zinc-800 rounded-3xl p-8 bg-zinc-900/50 hover:bg-zinc-900/80 transition-colors">
            <div className="w-12 h-12 rounded-2xl bg-cortex-500/20 text-cortex-400 flex items-center justify-center mb-6">
              <GitBranch className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-white mb-3">Idempotent Retries</h3>
            <p className="text-zinc-400 text-sm leading-relaxed">
              All automation scripts are wrapped in exponential backoffs with jitter. Safe, deterministic retries guarantee that temporary network blips don&apos;t break the recovery process.
            </p>
          </div>
        </div>

        {/* CTA */}
        <div className="flex justify-center">
          <Link 
            href="/docs/operations"
            className="flex items-center gap-2 px-8 py-4 rounded-full bg-white text-black font-semibold hover:bg-zinc-200 hover:scale-105 active:scale-95 transition-all duration-300 shadow-[0_0_20px_rgba(255,255,255,0.2)]"
          >
            Read the Operations Manual
          </Link>
        </div>

      </div>
    </div>
  );
}
