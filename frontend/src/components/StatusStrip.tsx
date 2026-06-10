"use client";

import { Activity, Layers, ShieldCheck, Zap, Repeat } from "lucide-react";

export function StatusStrip() {
  const metrics = [
    { icon: <Zap className="w-5 h-5 text-cortex-400" />, label: "<50ms", sub: "Telemetry Ingestion" },
    { icon: <Layers className="w-5 h-5 text-cortex-400" />, label: "1000+", sub: "Concurrent Windows" },
    { icon: <ShieldCheck className="w-5 h-5 text-cortex-400" />, label: "Exactly-Once", sub: "Event Processing" },
    { icon: <Activity className="w-5 h-5 text-cortex-400" />, label: "Real-Time", sub: "Blast Radius Analysis" },
    { icon: <Repeat className="w-5 h-5 text-cortex-400" />, label: "Deterministic", sub: "Workflow Recovery" },
  ];

  return (
    <div className="w-full max-w-7xl mx-auto px-6 -mt-10 relative z-30">
      <div className="glass-panel rounded-2xl p-6 flex flex-wrap lg:flex-nowrap items-center justify-between gap-6 border-t border-t-cortex-500/20">
        
        {/* Header inside strip */}
        <div className="flex items-center gap-3 w-full lg:w-auto mb-4 lg:mb-0 border-b lg:border-b-0 lg:border-r border-zinc-800 pb-4 lg:pb-0 lg:pr-6">
          <div className="w-2 h-2 rounded-full bg-cortex-400 animate-slow-pulse"></div>
          <span className="text-sm font-medium text-zinc-300 tracking-wide uppercase">System Status</span>
        </div>

        {/* Metrics */}
        <div className="flex-1 flex flex-wrap sm:flex-nowrap justify-between gap-4 w-full">
          {metrics.map((m, i) => (
            <div key={i} className="flex items-center gap-3 min-w-[140px]">
              <div className="w-10 h-10 rounded-xl bg-cortex-900/50 border border-cortex-500/20 flex items-center justify-center shadow-[0_0_10px_rgba(168,85,247,0.1)]">
                {m.icon}
              </div>
              <div className="flex flex-col">
                <span className="text-white font-semibold text-sm">{m.label}</span>
                <span className="text-zinc-500 text-xs">{m.sub}</span>
              </div>
            </div>
          ))}
        </div>

      </div>
    </div>
  );
}
