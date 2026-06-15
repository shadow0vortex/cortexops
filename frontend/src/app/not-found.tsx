"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { ChevronRight, Radio } from "lucide-react";

export default function NotFound() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center px-6 relative overflow-hidden">
      {/* Background effects */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-cortex-600/10 rounded-full blur-[160px]"></div>
        <div className="absolute top-1/4 right-1/4 w-[200px] h-[200px] bg-cortex-500/5 rounded-full blur-[80px] animate-slow-pulse"></div>
        <div className="absolute bottom-1/3 left-1/3 w-[300px] h-[300px] bg-cortex-700/5 rounded-full blur-[100px] animate-slow-pulse"></div>
      </div>

      {/* Scanning grid lines */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808008_1px,transparent_1px),linear-gradient(to_bottom,#80808008_1px,transparent_1px)] bg-[size:40px_40px] pointer-events-none"></div>

      <motion.div
        initial={{ opacity: 0, y: 30 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
        className="relative z-10 flex flex-col items-center text-center max-w-xl"
      >
        {/* 404 number */}
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 1, ease: [0.16, 1, 0.3, 1] }}
          className="relative mb-8"
        >
          <span className="text-[10rem] md:text-[14rem] font-bold leading-none tracking-tighter text-transparent bg-clip-text bg-gradient-to-b from-zinc-800 to-zinc-900 select-none">
            404
          </span>
          <span className="absolute inset-0 flex items-center justify-center text-[10rem] md:text-[14rem] font-bold leading-none tracking-tighter text-transparent bg-clip-text bg-gradient-to-b from-cortex-400/60 via-cortex-500/40 to-transparent select-none blur-[1px]">
            404
          </span>
        </motion.div>

        {/* Status indicator */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.6 }}
          className="flex items-center gap-2 px-4 py-2 rounded-full border border-red-500/30 bg-red-500/10 text-red-400 text-sm font-medium mb-6 backdrop-blur-md"
        >
          <Radio className="w-4 h-4 animate-pulse" />
          Signal Lost
        </motion.div>

        {/* Message */}
        <motion.h1
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3, duration: 0.6 }}
          className="text-2xl md:text-3xl font-bold text-white mb-4"
        >
          The requested route could not be correlated.
        </motion.h1>
        <motion.p
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4, duration: 0.6 }}
          className="text-zinc-400 text-lg mb-10 leading-relaxed"
        >
          This endpoint does not exist in the CortexOps topology graph.
          The signal may have been deprecated or the route was never registered.
        </motion.p>

        {/* CTAs */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5, duration: 0.6 }}
          className="flex flex-col sm:flex-row items-center gap-4"
        >
          <Link
            href="/"
            className="flex items-center gap-2 px-6 py-3 rounded-full bg-white text-black font-semibold hover:bg-zinc-200 hover:scale-105 active:scale-95 transition-all duration-300 shadow-[0_0_20px_rgba(255,255,255,0.2)]"
          >
            Return Home <ChevronRight className="w-4 h-4" />
          </Link>
          <Link
            href="/docs"
            className="flex items-center gap-2 px-6 py-3 rounded-full border border-zinc-700 bg-zinc-900/50 text-white font-medium hover:bg-zinc-800 hover:border-zinc-500 transition-all duration-300"
          >
            View Documentation
          </Link>
        </motion.div>
      </motion.div>
    </div>
  );
}
