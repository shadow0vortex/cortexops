import { useDroppable } from "@dnd-kit/core";
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { IncidentCard, type Incident } from "./IncidentCard";

interface ColumnProps {
  id: string;
  title: string;
  incidents: Incident[];
}

export function Column({ id, title, incidents }: ColumnProps) {
  const { setNodeRef } = useDroppable({ id });

  return (
    <div className="flex flex-col bg-zinc-950/50 border border-zinc-800/50 rounded-2xl min-w-[300px] w-[300px] h-[70vh] max-h-[800px]">
      <div className="p-4 border-b border-zinc-800/50 flex items-center justify-between bg-zinc-900/30 rounded-t-2xl">
        <h3 className="text-sm font-bold text-zinc-300 uppercase tracking-wider">{title}</h3>
        <span className="text-xs font-mono bg-zinc-800 text-zinc-400 px-2 py-1 rounded-full">{incidents.length}</span>
      </div>
      
      <div ref={setNodeRef} className="p-4 flex-1 overflow-y-auto custom-scrollbar flex flex-col gap-1">
        <SortableContext items={incidents.map((i) => i.id)} strategy={verticalListSortingStrategy}>
          {incidents.map((incident) => (
            <IncidentCard key={incident.id} incident={incident} />
          ))}
        </SortableContext>
      </div>
    </div>
  );
}
