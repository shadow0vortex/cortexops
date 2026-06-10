"use client";

import { motion } from "framer-motion";
import { Sparkles, RefreshCcw, ShieldAlert } from "lucide-react";

export default function SecurityPage() {
  const guarantees = [
    {
      id: "deterministic",
      title: "Deterministic Execution",
      icon: Sparkles,
      desc: "Topology-aware correlation, replay-safe workflows, and policy-governed remediation for modern infrastructure operations.",
    },
    {
      id: "replay",
      title: "Replay Safety",
      icon: RefreshCcw,
      desc: "Topology-aware correlation and replay-safe workflows and policy-governed remediation for modern infrastructure operations.",
    },
    {
      id: "fail-closed",
      title: "Fail-Closed Governance",
      icon: ShieldAlert,
      desc: "Deterministic governance ensures only safe, compliant, and policy-approved actions are executed with automatic rollback validation.",
    }
  ];

  return (
    <div className="pt-40 pb-24 px-6 max-w-6xl mx-auto flex flex-col items-center min-h-screen">
      
      <div className="flex flex-col items-center text-center mb-16 relative">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-cortex-600/20 rounded-full blur-[100px] -z-10 mix-blend-screen pointer-events-none"></div>
        <div className="flex items-center gap-2 px-3 py-1.5 rounded-full border border-cortex-500/30 bg-cortex-500/10 text-cortex-300 text-xs font-medium mb-6">
          <span className="w-1.5 h-1.5 rounded-full bg-cortex-400"></span>
          Kubernetes-Native Infrastructure Intelligence
        </div>
        <h1 className="text-4xl md:text-5xl font-bold mb-4 text-white">
          Deterministic <br/>Execution & Guarantees
        </h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
        {guarantees.map((item, index) => {
          const Icon = item.icon;
          return (
            <motion.div
              key={item.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className="glass-panel p-8 rounded-2xl flex flex-col items-start gap-4 group hover:border-cortex-500/50 transition-all relative overflow-hidden"
            >
              <div className="absolute top-0 right-0 w-32 h-32 bg-cortex-500/10 rounded-bl-full opacity-0 group-hover:opacity-100 transition-opacity"></div>
              
              <Icon className="w-8 h-8 text-white group-hover:text-cortex-300 transition-colors" strokeWidth={1.5} />
              <h3 className="text-xl font-bold text-white">{item.title}</h3>
              <p className="text-zinc-400 text-sm leading-relaxed">{item.desc}</p>
            </motion.div>
          )
        })}
      </div>

    </div>
  );
}
