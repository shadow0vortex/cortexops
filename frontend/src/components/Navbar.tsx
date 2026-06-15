"use client";

import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Star, Menu, X } from "lucide-react";
import { useState, useEffect, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";

export function Navbar() {
  const pathname = usePathname();
  const [mobileOpen, setMobileOpen] = useState(false);

  const links = [
    { href: "/", label: "Platform" },
    { href: "/architecture", label: "Architecture" },
    { href: "/workflows", label: "Workflows" },
    { href: "/observability", label: "Observability" },
    { href: "/docs", label: "Docs" },
    { href: "https://github.com/shadow0vortex/cortexops", label: "GitHub", external: true },
  ];

  // Close on route change
  useEffect(() => {
    setMobileOpen(false);
  }, [pathname]);

  // Close on Escape key
  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    if (e.key === "Escape") setMobileOpen(false);
  }, []);

  useEffect(() => {
    if (mobileOpen) {
      document.addEventListener("keydown", handleKeyDown);
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => {
      document.removeEventListener("keydown", handleKeyDown);
      document.body.style.overflow = "";
    };
  }, [mobileOpen, handleKeyDown]);

  return (
    <>
      <div className="fixed top-6 left-0 right-0 z-50 flex justify-center px-4 pointer-events-none">
        <nav
          className="pointer-events-auto flex items-center justify-between px-6 py-3 rounded-[32px] glass-panel w-full max-w-6xl"
          role="navigation"
          aria-label="Main navigation"
        >
          {/* Left: Logo */}
          <div className="flex items-center">
            <Link href="/" aria-label="CortexOps home">
              <div className="w-20 h-16 flex items-center justify-center mix-blend-screen overflow-hidden group hover:scale-110 transition-transform duration-500 origin-left">
                <Image src="/logo.png" alt="CortexOps Logo" width={80} height={64} className="w-full h-full object-contain drop-shadow-[0_0_20px_rgba(168,85,247,0.6)]" />
              </div>
            </Link>
          </div>

          {/* Center: Links (desktop) */}
          <ul className="hidden md:flex items-center gap-1 lg:gap-2">
            {links.map((link) => (
              <li key={link.label}>
                <Link 
                  href={link.href}
                  target={link.external ? "_blank" : "_self"}
                  rel={link.external ? "noopener noreferrer" : undefined}
                  className={cn(
                    "px-4 py-2 rounded-full text-sm font-medium transition-all duration-300",
                    pathname === link.href 
                      ? "text-white bg-white/10" 
                      : "text-zinc-400 hover:text-white hover:bg-white/5"
                  )}
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>

          {/* Right: Buttons (desktop) + Hamburger (mobile) */}
          <div className="flex items-center gap-4">
            <Link
              href="https://github.com/shadow0vortex/cortexops"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="Star CortexOps on GitHub"
              className="hidden sm:flex items-center gap-2 px-4 py-2 rounded-full border border-zinc-700/50 bg-zinc-900/50 hover:bg-zinc-800 hover:border-zinc-600 transition-all text-sm font-medium text-zinc-300 hover:text-white group"
            >
              <Star className="w-4 h-4 group-hover:text-cortex-400 transition-colors" />
              <span>Star on GitHub</span>
            </Link>
            <Link
              href="/deploy"
              className="hidden md:flex items-center gap-2 px-5 py-2 rounded-full bg-white text-black hover:bg-zinc-200 transition-all text-sm font-semibold shadow-[0_0_20px_rgba(255,255,255,0.2)]"
            >
              Deploy CortexOps &rarr;
            </Link>

            {/* Mobile hamburger */}
            <button
              onClick={() => setMobileOpen(!mobileOpen)}
              className="md:hidden flex items-center justify-center w-10 h-10 rounded-full text-zinc-300 hover:text-white hover:bg-white/10 transition-all"
              aria-label={mobileOpen ? "Close menu" : "Open menu"}
              aria-expanded={mobileOpen}
            >
              {mobileOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
          </div>
        </nav>
      </div>

      {/* Mobile overlay */}
      <AnimatePresence>
        {mobileOpen && (
          <>
            {/* Backdrop */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm md:hidden"
              onClick={() => setMobileOpen(false)}
              aria-hidden="true"
            />

            {/* Mobile menu panel */}
            <motion.div
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3, ease: [0.16, 1, 0.3, 1] }}
              className="fixed top-[100px] left-4 right-4 z-50 md:hidden"
              role="dialog"
              aria-modal="true"
              aria-label="Mobile navigation menu"
            >
              <div className="glass-panel rounded-2xl p-6 shadow-2xl border border-zinc-800">
                <ul className="flex flex-col gap-2">
                  {links.map((link) => (
                    <li key={link.label}>
                      <Link
                        href={link.href}
                        target={link.external ? "_blank" : "_self"}
                        rel={link.external ? "noopener noreferrer" : undefined}
                        onClick={() => setMobileOpen(false)}
                        className={cn(
                          "block px-4 py-3 rounded-xl text-base font-medium transition-all duration-200",
                          pathname === link.href
                            ? "text-white bg-cortex-500/10 border-l-2 border-cortex-400"
                            : "text-zinc-400 hover:text-white hover:bg-white/5"
                        )}
                      >
                        {link.label}
                      </Link>
                    </li>
                  ))}
                </ul>

                <div className="mt-6 pt-6 border-t border-zinc-800 flex flex-col gap-3">
                  <Link
                    href="/deploy"
                    onClick={() => setMobileOpen(false)}
                    className="flex items-center justify-center gap-2 px-5 py-3 rounded-full bg-white text-black font-semibold text-sm shadow-[0_0_20px_rgba(255,255,255,0.2)]"
                  >
                    Deploy CortexOps &rarr;
                  </Link>
                  <Link
                    href="https://github.com/shadow0vortex/cortexops"
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={() => setMobileOpen(false)}
                    className="flex items-center justify-center gap-2 px-5 py-3 rounded-full border border-zinc-700/50 bg-zinc-900/50 text-zinc-300 font-medium text-sm"
                  >
                    <Star className="w-4 h-4" />
                    Star on GitHub
                  </Link>
                </div>
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </>
  );
}
