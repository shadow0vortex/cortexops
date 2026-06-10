"use client";

import { motion } from "framer-motion";
import { Server, Database, GitMerge, Network, Activity, Workflow, ShieldAlert, BrainCircuit, Cpu } from "lucide-react";

import React from 'react';

export function SystemArchitecture() {
  const Node = ({ title, subtitle, icon: Icon, delay = 0, color = "cortex" }: { title: string, subtitle: string, icon: React.ElementType, delay?: number, color?: string }) => (
    <motion.div 
      initial={{ opacity: 0, translateZ: 200 }}
      whileInView={{ opacity: 1, translateZ: 0 }}
      viewport={{ once: true }}
      transition={{ delay, duration: 0.8, type: "spring" }}
      style={{ transformStyle: 'preserve-3d' }}
      className="relative w-40 h-28 group z-20 cursor-pointer"
    >
      <div className={`absolute inset-0 border ${color === 'cortex' ? 'border-cortex-500/50 bg-zinc-900/90' : 'border-zinc-700 bg-zinc-950/90'} flex flex-col items-center justify-center gap-2 z-10 transform-gpu transition-all duration-300 group-hover:-translate-y-2 group-hover:-translate-x-2 shadow-[inset_0_0_20px_rgba(168,85,247,0.1)] group-hover:shadow-[inset_0_0_30px_rgba(168,85,247,0.4)] backdrop-blur-md`}>
        <div className={`p-2 rounded-lg ${color === 'cortex' ? 'bg-cortex-500/20 text-cortex-400' : 'bg-zinc-800 text-zinc-400'}`}>
          <Icon className="w-6 h-6" />
        </div>
        <div className="text-center">
          <div className="text-sm font-bold text-white tracking-wide">{title}</div>
          <div className="text-[10px] text-zinc-400 uppercase">{subtitle}</div>
        </div>
      </div>
      {/* 3D Faces */}
      <div className={`absolute top-0 right-0 w-4 h-full ${color === 'cortex' ? 'bg-zinc-950 border-cortex-500/50' : 'bg-black border-zinc-700'} border-y border-r origin-right brightness-50 group-hover:brightness-75 transition-all`} style={{ transform: 'translateX(100%) rotateY(90deg)' }}></div>
      <div className={`absolute bottom-0 left-0 w-full h-4 ${color === 'cortex' ? 'bg-zinc-900 border-cortex-500/50' : 'bg-zinc-950 border-zinc-700'} border-x border-b origin-bottom brightness-75 group-hover:brightness-100 transition-all`} style={{ transform: 'translateY(100%) rotateX(-90deg)' }}></div>
    </motion.div>
  );

  return (
    <div className="w-full bg-[#0B0B0F] border border-zinc-800 rounded-3xl p-8 my-8 relative overflow-hidden flex flex-col items-center perspective-[2500px]">
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,rgba(168,85,247,0.05),transparent_70%)]"></div>
      
      <motion.div 
        initial={{ rotateX: 55, rotateZ: -45 }}
        whileHover={{ rotateZ: -40 }}
        transition={{ duration: 1, ease: "easeOut" }}
        style={{ transformStyle: 'preserve-3d' }}
        className="w-full max-w-4xl py-24 flex items-center justify-center"
      >
        <div className="grid grid-cols-1 md:grid-cols-3 gap-x-20 gap-y-24 relative place-items-center w-full" style={{ transformStyle: 'preserve-3d' }}>
          
          {/* Base Floor Plane */}
          <div className="absolute inset-[-20%] bg-cortex-500/5 border border-cortex-500/20 rounded-3xl" style={{ transform: 'translateZ(-20px)' }}></div>
          <div className="absolute inset-[-20%] bg-[linear-gradient(to_right,rgba(168,85,247,0.1)_1px,transparent_1px),linear-gradient(to_bottom,rgba(168,85,247,0.1)_1px,transparent_1px)] bg-[size:40px_40px]" style={{ transform: 'translateZ(-19px)' }}></div>

          {/* Row 1 */}
          <div className="md:col-start-2" style={{ transformStyle: 'preserve-3d' }}>
            <Node title="K8s API Server" subtitle="Watch Events" icon={Server} color="zinc" delay={0.1} />
          </div>

          {/* Row 2 */}
          <div className="md:col-start-2 relative" style={{ transformStyle: 'preserve-3d' }}>
            <div className="absolute -top-24 left-1/2 w-1 h-24 bg-gradient-to-b from-zinc-700 to-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
            <Node title="Collector" subtitle="Normalize" icon={Cpu} delay={0.2} />
          </div>

          {/* Row 3 */}
          <div className="md:col-start-2 relative" style={{ transformStyle: 'preserve-3d' }}>
            <div className="absolute -top-24 left-1/2 w-1 h-24 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
            <Node title="NATS JetStream" subtitle="Durable Events" icon={Database} delay={0.3} />
          </div>

          {/* Row 4 */}
          <div className="md:col-start-2 relative" style={{ transformStyle: 'preserve-3d' }}>
            <div className="absolute -top-24 left-1/2 w-1 h-24 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
            <Node title="Correlation Engine" subtitle="Incident Grouping" icon={GitMerge} delay={0.4} />
          </div>

          {/* Row 5 - Split */}
          <div className="md:col-span-3 w-full flex justify-between relative px-4" style={{ transformStyle: 'preserve-3d' }}>
            {/* Connector lines on the floor */}
            <div className="absolute -top-24 left-1/2 w-1 h-12 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
            <div className="absolute -top-12 left-[15%] right-[15%] h-1 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
            <div className="absolute -top-12 left-[15%] w-1 h-12 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
            <div className="absolute -top-12 right-[15%] w-1 h-12 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
            <div className="absolute -top-12 left-1/2 w-1 h-12 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>

            <Node title="Topology Graph" subtitle="Blast Radius" icon={Network} delay={0.5} />
            <Node title="RCA Engine" subtitle="Advisory Report" icon={BrainCircuit} delay={0.5} />
            <Node title="Remediation" subtitle="Orchestration" icon={Activity} delay={0.5} />
          </div>

          {/* Row 6 - Remediation Dependencies */}
          <div className="md:col-start-3 w-full flex justify-between gap-8 relative" style={{ transformStyle: 'preserve-3d' }}>
             {/* Connectors from Remediation */}
             <div className="absolute -top-24 left-[20%] w-1 h-24 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
             <div className="absolute -top-24 right-[20%] w-1 h-24 bg-cortex-500/50 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>
             
             <Node title="Temporal" subtitle="Durable Execution" icon={Workflow} color="zinc" delay={0.6} />
             <Node title="OPA Engine" subtitle="Policy Gates" icon={ShieldAlert} color="zinc" delay={0.6} />
          </div>

        </div>
      </motion.div>
    </div>
  );
}
