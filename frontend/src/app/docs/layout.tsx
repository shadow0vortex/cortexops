"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { ChevronDown, BookOpen } from "lucide-react";

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  // Close sidebar on route change (mobile)
  useEffect(() => {
    setSidebarOpen(false);
  }, [pathname]);

  const sidebarLinks = [
    { label: "Getting Started", href: "/docs/getting-started" },
    { label: "Installation", href: "/docs/installation" },
    { label: "Architecture", href: "/docs/architecture" },
    { label: "Platform Components", href: "/docs/platform-components" },
    { label: "Infrastructure", href: "/docs/infrastructure" },
    { label: "Operations Overview", href: "/docs/operations" },
    { label: "  ↳ Incident Response", href: "/docs/operations/incident-response" },
    { label: "  ↳ Remediation", href: "/docs/operations/remediation" },
    { label: "  ↳ Temporal Recovery", href: "/docs/operations/temporal-recovery" },
    { label: "  ↳ Chaos Operations", href: "/docs/operations/chaos" },
    { label: "  ↳ Backup & Restore", href: "/docs/operations/backup" },
    { label: "Security & Governance", href: "/docs/security-governance" },
    { label: "Reference", href: "/docs/reference" },
    { label: "Engineering", href: "/docs/engineering" },
    { label: "Deployment", href: "/docs/deployment" },
  ];

  const currentLabel = sidebarLinks.find((l) => l.href === pathname)?.label || "Documentation";

  const sidebarContent = (
    <>
      <h3 className="text-sm font-semibold text-zinc-500 uppercase tracking-widest mb-6">
        Documentation
      </h3>
      <ul className="flex flex-col gap-1" role="list">
        {sidebarLinks.map((link) => {
          const isActive = pathname === link.href;
          return (
            <li key={link.label}>
              <Link
                href={link.href}
                aria-current={isActive ? "page" : undefined}
                className={`block px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? "bg-cortex-500/10 text-cortex-300 border-l-2 border-cortex-400"
                    : "text-zinc-400 hover:bg-zinc-900 hover:text-zinc-200 border-l-2 border-transparent"
                }`}
              >
                {link.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </>
  );

  return (
    <div className="pt-32 pb-24 px-6 max-w-7xl mx-auto flex flex-col md:flex-row gap-12 min-h-screen relative">
      {/* Mobile sidebar toggle */}
      <div className="md:hidden">
        <button
          onClick={() => setSidebarOpen(!sidebarOpen)}
          className="flex items-center justify-between w-full px-4 py-3 rounded-xl glass-panel text-sm font-medium text-zinc-300 hover:text-white transition-colors"
          aria-expanded={sidebarOpen}
          aria-controls="docs-sidebar"
        >
          <div className="flex items-center gap-2">
            <BookOpen className="w-4 h-4 text-cortex-400" />
            <span>{currentLabel}</span>
          </div>
          <ChevronDown
            className={`w-4 h-4 transition-transform duration-200 ${sidebarOpen ? "rotate-180" : ""}`}
          />
        </button>

        <AnimatePresence>
          {sidebarOpen && (
            <motion.div
              id="docs-sidebar"
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.25, ease: [0.16, 1, 0.3, 1] }}
              className="overflow-hidden mt-3"
            >
              <div className="glass-panel rounded-xl p-4">{sidebarContent}</div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* Desktop sidebar */}
      <aside className="hidden md:block w-64 flex-shrink-0" role="navigation" aria-label="Documentation sidebar">
        {sidebarContent}
      </aside>

      {/* Main Content */}
      <main className="flex-1 glass-panel rounded-3xl p-8 md:p-12 prose prose-invert max-w-none">
        {children}
      </main>
    </div>
  );
}
