"use client";

import { motion } from "framer-motion";

export function ReplayRecovery() {
  return (
    <div className="w-full bg-[#0B0B0F] border border-zinc-800 rounded-3xl p-8 my-8 relative">
      <div className="flex flex-col gap-8 max-w-3xl mx-auto">
        
        {/* Entities */}
        <div className="flex justify-between border-b border-zinc-800 pb-4">
          <div className="w-32 text-center font-mono text-sm text-cortex-400 font-bold">NATS JetStream</div>
          <div className="w-32 text-center font-mono text-sm text-white font-bold">Correlation Engine</div>
          <div className="w-32 text-center font-mono text-sm text-zinc-400 font-bold">Audit DB</div>
        </div>

        <div className="relative">
          {/* Vertical Lifelines */}
          <div className="absolute top-0 bottom-0 left-[4.2rem] sm:left-[5.5rem] w-px bg-zinc-800 border-dashed"></div>
          <div className="absolute top-0 bottom-0 left-1/2 w-px bg-zinc-800 border-dashed -translate-x-1/2"></div>
          <div className="absolute top-0 bottom-0 right-[4.2rem] sm:right-[5.5rem] w-px bg-zinc-800 border-dashed"></div>

          {/* Sequence 1: Standard Execution */}
          <div className="mb-6 relative z-10">
            <div className="text-xs text-zinc-500 mb-2 uppercase tracking-widest text-center font-semibold">Standard Execution</div>
            
            <motion.div initial={{ opacity: 0, x: -20 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true }} className="flex justify-between items-center w-1/2 ml-[5.5rem] mb-3">
              <div className="h-px bg-cortex-500 w-full relative">
                <div className="absolute right-0 top-1/2 -translate-y-1/2 border-[4px] border-transparent border-l-cortex-500"></div>
                <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] text-zinc-300 whitespace-nowrap bg-[#0B0B0F] px-2">Event A (id: 123)</div>
              </div>
            </motion.div>

            <motion.div initial={{ opacity: 0, x: -20 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true }} transition={{ delay: 0.2 }} className="flex justify-between items-center w-1/2 ml-auto mr-[5.5rem]">
              <div className="h-px bg-zinc-500 w-full relative">
                <div className="absolute right-0 top-1/2 -translate-y-1/2 border-[4px] border-transparent border-l-zinc-500"></div>
                <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] text-zinc-400 whitespace-nowrap bg-[#0B0B0F] px-2">Save Incident (id: 123)</div>
              </div>
            </motion.div>
          </div>

          {/* Sequence 2: Chaos Partition */}
          <div className="mb-6 relative z-10">
            <div className="text-xs text-red-500/80 mb-2 uppercase tracking-widest text-center font-semibold">Chaos Partition Occurs</div>
            <motion.div initial={{ opacity: 0 }} whileInView={{ opacity: 1 }} viewport={{ once: true }} transition={{ delay: 0.4 }} className="flex justify-center mb-4">
              <div className="bg-red-950/40 border border-red-900/50 text-red-400 text-[10px] px-3 py-1 rounded">Connection Dropped / Pod Crashes</div>
            </motion.div>
          </div>

          {/* Sequence 3: Recovery */}
          <div className="relative z-10">
            <div className="text-xs text-cortex-400 mb-2 uppercase tracking-widest text-center font-semibold">Recovery & Replay Phase</div>
            
            <motion.div initial={{ opacity: 0, x: 20 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true }} transition={{ delay: 0.6 }} className="flex justify-between items-center w-1/2 ml-[5.5rem] mb-4">
              <div className="h-px bg-zinc-600 w-full relative border-dashed border-b">
                <div className="absolute left-0 top-1/2 -translate-y-1/2 border-[4px] border-transparent border-r-zinc-600"></div>
                <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] text-zinc-400 whitespace-nowrap bg-[#0B0B0F] px-2">Reconnect</div>
              </div>
            </motion.div>

            <motion.div initial={{ opacity: 0, x: -20 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true }} transition={{ delay: 0.8 }} className="flex justify-between items-center w-1/2 ml-[5.5rem] mb-4">
              <div className="h-px bg-cortex-500 w-full relative">
                <div className="absolute right-0 top-1/2 -translate-y-1/2 border-[4px] border-transparent border-l-cortex-500"></div>
                <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] text-cortex-300 whitespace-nowrap bg-[#0B0B0F] px-2">Event A (id: 123) [Redelivered]</div>
              </div>
            </motion.div>

            <motion.div initial={{ opacity: 0, x: -20 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true }} transition={{ delay: 1.0 }} className="flex justify-between items-center w-1/2 ml-auto mr-[5.5rem] mb-4">
              <div className="h-px bg-zinc-500 w-full relative">
                <div className="absolute right-0 top-1/2 -translate-y-1/2 border-[4px] border-transparent border-l-zinc-500"></div>
                <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] text-zinc-400 whitespace-nowrap bg-[#0B0B0F] px-2">Query Exists? (id: 123)</div>
              </div>
            </motion.div>

            <motion.div initial={{ opacity: 0, x: 20 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true }} transition={{ delay: 1.2 }} className="flex justify-between items-center w-1/2 ml-auto mr-[5.5rem] mb-4">
              <div className="h-px bg-zinc-500 w-full relative border-dashed border-b">
                <div className="absolute left-0 top-1/2 -translate-y-1/2 border-[4px] border-transparent border-r-zinc-500"></div>
                <div className="absolute -top-5 left-1/2 -translate-x-1/2 text-[10px] text-zinc-400 whitespace-nowrap bg-[#0B0B0F] px-2">True</div>
              </div>
            </motion.div>

            <motion.div initial={{ opacity: 0 }} whileInView={{ opacity: 1 }} viewport={{ once: true }} transition={{ delay: 1.4 }} className="flex justify-center">
              <div className="bg-cortex-900/30 border border-cortex-500/50 text-cortex-300 text-[10px] px-3 py-1 rounded">Drop Duplicate deterministically</div>
            </motion.div>

          </div>
        </div>
      </div>
    </div>
  );
}
