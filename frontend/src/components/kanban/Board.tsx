"use client";

import { useState } from "react";
import { DndContext, DragOverlay, closestCorners, KeyboardSensor, PointerSensor, useSensor, useSensors, DragStartEvent, DragOverEvent, DragEndEvent } from "@dnd-kit/core";
import { sortableKeyboardCoordinates, arrayMove } from "@dnd-kit/sortable";
import { Column } from "./Column";
import { IncidentCard, type Incident } from "./IncidentCard";

const initialData: Record<string, Incident[]> = {
  detected: [
    { id: "inc-1", title: "Latency spike on checkout API", service: "checkout-svc", severity: "high", time: "2m ago" },
    { id: "inc-2", title: "OOMKilled pod detected", service: "payment-worker", severity: "critical", time: "5m ago" },
  ],
  correlating: [
    { id: "inc-3", title: "Database connection timeouts", service: "postgres-primary", severity: "critical", time: "10m ago" },
  ],
  awaiting: [
    { id: "inc-4", title: "Restart Kube-system pods", service: "kube-dns", severity: "high", time: "15m ago" },
  ],
  resolved: [
    { id: "inc-5", title: "CPU throttling resolved", service: "search-api", severity: "medium", time: "1h ago" },
  ],
};

const columnTitles: Record<string, string> = {
  detected: "Detected",
  correlating: "Correlating",
  awaiting: "Awaiting Human",
  resolved: "Resolved"
};

export function Board() {
  const [columns, setColumns] = useState(initialData);
  const [activeId, setActiveId] = useState<string | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  );

  const findContainer = (id: string) => {
    if (id in columns) return id;
    return Object.keys(columns).find((key) => columns[key].some((item) => item.id === id));
  };

  const handleDragStart = (event: DragStartEvent) => {
    setActiveId(event.active.id as string);
  };

  const handleDragOver = (event: DragOverEvent) => {
    const { active, over } = event;
    const overId = over?.id;
    
    if (!overId || active.id === overId) return;

    const activeContainer = findContainer(active.id as string);
    const overContainer = findContainer(overId as string);

    if (!activeContainer || !overContainer || activeContainer === overContainer) return;

    setColumns((prev) => {
      const activeItems = prev[activeContainer];
      const overItems = prev[overContainer];
      
      const activeIndex = activeItems.findIndex((item) => item.id === active.id);
      const overIndex = overId in prev 
        ? overItems.length + 1 
        : overItems.findIndex((item) => item.id === overId);

      return {
        ...prev,
        [activeContainer]: [...prev[activeContainer].filter((item) => item.id !== active.id)],
        [overContainer]: [
          ...prev[overContainer].slice(0, overIndex),
          activeItems[activeIndex],
          ...prev[overContainer].slice(overIndex, prev[overContainer].length)
        ]
      };
    });
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    const activeContainer = findContainer(active.id as string);
    const overContainer = over?.id ? findContainer(over.id as string) : null;

    if (!activeContainer || !overContainer || activeContainer !== overContainer) {
      setActiveId(null);
      return;
    }

    const activeIndex = columns[activeContainer].findIndex((item) => item.id === active.id);
    const overIndex = columns[overContainer].findIndex((item) => item.id === over?.id);

    if (activeIndex !== overIndex) {
      setColumns((prev) => ({
        ...prev,
        [overContainer]: arrayMove(prev[overContainer], activeIndex, overIndex)
      }));
    }

    setActiveId(null);
  };

  const activeItem = activeId 
    ? Object.values(columns).flat().find((i) => i.id === activeId) 
    : null;

  return (
    <DndContext sensors={sensors} collisionDetection={closestCorners} onDragStart={handleDragStart} onDragOver={handleDragOver} onDragEnd={handleDragEnd}>
      <div className="flex gap-6 overflow-x-auto pb-8 custom-scrollbar">
        {Object.keys(columns).map((key) => (
          <Column key={key} id={key} title={columnTitles[key]} incidents={columns[key]} />
        ))}
      </div>
      <DragOverlay>
        {activeItem ? <IncidentCard incident={activeItem} /> : null}
      </DragOverlay>
    </DndContext>
  );
}
