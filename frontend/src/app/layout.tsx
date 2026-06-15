import type { Metadata, Viewport } from "next";
import { Inter, Space_Grotesk } from "next/font/google";
import "./globals.css";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { PostHogProvider } from "@/components/PostHogProvider";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });
const spaceGrotesk = Space_Grotesk({ subsets: ["latin"], variable: "--font-space" });

export const viewport: Viewport = {
  themeColor: "#000000",
  width: "device-width",
  initialScale: 1,
};

export const metadata: Metadata = {
  metadataBase: new URL('https://cortexops.amshithnair.in'),
  title: "CortexOps | Deterministic Infrastructure Intelligence",
  description: "Topology-aware incident correlation, replay-safe workflows, and policy-governed remediation for Kubernetes.",
  keywords: ["Kubernetes", "Platform Engineering", "SRE", "Temporal", "NATS", "Infrastructure Automation", "Incident Management", "Topology Intelligence"],
  icons: {
    icon: [
      { url: "/favicon-16x16.png", sizes: "16x16", type: "image/png" },
      { url: "/favicon-32x32.png", sizes: "32x32", type: "image/png" },
    ],
    apple: [
      { url: "/apple-touch-icon.png", sizes: "180x180", type: "image/png" },
    ],
  },
  openGraph: {
    title: "CortexOps | Deterministic Infrastructure Intelligence",
    description: "Topology-aware incident correlation, replay-safe workflows, and policy-governed remediation for Kubernetes.",
    url: "https://cortexops.amshithnair.in",
    siteName: "CortexOps",
    images: [
      {
        url: "/og-image.png",
        width: 1200,
        height: 630,
      },
    ],
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "CortexOps | Deterministic Infrastructure Intelligence",
    description: "Topology-aware incident correlation, replay-safe workflows, and policy-governed remediation for Kubernetes.",
    images: ["/og-image.png"],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const jsonLd = [
    {
      "@context": "https://schema.org",
      "@type": "SoftwareApplication",
      "name": "CortexOps",
      "description": "Topology-aware incident correlation, replay-safe workflows, and policy-governed remediation for Kubernetes.",
      "applicationCategory": "DeveloperApplication",
      "operatingSystem": "Kubernetes",
      "url": "https://cortexops.amshithnair.in",
      "offers": {
        "@type": "Offer",
        "price": "0",
        "priceCurrency": "USD"
      }
    },
    {
      "@context": "https://schema.org",
      "@type": "Organization",
      "name": "CortexOps",
      "url": "https://cortexops.amshithnair.in",
      "logo": "https://cortexops.amshithnair.in/logo.png",
      "sameAs": [
        "https://github.com/shadow0vortex/cortexops"
      ]
    },
    {
      "@context": "https://schema.org",
      "@type": "WebSite",
      "name": "CortexOps",
      "url": "https://cortexops.amshithnair.in",
      "potentialAction": {
        "@type": "SearchAction",
        "target": "https://cortexops.amshithnair.in/docs?q={search_term_string}",
        "query-input": "required name=search_term_string"
      }
    }
  ];

  return (
    <html lang="en" className="dark">
      <head>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
      <body className={`${inter.variable} ${spaceGrotesk.variable} antialiased bg-black text-white selection:bg-cortex-500/30 selection:text-cortex-100`}>
        <PostHogProvider>
          {/* Skip to content link for accessibility */}
          <a
            href="#main-content"
            className="sr-only focus:not-sr-only focus:fixed focus:top-4 focus:left-4 focus:z-[100] focus:px-4 focus:py-2 focus:bg-cortex-600 focus:text-white focus:rounded-lg focus:outline-none focus:ring-2 focus:ring-cortex-400 focus:ring-offset-2 focus:ring-offset-black"
          >
            Skip to main content
          </a>
          <Navbar />
          <main id="main-content" className="min-h-screen">
            {children}
          </main>
          <Footer />
        </PostHogProvider>
      </body>
    </html>
  );
}
