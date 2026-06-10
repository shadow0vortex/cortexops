"use client";

import { motion, useMotionValue, useSpring, useTransform } from "framer-motion";
import Link from "next/link";
import { ChevronRight } from "lucide-react";

export function Hero() {
  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);

  // 3D Tilt based on mouse position
  const rotateX = useSpring(useTransform(mouseY, [-0.5, 0.5], [25, -25]), { stiffness: 150, damping: 20 });
  const rotateY = useSpring(useTransform(mouseX, [-0.5, 0.5], [-25, 25]), { stiffness: 150, damping: 20 });

  function handleMouseMove(event: React.MouseEvent<HTMLElement>) {
    const rect = event.currentTarget.getBoundingClientRect();
    const xPct = (event.clientX - rect.left) / rect.width - 0.5;
    const yPct = (event.clientY - rect.top) / rect.height - 0.5;
    mouseX.set(xPct);
    mouseY.set(yPct);
  }

  function handleMouseLeave() {
    mouseX.set(0);
    mouseY.set(0);
  }

  return (
    <section 
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      className="relative pt-40 pb-20 px-6 max-w-7xl mx-auto flex flex-col lg:flex-row items-center gap-16 min-h-[90vh]"
      style={{ perspective: 1500 }}
    >
      
      {/* Left Content */}
      <div className="flex-1 flex flex-col items-start z-10">
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
          className="flex items-center gap-2 px-3 py-1.5 rounded-full border border-cortex-500/30 bg-cortex-500/10 text-cortex-300 text-sm font-medium mb-8 backdrop-blur-md"
        >
          <span className="w-2 h-2 rounded-full bg-cortex-400 animate-slow-pulse shadow-[0_0_8px_rgba(168,85,247,0.8)]"></span>
          A Control Plane for Infrastructure Reliability
        </motion.div>

        <motion.h1 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.1, ease: [0.16, 1, 0.3, 1] }}
          className="text-5xl md:text-7xl font-bold tracking-tight text-white mb-6 leading-[1.1]"
        >
          Infrastructure <br />
          Intelligence <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-cortex-300 via-cortex-500 to-cortex-600">
            Platform
          </span>
        </motion.h1>

        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2, ease: [0.16, 1, 0.3, 1] }}
          className="text-lg text-zinc-400 max-w-xl mb-10 leading-relaxed font-light"
        >
          <p className="mb-4">
            CortexOps continuously analyzes infrastructure telemetry, correlates distributed failures, evaluates remediation policies, and orchestrates deterministic recovery workflows across Kubernetes environments.
          </p>
          <p>
            Built around event-driven architecture, topology awareness, replay-safe execution, and policy-governed automation.
          </p>
        </motion.div>

        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.3, ease: [0.16, 1, 0.3, 1] }}
          className="flex flex-wrap items-center gap-4"
        >
          <Link 
            href="/platform"
            className="flex items-center gap-2 px-6 py-3 rounded-full bg-white text-black font-semibold hover:bg-zinc-200 hover:scale-105 active:scale-95 transition-all duration-300 shadow-[0_0_20px_rgba(255,255,255,0.2)]"
          >
            Explore Platform <ChevronRight className="w-4 h-4" />
          </Link>
          <Link 
            href="/architecture"
            className="flex items-center gap-2 px-6 py-3 rounded-full border border-zinc-700 bg-zinc-900/50 text-white font-medium hover:bg-zinc-800 hover:border-zinc-500 transition-all duration-300"
          >
            View Architecture
          </Link>
        </motion.div>
      </div>

      {/* Right Graphic (Topology Visualization) */}
      <motion.div 
        style={{ rotateX, rotateY, transformStyle: "preserve-3d" }}
        className="flex-1 relative w-full aspect-square max-w-2xl lg:max-w-none flex items-center justify-center pointer-events-none z-0"
      >
        <motion.div 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 1.5 }}
          className="absolute inset-0 bg-cortex-600/10 rounded-full blur-[120px] mix-blend-screen"
        ></motion.div>
        
        {/* Core Node */}
        <motion.div 
          animate={{ scale: [1, 1.05, 1], boxShadow: ["0 0 40px rgba(168,85,247,0.4)", "0 0 100px rgba(168,85,247,0.8)", "0 0 40px rgba(168,85,247,0.4)"] }}
          transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
          style={{ translateZ: 150 }}
          className="absolute w-32 h-32 rounded-3xl border border-cortex-400/50 bg-black/60 backdrop-blur-xl flex items-center justify-center z-20"
        >
          <div className="w-24 h-24 rounded-xl flex items-center justify-center relative overflow-hidden mix-blend-screen">
            <img src="/logo.png" alt="CortexOps Logo" className="w-full h-full object-contain drop-shadow-[0_0_20px_rgba(168,85,247,0.8)] scale-110" />
          </div>
        </motion.div>

        {/* Orbit Rings with 3D Depth */}
        <motion.div style={{ translateZ: -100 }} className="absolute w-[120%] h-[120%] rounded-full border border-cortex-500/20 rotate-x-60 animate-[spin_40s_linear_infinite]"></motion.div>
        <motion.div style={{ translateZ: 0 }} className="absolute w-[90%] h-[90%] rounded-full border border-dashed border-cortex-500/40 rotate-x-60 animate-[spin_25s_linear_infinite_reverse]"></motion.div>
        <motion.div style={{ translateZ: 80 }} className="absolute w-[60%] h-[60%] rounded-full border border-cortex-500/50 rotate-x-60 animate-[spin_15s_linear_infinite]"></motion.div>
        
        {/* Distributed Nodes & Telemetry Pulses */}
        {[...Array(6)].map((_, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, scale: 0 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.5 + i * 0.1, duration: 0.8, ease: "easeOut" }}
            className="absolute w-12 h-12 rounded-xl border border-zinc-700 bg-zinc-900/80 backdrop-blur-md z-10 flex items-center justify-center shadow-lg"
            style={{
              translateZ: (i % 2 === 0 ? 50 : 200),
              transform: `rotate(${i * 60}deg) translateY(-220px) rotate(-${i * 60}deg)`
            }}
          >
            <motion.div 
              animate={{ opacity: [0.3, 1, 0.3] }}
              transition={{ duration: 2, delay: i * 0.3, repeat: Infinity }}
              className="w-2 h-2 rounded-full bg-cortex-400 shadow-[0_0_10px_rgba(168,85,247,0.8)]"
            ></motion.div>
            
            {/* Connecting line */}
            <div 
              className="absolute top-1/2 left-1/2 w-[220px] h-[1px] bg-gradient-to-r from-cortex-500/40 to-transparent -z-10"
              style={{ transform: `rotate(${i * 60 + 90}deg)`, transformOrigin: '0 0' }}
            ></div>

            {/* Telemetry packet travelling along line */}
            <motion.div
               animate={{ x: [0, 220], opacity: [0, 1, 0] }}
               transition={{ duration: 1.5, repeat: Infinity, delay: i * 0.5, ease: "linear" }}
               className="absolute top-1/2 left-1/2 w-8 h-[2px] bg-gradient-to-r from-transparent via-white to-transparent -z-10"
               style={{ transformOrigin: '0 0', rotate: `${i * 60 + 90}deg` }}
            />
          </motion.div>
        ))}
      </motion.div>
    </section>
  );
}
