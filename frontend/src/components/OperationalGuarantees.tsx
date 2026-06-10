"use client";

import { motion } from "framer-motion";
import { CheckCircle } from "lucide-react";

const guarantees = [
  { title: "Deterministic Execution", desc: "Workflows produce predictable outcomes under retries and failures." },
  { title: "Replay Safety", desc: "Workflow re-execution does not create unintended side effects." },
  { title: "Fail-Closed Governance", desc: "Unsafe actions are blocked before infrastructure mutation occurs." },
  { title: "Rollback Protection", desc: "Remediation workflows verify stabilization before completion." }
];

export function OperationalGuarantees() {
  return (
    <section className="py-24 px-6 max-w-7xl mx-auto">
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.6 }}
        className="text-center mb-16"
      >
        <h2 className="text-3xl font-bold text-white mb-4">Operational Guarantees</h2>
        <p className="text-zinc-400 max-w-2xl mx-auto">Safety and predictability built into every remediation action.</p>
      </motion.div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {guarantees.map((item, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5, delay: i * 0.1 }}
            className="glass-panel p-8 rounded-2xl border border-zinc-800 flex flex-col gap-4 group hover:border-cortex-500/50 transition-colors"
          >
            <div className="flex items-center gap-3">
              <CheckCircle className="w-6 h-6 text-cortex-400 group-hover:text-cortex-300 transition-colors" />
              <h3 className="text-xl font-semibold text-white">{item.title}</h3>
            </div>
            <p className="text-zinc-400 leading-relaxed pl-9">
              {item.desc}
            </p>
          </motion.div>
        ))}
      </div>
    </section>
  );
}
