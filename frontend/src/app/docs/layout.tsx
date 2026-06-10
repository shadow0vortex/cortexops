"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  const sidebarLinks = [
    { label: "Getting Started", href: "/docs/getting-started" },
    { label: "Installation", href: "/docs/installation" },
    { label: "Architecture", href: "/docs/architecture" },
    { label: "Platform Components", href: "/docs/platform-components" },
    { label: "Infrastructure", href: "/docs/infrastructure" },
    { label: "Operations", href: "/docs/operations" },
    { label: "Security & Governance", href: "/docs/security-governance" },
    { label: "Reference", href: "/docs/reference" },
    { label: "Engineering", href: "/docs/engineering" },
  ];

  return (
    <div className="pt-32 pb-24 px-6 max-w-7xl mx-auto flex flex-col md:flex-row gap-12 min-h-screen relative">
      
      {/* Sidebar */}
      <aside className="w-full md:w-64 flex-shrink-0">
        <h3 className="text-sm font-semibold text-zinc-500 uppercase tracking-widest mb-6">Documentation</h3>
        <ul className="flex flex-col gap-2">
          {sidebarLinks.map(link => {
            const isActive = pathname === link.href;
            return (
              <li key={link.label}>
                <Link 
                  href={link.href} 
                  className={`block px-4 py-2 rounded-lg text-sm font-medium transition-colors ${isActive ? 'bg-cortex-500/10 text-cortex-300 border-l-2 border-cortex-400' : 'text-zinc-400 hover:bg-zinc-900 hover:text-zinc-200 border-l-2 border-transparent'}`}
                >
                  {link.label}
                </Link>
              </li>
            );
          })}
        </ul>
      </aside>

      {/* Main Content */}
      <main className="flex-1 glass-panel rounded-3xl p-8 md:p-12 prose prose-invert max-w-none">
        {children}
      </main>
    </div>
  );
}
