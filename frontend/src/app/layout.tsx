import type { Metadata } from "next";
import { Inter, Space_Grotesk } from "next/font/google";
import "./globals.css";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });
const spaceGrotesk = Space_Grotesk({ subsets: ["latin"], variable: "--font-space" });

export const metadata: Metadata = {
  metadataBase: new URL('https://cortexops.amshithnair.in'),
  title: "CortexOps | Deterministic Infrastructure Intelligence",
  description: "Topology-aware incident correlation, replay-safe workflows, and policy-governed remediation for Kubernetes.",
  keywords: ["Kubernetes", "Platform Engineering", "SRE", "Temporal", "NATS", "Infrastructure Automation", "Incident Management", "Topology Intelligence"],
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
  const jsonLd = {
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
  };

  return (
    <html lang="en" className="dark">
      <head>
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
      <body className={`${inter.variable} ${spaceGrotesk.variable} antialiased bg-black text-white selection:bg-cortex-500/30 selection:text-cortex-100`}>
        <Navbar />
        <main className="min-h-screen">
          {children}
        </main>
        <Footer />
      </body>
    </html>
  );
}
