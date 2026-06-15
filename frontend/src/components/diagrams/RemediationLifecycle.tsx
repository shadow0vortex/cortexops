"use client";

import { motion } from "framer-motion";
import { CheckCircle2, ShieldAlert, Timer, UserCheck, Play, XCircle, ActivitySquare } from "lucide-react";

export function RemediationLifecycle() {
  const steps = [
    { id: "PROPOSED", icon: Timer, status: "pending" },
    { id: "POLICY EVALUATING", icon: ShieldAlert, status: "active" },
    { id: "APPROVAL PENDING", icon: UserCheck, status: "pending" },
    { id: "DRY RUNNING", icon: Play, status: "active" },
    { id: "EXECUTING", icon: Play, status: "active" },
    { id: "VERIFYING", icon: ActivitySquare, status: "active" },
    { id: "SUCCESS", icon: CheckCircle2, status: "success" },
  ];

  return (
    <div className="w-full bg-[#0B0B0F] border border-zinc-800 rounded-3xl p-8 my-8 overflow-x-auto">
      <div className="min-w-[800px] flex flex-col items-center justify-center py-12 relative">
        
        {/* Main Flow Line */}
        <div className="absolute top-1/2 left-10 right-10 h-1 bg-zinc-800 -translate-y-1/2 z-0"></div>
        <div className="absolute top-1/2 left-10 right-1/2 h-1 bg-cortex-500 -translate-y-1/2 z-0 shadow-[0_0_10px_rgba(168,85,247,0.5)]"></div>

        <div className="flex justify-between w-full relative z-10 px-4">
          {steps.map((step, i) => {
            const Icon = step.icon;
            return (
              <motion.div 
                key={step.id}
                initial={{ opacity: 0, scale: 0.8 }}
                whileInView={{ opacity: 1, scale: 1 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.15 }}
                className="flex flex-col items-center gap-4 relative group"
              >
                <div className={`w-12 h-12 rounded-full flex items-center justify-center border-2 bg-black transition-colors ${
                  step.status === 'success' ? 'border-green-500 text-green-500 shadow-[0_0_15px_rgba(34,197,94,0.3)]' :
                  step.status === 'active' ? 'border-cortex-500 text-cortex-400 shadow-[0_0_15px_rgba(168,85,247,0.4)]' :
                  'border-zinc-700 text-zinc-500'
                }`}>
                  <Icon className="w-5 h-5" />
                </div>
                <div className="text-[10px] font-bold text-zinc-400 tracking-wider text-center w-24 uppercase">
                  {step.id}
                </div>
                
                {/* Failure Branch Examples */}
                {(step.id === "POLICY EVALUATING" || step.id === "VERIFYING") && (
                  <motion.div 
                    initial={{ opacity: 0 }}
                    whileInView={{ opacity: 1 }}
                    transition={{ delay: i * 0.15 + 0.5 }}
                    className="absolute top-20 flex flex-col items-center"
                  >
                    <div className="w-0.5 h-8 bg-red-900/50 mb-2"></div>
                    <div className="flex items-center gap-1 text-[10px] text-red-400 bg-red-950/30 px-2 py-1 rounded border border-red-900/50 whitespace-nowrap">
                      <XCircle className="w-3 h-3" /> {step.id === "POLICY EVALUATING" ? "OPA Denied" : "Rollback"}
                    </div>
                  </motion.div>
                )}
              </motion.div>
            )
          })}
        </div>
      </div>
    </div>
  );
}
