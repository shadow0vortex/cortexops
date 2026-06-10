import { Hero } from "@/components/Hero";
import { StatusStrip } from "@/components/StatusStrip";
import { FeatureCarousel } from "@/components/FeatureCarousel";
import { PlatformIntro } from "@/components/PlatformIntro";
import dynamic from "next/dynamic";
const ArchitecturePreview = dynamic(() => import("@/components/ArchitecturePreview").then(mod => mod.ArchitecturePreview));
import { OperationalGuarantees } from "@/components/OperationalGuarantees";
import { PlatformCta } from "@/components/PlatformCta";

export default function Home() {
  return (
    <div className="flex flex-col relative overflow-hidden">
      <Hero />
      <StatusStrip />
      <div className="relative z-20 space-y-12">
        <PlatformIntro />
        <FeatureCarousel />
        <ArchitecturePreview />
        <OperationalGuarantees />
        <PlatformCta />
      </div>
    </div>
  );
}
