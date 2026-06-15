"use client";

import { Activity, Layers, ShieldCheck, Zap, Repeat } from "lucide-react";

export function StatusStrip() {
  const metrics = [
    { icon: <Zap className="w-5 h-5 text-cortex-400" />, label: "39,000+", sub: "Events/sec" },
    { icon: <Layers className="w-5 h-5 text-cortex-400" />, label: "100k", sub: "Event Storm Tested" },
    { icon: <ShieldCheck className="w-5 h-5 text-cortex-400" />, label: "Validated", sub: "Replay Safety" },
    { icon: <Activity className="w-5 h-5 text-cortex-400" />, label: "Topology-Aware", sub: "Correlation" },
    { icon: <Repeat className="w-5 h-5 text-cortex-400" />, label: "Durable", sub: "Temporal Workflows" },
    { icon: <ShieldCheck className="w-5 h-5 text-cortex-400" />, label: "OPA", sub: "Policy Enforcement" },
  ];

  return (
    <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 -mt-10 relative z-30">
      <div
        className="glass-panel rounded-2xl p-4 sm:p-6 flex flex-wrap lg:flex-nowrap items-center justify-between gap-4 sm:gap-6 border-t border-t-cortex-500/20"
        role="region"
        aria-label="Platform metrics"
      >
        {/* Header inside strip */}
        <div className="flex items-center gap-3 w-full lg:w-auto mb-2 lg:mb-0 border-b lg:border-b-0 lg:border-r border-zinc-800 pb-3 lg:pb-0 lg:pr-6">
          <div className="w-2 h-2 rounded-full bg-cortex-400 animate-slow-pulse" aria-hidden="true"></div>
          <span className="text-sm font-medium text-zinc-300 tracking-wide uppercase">System Status</span>
        </div>

        {/* Metrics */}
        <div className="flex-1 grid grid-cols-2 sm:grid-cols-3 lg:flex lg:flex-nowrap lg:justify-between gap-3 sm:gap-4 w-full">
          {metrics.map((m, i) => (
            <div key={i} className="flex items-center gap-2 sm:gap-3">
              <div className="w-8 h-8 sm:w-10 sm:h-10 rounded-xl bg-cortex-900/50 border border-cortex-500/20 flex items-center justify-center shadow-[0_0_10px_rgba(168,85,247,0.1)] flex-shrink-0" aria-hidden="true">
                {m.icon}
              </div>
              <div className="flex flex-col min-w-0">
                <span className="text-white font-semibold text-xs sm:text-sm truncate">{m.label}</span>
                <span className="text-zinc-500 text-[10px] sm:text-xs truncate">{m.sub}</span>
              </div>
            </div>
          ))}
        </div>

      </div>
    </div>
  );
}
