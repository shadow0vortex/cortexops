"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { ChevronRight } from "lucide-react";

export function PlatformCta() {
  return (
    <section className="py-32 px-6 max-w-4xl mx-auto text-center">
      <motion.div 
        initial={{ opacity: 0, scale: 0.95 }}
        whileInView={{ opacity: 1, scale: 1 }}
        viewport={{ once: true }}
        transition={{ duration: 0.8 }}
        className="glass-panel p-12 md:p-16 rounded-3xl border border-cortex-500/30 relative overflow-hidden"
      >
        <div className="absolute inset-0 bg-gradient-to-t from-cortex-600/20 to-transparent"></div>
        <div className="relative z-10">
          <h2 className="text-4xl md:text-5xl font-bold text-white mb-6">Operate Infrastructure With Confidence</h2>
          <p className="text-xl text-zinc-300 mb-10 max-w-2xl mx-auto">
            Move from reactive incident response to deterministic infrastructure operations.
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link 
              href="/docs"
              className="flex items-center gap-2 px-8 py-4 rounded-full bg-white text-black font-semibold hover:bg-zinc-200 transition-colors shadow-[0_0_20px_rgba(255,255,255,0.2)]"
            >
              Deploy CortexOps <ChevronRight className="w-4 h-4" />
            </Link>
            <Link 
              href="/architecture"
              className="flex items-center gap-2 px-8 py-4 rounded-full border border-zinc-700 bg-zinc-900/80 text-white font-medium hover:bg-zinc-800 transition-colors"
            >
              View Architecture
            </Link>
          </div>
        </div>
      </motion.div>
    </section>
  );
}
