"use client";

import { motion } from "framer-motion";

const steps = [
  "Kubernetes",
  "Telemetry Collection",
  "NATS JetStream",
  "Correlation Engine",
  "Topology Intelligence",
  "AI Analysis",
  "Temporal Workflows",
  "OPA Governance",
  "Remediation Execution"
];

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

      <div className="w-full h-[600px] flex items-center justify-center perspective-[2000px] mt-12">
         <motion.div 
           initial={{ rotateX: 60, rotateZ: -45 }}
           animate={{ rotateZ: [-45, -35, -45] }}
           transition={{ duration: 30, repeat: Infinity, ease: "linear" }}
           style={{ transformStyle: 'preserve-3d' }}
           className="relative w-72 h-72"
         >
           {/* Base Floor */}
           <div className="absolute inset-[-100%] bg-cortex-500/5 blur-[100px] rounded-full translate-z-[-50px]"></div>

           {[...steps].reverse().map((step, i) => (
             <motion.div
               key={step}
               initial={{ translateZ: i * 80 + 800, opacity: 0 }}
               whileInView={{ translateZ: i * 80, opacity: 1 }}
               viewport={{ once: true, margin: "-200px" }}
               transition={{ duration: 1.2, delay: i * 0.15, type: "spring", stiffness: 60 }}
               style={{ transformStyle: 'preserve-3d' }}
               className="absolute inset-0 group cursor-pointer"
             >
                <div className="absolute inset-0 bg-black/40 group-hover:bg-black/0 transition-colors duration-500 z-20 pointer-events-none" style={{ transform: 'translateZ(1px)' }}></div>
                
                {/* Top Face */}
                <div className="absolute inset-0 bg-zinc-900/90 border border-cortex-500/40 flex items-center justify-center group-hover:bg-cortex-900/90 transition-all duration-300 shadow-[inset_0_0_30px_rgba(168,85,247,0.15)] group-hover:shadow-[inset_0_0_50px_rgba(168,85,247,0.4)] backdrop-blur-md">
                  <span className="text-white font-bold tracking-wider text-lg">{step}</span>
                </div>
                
                {/* Right Face */}
                <div className="absolute top-0 right-0 w-6 h-full bg-zinc-950 border-y border-r border-cortex-500/40 origin-right brightness-50 group-hover:brightness-75 transition-all" style={{ transform: 'translateX(100%) rotateY(90deg)' }}></div>
                
                {/* Front Face */}
                <div className="absolute bottom-0 left-0 w-full h-6 bg-black border-x border-b border-cortex-500/40 origin-bottom brightness-75 group-hover:brightness-100 transition-all" style={{ transform: 'translateY(100%) rotateX(-90deg)' }}></div>
                
                {/* Connection Beam (except top) */}
                {i < steps.length - 1 && (
                  <div className="absolute top-1/2 left-1/2 w-2 h-[80px] bg-cortex-500/30 blur-sm" style={{ transform: 'translate(-50%, -50%) rotateX(-90deg) translateZ(40px)' }}></div>
                )}
             </motion.div>
           ))}
         </motion.div>
      </div>
    </section>
  );
}
