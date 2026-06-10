"use client";

import { motion } from "framer-motion";
import { Server, Database, GitMerge, RefreshCcw, ShieldCheck, Zap } from "lucide-react";

export default function ArchitecturePage() {
  const nodes = [
    { id: 1, title: "Telemetry Ingestion", icon: Server, desc: "client-go Informers outputting strongly-typed Protobufs via NATS JetStream" },
    { id: 2, title: "Topology Graph", icon: Database, desc: "In-memory Directed Graph + Asynchronous PostgreSQL for blast-radius analysis" },
    { id: 3, title: "Correlation Lifecycle", icon: GitMerge, desc: "Temporal windowing and heuristic scoring to detect causality" },
    { id: 4, title: "Remediation Workflow", icon: ShieldCheck, desc: "Orchestrated OPA execution with Dry-Run -> Execute -> Verify -> Rollback" },
    { id: 5, title: "Replay Recovery", icon: RefreshCcw, desc: "Durable Temporal workflows with state persistence and automatic retries" },
  ];

  return (
    <div className="pt-32 pb-24 px-6 max-w-5xl mx-auto flex flex-col items-center min-h-screen">
      <div className="text-center mb-16 relative">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-64 h-64 bg-cortex-600/30 rounded-full blur-[100px] -z-10 mix-blend-screen pointer-events-none"></div>
        <h1 className="text-4xl md:text-5xl font-bold mb-4 text-white">Distributed Systems Architecture</h1>
        <p className="text-zinc-400 max-w-xl mx-auto">A high-level view of the CortexOps control plane.</p>
      </div>

      <div className="relative w-full max-w-2xl flex flex-col items-center gap-10">
        
        {/* Vertical glowing pathway with flowing telemetry particles */}
        <div className="absolute top-0 bottom-0 w-[2px] bg-gradient-to-b from-transparent via-cortex-500/20 to-transparent -z-10 overflow-hidden">
          <motion.div 
            className="w-full h-32 bg-gradient-to-b from-transparent via-cortex-400 to-transparent blur-[2px]"
            animate={{ y: [-150, 1000] }}
            transition={{ duration: 3, repeat: Infinity, ease: "linear" }}
          />
          {/* Secondary smaller particle */}
          <motion.div 
            className="w-full h-12 bg-gradient-to-b from-transparent via-white to-transparent blur-[1px] absolute top-0"
            animate={{ y: [-50, 1000] }}
            transition={{ duration: 2.2, repeat: Infinity, ease: "linear", delay: 1 }}
          />
        </div>

        {nodes.map((node, index) => {
          const Icon = node.icon;
          return (
            <motion.div 
              key={node.id}
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-50px" }}
              transition={{ delay: index * 0.1, duration: 0.6, type: "spring", stiffness: 200, damping: 20 }}
              className="relative w-full glass-panel rounded-2xl p-6 flex flex-col items-center text-center group hover:bg-zinc-900/80 transition-colors duration-500"
            >
              <div className="absolute inset-0 bg-gradient-to-b from-cortex-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 rounded-2xl"></div>
              
              <div className="w-16 h-16 rounded-2xl bg-zinc-900 border border-zinc-700 shadow-xl flex items-center justify-center mb-4 z-10 group-hover:border-cortex-400/50 group-hover:bg-black group-hover:shadow-[0_0_30px_rgba(168,85,247,0.4)] group-hover:-translate-y-1 transition-all duration-300">
                <Icon className="w-8 h-8 text-cortex-300 group-hover:text-cortex-400 group-hover:scale-110 transition-all duration-300" strokeWidth={1.5} />
              </div>
              <h3 className="text-xl font-bold text-white mb-2 z-10">{node.title}</h3>
              <p className="text-zinc-400 text-sm z-10">{node.desc}</p>
            </motion.div>
          )
        })}

      </div>
    </div>
  );
}
