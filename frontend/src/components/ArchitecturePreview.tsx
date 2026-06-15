"use client";

import { motion } from "framer-motion";

export function ArchitecturePreview() {
  return (
    <section className="py-32 px-6 max-w-7xl mx-auto flex flex-col items-center overflow-hidden">
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.6 }}
        className="text-center mb-24 z-10"
      >
        <h2 className="text-sm font-semibold tracking-widest text-cortex-400 uppercase mb-3">Distributed Systems by Design</h2>
        <p className="text-3xl font-bold text-white max-w-2xl mx-auto mb-4">
          CortexOps follows an event-driven architecture designed for resilience and operational correctness.
        </p>
        <p className="text-zinc-400">
          Every component is independently deployable and horizontally scalable.
        </p>
      </motion.div>

      {/* High-Quality Animated Data Flow Diagram */}
      <div className="w-full max-w-5xl mx-auto mt-16 relative">
        {/* Background Grid */}
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_50%,#000_70%,transparent_100%)]"></div>
        
        <div className="relative flex flex-col md:flex-row items-center justify-between gap-8 md:gap-4 py-12 px-4 z-10">
          
          {[
            { title: "Telemetry Ingestion", desc: "K8s Events", icon: "activity" },
            { title: "NATS JetStream", desc: "Event Bus", icon: "network" },
            { title: "Correlation Engine", desc: "Topology Intelligence", icon: "brain" },
            { title: "Temporal", desc: "Durable Workflows", icon: "history" },
            { title: "Remediation", desc: "Policy Executed", icon: "shield" }
          ].map((node, i, arr) => (
            <div key={node.title} className="flex-1 flex flex-col items-center relative w-full md:w-auto">
              {/* Connection Line */}
              {i < arr.length - 1 && (
                <div className="hidden md:block absolute top-10 left-[60%] w-[80%] h-[2px] bg-zinc-800">
                  <motion.div 
                    className="h-full bg-gradient-to-r from-transparent via-cortex-400 to-transparent w-1/2"
                    animate={{ x: ["-100%", "200%"] }}
                    transition={{ duration: 2, repeat: Infinity, ease: "linear", delay: i * 0.4 }}
                  />
                </div>
              )}
              
              <motion.div 
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.2 }}
                className="relative group w-20 h-20 mb-4 flex items-center justify-center"
              >
                {/* Node Glow */}
                <div className="absolute inset-0 bg-cortex-500/20 rounded-2xl blur-xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                
                {/* Node Container */}
                <div className="absolute inset-0 bg-zinc-900 border border-zinc-700 rounded-2xl shadow-xl group-hover:border-cortex-400/50 group-hover:scale-110 transition-all duration-300 flex items-center justify-center overflow-hidden">
                  <div className="absolute inset-0 bg-gradient-to-b from-white/5 to-transparent"></div>
                  
                  {node.icon === "activity" && <svg className="w-8 h-8 text-cortex-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path strokeLinecap="round" strokeLinejoin="round" d="M22 12h-4l-3 9L9 3l-3 9H2" /></svg>}
                  {node.icon === "network" && <svg className="w-8 h-8 text-cortex-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path strokeLinecap="round" strokeLinejoin="round" d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>}
                  {node.icon === "brain" && <svg className="w-8 h-8 text-cortex-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path strokeLinecap="round" strokeLinejoin="round" d="M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44 2.5 2.5 0 0 1-2.96-3.08 3 3 0 0 1-.34-5.58 2.5 2.5 0 0 1 1.32-4.24 2.5 2.5 0 0 1 1.98-3A2.5 2.5 0 0 1 9.5 2Z"></path><path strokeLinecap="round" strokeLinejoin="round" d="M14.5 2A2.5 2.5 0 0 0 12 4.5v15a2.5 2.5 0 0 0 4.96.44 2.5 2.5 0 0 0 2.96-3.08 3 3 0 0 0 .34-5.58 2.5 2.5 0 0 0-1.32-4.24 2.5 2.5 0 0 0-1.98-3A2.5 2.5 0 0 0 14.5 2Z"></path></svg>}
                  {node.icon === "history" && <svg className="w-8 h-8 text-cortex-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>}
                  {node.icon === "shield" && <svg className="w-8 h-8 text-cortex-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path strokeLinecap="round" strokeLinejoin="round" d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>}
                </div>
              </motion.div>
              
              <motion.div 
                initial={{ opacity: 0 }}
                whileInView={{ opacity: 1 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.2 + 0.3 }}
                className="text-center"
              >
                <h3 className="text-white font-semibold text-sm mb-1">{node.title}</h3>
                <p className="text-zinc-500 text-xs uppercase tracking-widest">{node.desc}</p>
              </motion.div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
