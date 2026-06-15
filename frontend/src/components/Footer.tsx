import Link from "next/link";
import Image from "next/image";

export function Footer() {
  const sections = [
    {
      title: "Platform",
      links: [
        { label: "Architecture", href: "/architecture" },
        { label: "Platform Overview", href: "/platform" },
        { label: "Feature Matrix", href: "/resources/feature-matrix" },
        { label: "Workflows", href: "/workflows" },
      ]
    },
    {
      title: "Resources",
      links: [
        { label: "Documentation", href: "/docs" },
        { label: "Production Readiness", href: "/resources/production-readiness" },
        { label: "Deployment Guide", href: "/docs/deployment" },
        { label: "Security", href: "/docs/security-governance" },
      ]
    },
    {
      title: "Community",
      links: [
        { label: "GitHub", href: "https://github.com/shadow0vortex/cortexops" },
        { label: "LinkedIn", href: "https://linkedin.com" },
        { label: "Engineering Decisions", href: "/docs/engineering" },
        { label: "Load Test Results", href: "/resources/load-testing" },
      ]
    }
  ];

  return (
    <footer className="border-t border-zinc-800/50 bg-black pt-16 pb-12 px-6" role="contentinfo">
      <div className="max-w-7xl mx-auto flex flex-col md:flex-row justify-between gap-12">
        
        {/* Left Branding */}
        <div className="flex flex-col gap-4 max-w-xs">
          <div className="flex items-center gap-2">
            <div className="w-10 h-8 flex items-center justify-center mix-blend-screen overflow-hidden">
              <Image src="/logo.png" alt="CortexOps Logo" width={40} height={32} className="w-full h-full object-contain drop-shadow-[0_0_10px_rgba(168,85,247,0.5)] scale-125" />
            </div>
            <span className="text-white font-semibold tracking-tight">CortexOps</span>
          </div>
          <p className="text-zinc-500 text-sm">
            © CortexOps {new Date().getFullYear()}<br/>
            Deterministic Infrastructure Intelligence.
          </p>
        </div>

        {/* Links */}
        <div className="flex flex-wrap gap-8 sm:gap-16">
          {sections.map(section => (
            <div key={section.title} className="flex flex-col gap-4">
              <h4 className="text-white font-semibold text-sm">{section.title}</h4>
              <ul className="flex flex-col gap-3">
                {section.links.map(link => (
                  <li key={link.label}>
                    <Link href={link.href} className="text-zinc-400 hover:text-cortex-300 text-sm transition-colors">
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
          
          {/* GitHub CTA */}
          <div className="flex flex-col gap-4">
            <h4 className="text-white font-semibold text-sm">Star us on GitHub</h4>
            <p className="text-zinc-500 text-sm max-w-[180px]">Help us build the future of infrastructure operations.</p>
            <Link 
              href="https://github.com/shadow0vortex/cortexops"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="Star CortexOps on GitHub"
              className="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-zinc-900 border border-zinc-800 hover:border-zinc-700 hover:bg-zinc-800 transition-all text-sm text-white mt-2 w-max"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/><path d="M9 18c-4.51 2-5-2-7-2"/></svg>
              Star on GitHub
            </Link>
          </div>
        </div>

      </div>
    </footer>
  );
}
