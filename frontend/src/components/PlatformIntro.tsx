"use client";

import { motion } from "framer-motion";
import { CheckCircle2 } from "lucide-react";

const capabilities = [
  "Understanding service relationships",
  "Detecting correlated failures",
  "Calculating blast radius",
  "Generating root cause hypotheses",
  "Coordinating safe remediation workflows"
];

export function PlatformIntro() {
  return (
    <section className="py-24 px-6 max-w-7xl mx-auto flex flex-col md:flex-row items-center gap-16">
      <div className="flex-1">
        <motion.h2 
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
          className="text-4xl md:text-5xl font-bold text-white mb-6"
        >
          A Control Plane for Infrastructure Reliability
        </motion.h2>
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6, delay: 0.1 }}
          className="text-lg text-zinc-400 font-light leading-relaxed mb-8"
        >
          <p className="mb-4">Modern cloud-native systems generate enormous volumes of events, metrics, traces, and operational signals.</p>
          <p className="mb-6">CortexOps transforms this telemetry into actionable intelligence by:</p>
          <ul className="space-y-4 mb-8">
            {capabilities.map((cap, i) => (
              <li key={i} className="flex items-start gap-3">
                <CheckCircle2 className="w-6 h-6 text-cortex-500 flex-shrink-0" />
                <span className="text-zinc-300">{cap}</span>
              </li>
            ))}
          </ul>
          <p className="text-white font-medium">The result is faster incident resolution without sacrificing operational safety.</p>
        </motion.div>
      </div>
      <div className="flex-1 relative w-full h-[400px] glass-panel rounded-3xl overflow-hidden flex items-center justify-center border border-zinc-800">
         <div className="absolute inset-0 bg-gradient-to-br from-cortex-500/10 to-transparent"></div>
         <div className="w-3/4 h-3/4 border border-zinc-800/50 rounded-2xl relative">
           <div className="absolute top-4 left-4 right-4 h-8 bg-zinc-900 rounded flex items-center px-4 gap-2 border border-zinc-800">
             <div className="w-2 h-2 rounded-full bg-red-500"></div>
             <div className="w-2 h-2 rounded-full bg-yellow-500"></div>
             <div className="w-2 h-2 rounded-full bg-green-500"></div>
             <div className="ml-4 h-2 w-24 bg-zinc-700 rounded"></div>
           </div>
           <div className="absolute top-16 left-4 right-4 bottom-4 flex flex-col gap-2">
             {[35, 56, 22, 47, 31].map((width, i) => (
                <div key={i} className="h-6 w-full bg-zinc-900/50 rounded border border-zinc-800/30 flex items-center px-2">
                  <div className="h-2 bg-cortex-500/50 rounded" style={{ width: `${width}%` }}></div>
                </div>
             ))}
           </div>
         </div>
      </div>
    </section>
  );
}
