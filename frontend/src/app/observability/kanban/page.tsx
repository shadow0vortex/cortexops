import { Board } from "@/components/kanban/Board";
import { Workflow, ArrowLeft } from "lucide-react";
import Link from "next/link";

export default function KanbanPage() {
  return (
    <div className="pt-32 pb-24 px-6 max-w-screen-2xl mx-auto flex flex-col min-h-screen relative overflow-hidden">
      
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1200px] h-[600px] bg-[radial-gradient(ellipse_at_top,rgba(168,85,247,0.1),transparent_50%)] pointer-events-none"></div>

      <div className="relative z-10 flex flex-col h-full">
        <div className="mb-8">
          <Link href="/observability" className="flex items-center gap-2 text-zinc-400 hover:text-white transition-colors mb-4 w-fit">
            <ArrowLeft className="w-4 h-4" />
            <span className="text-sm font-medium">Back to Observability</span>
          </Link>
          
          <div className="flex items-center gap-3 mb-2">
            <div className="p-2 rounded-lg bg-cortex-500/20 text-cortex-400 border border-cortex-500/30">
              <Workflow className="w-5 h-5" />
            </div>
            <h1 className="text-3xl md:text-4xl font-bold text-white">Incident Orchestration Board</h1>
          </div>
          <p className="text-zinc-400">Drag and drop incidents through the Temporal workflow states.</p>
        </div>

        <Board />
      </div>
    </div>
  );
}
