import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { AlertCircle, Clock } from "lucide-react";
import { cn } from "@/lib/utils";

export type Incident = {
  id: string;
  title: string;
  service: string;
  severity: "critical" | "high" | "medium";
  time: string;
};

export function IncidentCard({ incident }: { incident: Incident }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: incident.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    zIndex: isDragging ? 50 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={cn(
        "glass-panel border rounded-xl p-4 cursor-grab active:cursor-grabbing mb-3 bg-zinc-900/80 backdrop-blur-md transition-shadow relative overflow-hidden",
        isDragging ? "opacity-50 shadow-[0_0_20px_rgba(168,85,247,0.4)] border-cortex-500 scale-105" : "border-zinc-700/50 hover:border-zinc-500 hover:bg-zinc-800/80"
      )}
    >
      <div className="absolute top-0 left-0 w-1 h-full" style={{ backgroundColor: incident.severity === 'critical' ? '#ef4444' : incident.severity === 'high' ? '#f97316' : '#eab308' }} />
      <div className="flex items-start justify-between mb-2">
        <h4 className="text-sm font-semibold text-white truncate pr-2">{incident.title}</h4>
        <AlertCircle className={cn("w-4 h-4 shrink-0", incident.severity === 'critical' ? 'text-red-500' : incident.severity === 'high' ? 'text-orange-500' : 'text-yellow-500')} />
      </div>
      <div className="text-xs text-zinc-400 mb-3">{incident.service}</div>
      <div className="flex items-center gap-1 text-[10px] text-zinc-500 font-medium bg-zinc-950/50 px-2 py-1 rounded w-fit">
        <Clock className="w-3 h-3" />
        {incident.time}
      </div>
    </div>
  );
}
