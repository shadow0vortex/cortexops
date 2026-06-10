"use client";

import { motion } from "framer-motion";
import { Activity, AlertCircle, Clock, LayoutDashboard, Plus } from "lucide-react";
import Link from "next/link";

export default function ObservabilityPage() {
  return (
    <div className="pt-32 pb-24 px-6 max-w-7xl mx-auto flex flex-col min-h-screen">
      
      <div className="mb-10 flex justify-between items-end">
        <div>
          <h1 className="text-3xl md:text-4xl font-bold mb-2 text-white">CortexOps Observability</h1>
          <p className="text-zinc-400">Operational Visibility & Real-time Telemetry</p>
        </div>
        <Link href="/observability/kanban" className="hidden sm:block">
          <motion.button 
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl border border-cortex-500/30 bg-cortex-500/10 hover:bg-cortex-500/20 transition-all text-sm font-semibold text-cortex-100 shadow-[0_0_15px_rgba(168,85,247,0.15)] hover:shadow-[0_0_25px_rgba(168,85,247,0.3)] relative overflow-hidden group"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-cortex-400/20 to-transparent translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-700 ease-out"></div>
            <LayoutDashboard className="w-4 h-4 text-cortex-400 group-hover:animate-pulse" />
            <span>Custom Dashboards</span>
            <Plus className="w-4 h-4 text-cortex-400 opacity-70" />
          </motion.button>
        </Link>
      </div>

      {/* Top Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        {/* Active Incidents */}
        <div className="glass-panel p-6 rounded-2xl flex flex-col justify-between h-48 relative overflow-hidden">
          <div className="flex items-center gap-2 text-zinc-400 mb-2">
            <AlertCircle className="w-4 h-4 text-cortex-400" />
            <span className="text-sm font-medium">Active Incidents</span>
          </div>
          <div className="text-6xl font-bold text-white z-10">14</div>
          <div className="text-xs text-zinc-500 mt-2 z-10">Last 24 hours</div>
          {/* Mock Graph line */}
          <div className="absolute bottom-0 left-0 right-0 h-24 opacity-50 pointer-events-none">
            <svg viewBox="0 0 100 40" preserveAspectRatio="none" className="w-full h-full stroke-cortex-400 fill-none" strokeWidth="2">
              <motion.path 
                initial={{ pathLength: 0 }}
                animate={{ pathLength: 1 }}
                transition={{ duration: 1.5, ease: "easeOut" }}
                d="M0,30 C20,30 30,10 50,20 C70,30 80,5 100,15" 
              />
            </svg>
            <div className="absolute inset-0 bg-gradient-to-t from-cortex-500/20 to-transparent"></div>
          </div>
        </div>

        {/* Replay Lag */}
        <div className="glass-panel p-6 rounded-2xl flex flex-col justify-between h-48">
          <div className="flex items-center gap-2 text-zinc-400 mb-2">
            <Clock className="w-4 h-4 text-cortex-400" />
            <span className="text-sm font-medium">Replay Lag</span>
          </div>
          <div className="text-5xl font-bold text-white flex items-baseline gap-1">
            320<span className="text-2xl text-zinc-500">ms</span>
          </div>
          {/* Mock Gauge */}
          <div className="w-full h-12 mt-4 relative">
            <div className="w-full h-2 bg-zinc-800 rounded-full overflow-hidden">
              <motion.div 
                className="h-full bg-gradient-to-r from-cortex-600 to-cortex-400" 
                initial={{ width: 0 }}
                animate={{ width: "32%" }}
                transition={{ duration: 1.5, ease: "easeOut" }}
              />
            </div>
          </div>
        </div>

        {/* Topology Health Score */}
        <div className="glass-panel p-6 rounded-2xl flex flex-col justify-between h-48 relative overflow-hidden">
          <div className="flex items-center gap-2 text-zinc-400 mb-2">
            <Activity className="w-4 h-4 text-emerald-400" />
            <span className="text-sm font-medium">Topology Health Score</span>
          </div>
          <div className="text-5xl font-bold text-emerald-400">98.5%</div>
          
          {/* Mock Radial */}
          <div className="absolute bottom-[-20px] right-[-20px] w-40 h-40 border-4 border-emerald-500/20 rounded-full flex items-center justify-center pointer-events-none">
             <motion.div 
               initial={{ rotate: -135 }}
               animate={{ rotate: 45 }}
               transition={{ duration: 1.5, ease: "easeOut" }}
               className="w-32 h-32 border-4 border-emerald-400 border-t-transparent rounded-full"
             ></motion.div>
          </div>
        </div>
      </div>

      {/* Main Graph (System Telemetry Throughput) */}
      <div className="glass-panel p-6 rounded-2xl mb-6 relative overflow-hidden h-72 flex flex-col">
         <div className="flex items-center justify-between mb-6 z-10">
            <span className="text-sm font-medium text-zinc-300">System Telemetry Throughput</span>
            <div className="flex items-center gap-2">
               <span className="w-2 h-2 rounded-full bg-cortex-400 animate-slow-pulse"></span>
               <span className="text-xs text-zinc-500">Live</span>
            </div>
         </div>
         {/* Simulated glowing wave graph */}
         <div className="flex-1 relative w-full h-full flex items-end">
           <svg viewBox="0 0 1000 200" preserveAspectRatio="none" className="absolute w-full h-full bottom-0">
             <defs>
               <linearGradient id="grad1" x1="0%" y1="0%" x2="0%" y2="100%">
                 <stop offset="0%" stopColor="rgba(168,85,247,0.3)" />
                 <stop offset="100%" stopColor="rgba(168,85,247,0)" />
               </linearGradient>
             </defs>
             <motion.path 
               initial={{ opacity: 0 }}
               animate={{ opacity: 1 }}
               transition={{ duration: 1 }}
               d="M0,150 Q100,50 250,120 T500,100 T750,150 T1000,80 L1000,200 L0,200 Z" 
               fill="url(#grad1)" 
             />
             <motion.path 
               initial={{ pathLength: 0 }}
               animate={{ pathLength: 1 }}
               transition={{ duration: 2, ease: "easeInOut" }}
               d="M0,150 Q100,50 250,120 T500,100 T750,150 T1000,80" 
               fill="none" 
               stroke="#a855f7" 
               strokeWidth="3" 
               className="drop-shadow-[0_0_10px_rgba(168,85,247,0.8)]" 
             />
             
             {/* Secondary line */}
             <motion.path 
               initial={{ pathLength: 0 }}
               animate={{ pathLength: 1 }}
               transition={{ duration: 2, ease: "easeInOut", delay: 0.5 }}
               d="M0,180 Q150,100 300,160 T600,120 T850,170 T1000,110" 
               fill="none" 
               stroke="#38bdf8" 
               strokeWidth="2" 
               strokeDasharray="5,5" 
               className="drop-shadow-[0_0_8px_rgba(56,189,248,0.5)] opacity-50" 
             />
           </svg>
         </div>
      </div>

      {/* Bottom Metrics Row */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
         {[
           { label: "Correlation Windows", val: "1050+" },
           { label: "Workflow Success Rate", val: "99.2%" },
           { label: "Remediation Success", val: "98.1%" },
           { label: "Mean Time To Detect", val: "4.8m" },
         ].map((m, i) => (
           <motion.div 
             key={i} 
             initial={{ opacity: 0, y: 10 }}
             animate={{ opacity: 1, y: 0 }}
             transition={{ delay: 0.5 + (i * 0.1) }}
             className="glass-panel p-4 rounded-xl flex flex-col justify-center"
           >
             <span className="text-zinc-400 text-xs mb-1">{m.label}</span>
             <span className="text-2xl font-bold text-white">{m.val}</span>
           </motion.div>
         ))}
      </div>

    </div>
  );
}
