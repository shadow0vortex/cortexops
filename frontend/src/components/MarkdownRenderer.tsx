"use client";

import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { motion } from "framer-motion";
import { Terminal, ChevronRight, FileText } from "lucide-react";
import React from "react";

export function MarkdownRenderer({ content }: { content: string }) {
  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, ease: "easeOut" }}
      className="w-full"
    >
      <ReactMarkdown 
        remarkPlugins={[remarkGfm]}
        components={{
          h1: ({ node, className, ...props }) => (
            <h1 className={`text-4xl font-bold text-white mb-4 ${className || ""}`} {...props} />
          ),
          h2: ({ node, className, ...props }) => (
            <h2 className={`text-2xl font-semibold text-white mt-12 mb-4 ${className || ""}`} {...props} />
          ),
          h3: ({ node, className, ...props }) => <h3 className={`text-xl font-semibold text-white mt-8 mb-3 ${className || ""}`} {...props} />,
          h4: ({ node, className, ...props }) => <h4 className={`text-lg font-semibold text-white mt-6 mb-2 ${className || ""}`} {...props} />,
          p: ({ node, className, ...props }) => <p className={`text-sm text-zinc-400 mb-6 leading-relaxed ${className || ""}`} {...props} />,
          a: ({ node, className, ...props }) => <a className={`text-cortex-400 hover:text-cortex-300 transition-all duration-300 underline decoration-cortex-500/30 hover:decoration-cortex-400 underline-offset-4 ${className || ""}`} {...props} />,
          code: ({ node, className, ...props }) => {
            const isInline = !node?.position?.start.line || node.position.start.line === node.position.end.line;
            if (isInline) {
              return <code className={`text-cortex-300 font-mono bg-cortex-500/10 border border-cortex-500/20 px-1.5 py-0.5 rounded text-[13px] shadow-inner ${className || ""}`} {...props} />;
            }
            return <code className={`font-mono text-[13px] text-zinc-300 ${className || ""}`} {...props} />;
          },
          pre: ({ node, className, ...props }) => (
            <div className="relative group mb-8">
              <div className="glass-panel rounded-xl overflow-hidden relative border border-white/10 shadow-lg">
                <div className="flex items-center gap-2 px-4 py-2 bg-white/5 border-b border-white/5">
                  <div className="flex gap-2">
                    <div className="w-2.5 h-2.5 rounded-full bg-red-500/50 hover:bg-red-500 transition-colors" />
                    <div className="w-2.5 h-2.5 rounded-full bg-yellow-500/50 hover:bg-yellow-500 transition-colors" />
                    <div className="w-2.5 h-2.5 rounded-full bg-green-500/50 hover:bg-green-500 transition-colors" />
                  </div>
                  <Terminal className="w-3.5 h-3.5 text-zinc-500 ml-auto" />
                </div>
                <pre className={`p-4 overflow-x-auto text-[13px] text-zinc-300 font-mono bg-[#0B0B0F]/80 leading-relaxed ${className || ""}`} {...props} />
              </div>
            </div>
          ),
          ul: ({ node, className, ...props }) => <ul className={`list-disc pl-5 space-y-2 mb-6 marker:text-cortex-500 text-sm text-zinc-400 ${className || ""}`} {...props} />,
          ol: ({ node, className, ...props }) => <ol className={`list-decimal pl-5 space-y-2 mb-6 marker:text-cortex-500 text-sm text-zinc-400 ${className || ""}`} {...props} />,
          li: ({ node, className, ...props }) => <li className={`pl-1 ${className || ""}`} {...props} />,
          blockquote: ({ node, className, ...props }) => (
            <blockquote 
              className={`relative bg-cortex-500/5 border-l-2 border-cortex-500 p-4 rounded-r-xl mb-6 text-zinc-300 text-sm italic ${className || ""}`} 
              {...props} 
            />
          ),
          table: ({ node, className, ...props }) => (
            <div className="glass-panel rounded-xl overflow-hidden mb-6 shadow-md border border-white/10">
              <div className="overflow-x-auto">
                <table className={`w-full text-left text-sm text-zinc-300 ${className || ""}`} {...props} />
              </div>
            </div>
          ),
          thead: ({ node, className, ...props }) => <thead className={`bg-white/5 text-zinc-100 border-b border-white/10 ${className || ""}`} {...props} />,
          th: ({ node, className, ...props }) => <th className={`px-4 py-3 font-semibold text-white ${className || ""}`} {...props} />,
          td: ({ node, className, ...props }) => <td className={`px-4 py-3 border-b border-white/5 bg-black/20 ${className || ""}`} {...props} />,
          strong: ({ node, className, ...props }) => <strong className={`font-semibold text-white ${className || ""}`} {...props} />,
        }}
      >
        {content}
      </ReactMarkdown>
    </motion.div>
  );
}
