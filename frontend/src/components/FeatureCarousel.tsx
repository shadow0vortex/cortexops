"use client";

import { motion, Variants } from "framer-motion";
import { Network, GitMerge, RefreshCcw, BrainCircuit, ShieldAlert, ActivitySquare } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";

const features = [
  { id: "telemetry", title: "Telemetry Ingestion", icon: ActivitySquare, desc: "Process massive streams of Kubernetes events with robust backpressure.", list: ["Protobuf normalization", "NATS JetStream routing", "High-throughput parsing", "Event buffering", "Metric extraction"] },
  { id: "topology", title: "Topology Intelligence", icon: Network, desc: "Maintain a live dependency graph of workloads, services, infrastructure resources, and operational relationships.", list: ["Dependency discovery", "Blast radius analysis", "Service relationship mapping", "Failure propagation modeling", "Infrastructure awareness"] },
  { id: "correlation", title: "Event Correlation", icon: GitMerge, desc: "Convert fragmented operational signals into coherent incidents.", list: ["Temporal correlation", "Trace affinity detection", "Topology-aware scoring", "Incident grouping", "Duplicate suppression"] },
  { id: "rca", title: "RCA Engine", icon: BrainCircuit, desc: "Operational recommendations grounded in telemetry and historical context.", list: ["Incident summarization", "Failure pattern recognition", "Context-aware recommendations", "Retrieval-augmented analysis", "Degraded-mode fallback"] },
  { id: "remediation", title: "Remediation Engine", icon: ShieldAlert, desc: "Every action is validated before execution.", list: ["Policy evaluation via OPA", "Action approval workflows", "Governance controls", "Rollback protection", "Fail-closed execution"] },
  { id: "replay", title: "Replay Safety", icon: RefreshCcw, desc: "Durable execution powered by Temporal.", list: ["Deterministic workflows", "Automatic retries", "Idempotent recovery", "State persistence", "Workflow replay guarantees"] },
];

const containerVariants: Variants = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: { staggerChildren: 0.1 }
  }
};

const itemVariants: Variants = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0, transition: { type: "spring", stiffness: 300, damping: 24 } }
};

export function FeatureCarousel() {
  const [activeIdx, setActiveIdx] = useState(0);

  return (
    <section className="py-24 px-6 max-w-7xl mx-auto flex flex-col items-center">
      <motion.div 
        initial={{ opacity: 0, y: -20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true, margin: "-100px" }}
        transition={{ duration: 0.6 }}
        className="flex flex-col items-center mb-16 text-center"
      >
        <h2 className="text-sm font-semibold tracking-widest text-cortex-400 uppercase mb-3">Core Capabilities</h2>
        <p className="text-3xl font-bold text-white max-w-2xl">
          Everything required for intelligent orchestration.
        </p>
      </motion.div>

      <motion.div 
        variants={containerVariants}
        initial="hidden"
        whileInView="show"
        viewport={{ once: true, margin: "-50px" }}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 w-full relative"
      >
        {/* Subtle background connection lines */}
        <div className="absolute inset-0 pointer-events-none overflow-hidden flex items-center justify-center opacity-20">
            <svg className="w-full h-full" viewBox="0 0 1000 400" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M100 200 L900 200" stroke="url(#paint0_linear)" strokeWidth="2" strokeDasharray="4 4" />
              <defs>
                <linearGradient id="paint0_linear" x1="100" y1="200" x2="900" y2="200" gradientUnits="userSpaceOnUse">
                  <stop stopColor="#A855F7" stopOpacity="0" />
                  <stop offset="0.5" stopColor="#A855F7" />
                  <stop offset="1" stopColor="#A855F7" stopOpacity="0" />
                </linearGradient>
              </defs>
            </svg>
        </div>

        {features.map((feat, idx) => {
          const isActive = activeIdx === idx;
          const Icon = feat.icon;
          
          return (
            <motion.div
              key={feat.id}
              variants={itemVariants}
              onHoverStart={() => setActiveIdx(idx)}
              className={cn(
                "relative group rounded-2xl p-[1px] transition-all duration-500 overflow-hidden cursor-pointer",
                isActive ? "bg-gradient-to-b from-cortex-400/50 to-transparent shadow-[0_0_30px_rgba(168,85,247,0.15)] scale-[1.02]" : "bg-zinc-800/50 hover:bg-zinc-700/50"
              )}
            >
              <div className="h-full bg-[#0B0B0F] rounded-2xl p-6 relative z-10 flex flex-col gap-4">
                <div className={cn(
                  "w-12 h-12 rounded-xl flex items-center justify-center border transition-all duration-500",
                  isActive ? "bg-cortex-900/50 border-cortex-400 shadow-[0_0_15px_rgba(168,85,247,0.4)] text-cortex-300 scale-110" : "bg-zinc-900 border-zinc-700 text-zinc-400 group-hover:text-zinc-200"
                )}>
                  <Icon className="w-6 h-6" />
                </div>
                <div>
                  <h3 className={cn("text-lg font-semibold mb-2 transition-colors duration-300", isActive ? "text-white" : "text-zinc-300")}>
                    {feat.title}
                  </h3>
                  <p className="text-zinc-400 text-sm leading-relaxed mb-4">
                    {feat.desc}
                  </p>
                  <ul className="flex flex-col gap-2">
                    {feat.list.map((item, i) => (
                      <li key={i} className="text-xs text-zinc-500 flex items-start gap-2">
                        <span className="w-1 h-1 rounded-full bg-cortex-500 mt-1.5 flex-shrink-0"></span>
                        <span className={isActive ? "text-zinc-300 transition-colors" : "transition-colors"}>{item}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </motion.div>
          );
        })}
      </motion.div>
    </section>
  );
}
