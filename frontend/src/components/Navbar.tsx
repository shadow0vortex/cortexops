"use client";

import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Star } from "lucide-react";

export function Navbar() {
  const pathname = usePathname();

  const links = [
    { href: "/", label: "Platform" },
    { href: "/architecture", label: "Architecture" },
    { href: "/workflows", label: "Workflows" },
    { href: "/observability", label: "Observability" },
    { href: "/docs", label: "Docs" },
    { href: "https://github.com/shadow0vortex/cortexops", label: "GitHub", external: true },
  ];

  return (
    <div className="fixed top-6 left-0 right-0 z-50 flex justify-center px-4 pointer-events-none">
      <nav className="pointer-events-auto flex items-center justify-between px-6 py-3 rounded-[32px] glass-panel w-full max-w-6xl">
        
        {/* Left: Logo */}
        <div className="flex items-center">
          <div className="w-20 h-16 flex items-center justify-center mix-blend-screen overflow-hidden group hover:scale-110 transition-transform duration-500 origin-left">
            <Image src="/logo.png" alt="Logo" width={80} height={64} className="w-full h-full object-contain drop-shadow-[0_0_20px_rgba(168,85,247,0.6)]" />
          </div>
        </div>

        {/* Center: Links */}
        <ul className="hidden md:flex items-center gap-1 lg:gap-2">
          {links.map((link) => (
            <li key={link.label}>
              <Link 
                href={link.href}
                target={link.external ? "_blank" : "_self"}
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

        {/* Right: Buttons */}
        <div className="flex items-center gap-4">
          <Link
            href="https://github.com/shadow0vortex/cortexops"
            target="_blank"
            className="hidden sm:flex items-center gap-2 px-4 py-2 rounded-full border border-zinc-700/50 bg-zinc-900/50 hover:bg-zinc-800 hover:border-zinc-600 transition-all text-sm font-medium text-zinc-300 hover:text-white group"
          >
            <Star className="w-4 h-4 group-hover:text-cortex-400 transition-colors" />
            <span>Star on GitHub</span>
          </Link>
          <Link
            href="/deploy"
            className="flex items-center gap-2 px-5 py-2 rounded-full bg-white text-black hover:bg-zinc-200 transition-all text-sm font-semibold shadow-[0_0_20px_rgba(255,255,255,0.2)]"
          >
            Deploy CortexOps &rarr;
          </Link>
        </div>

      </nav>
    </div>
  );
}
