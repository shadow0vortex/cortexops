export default function SectionLoading() {
  return (
    <div className="pt-32 pb-24 px-6 max-w-7xl mx-auto min-h-screen">
      <div className="animate-pulse space-y-12">
        {/* Header */}
        <div className="flex flex-col items-center text-center space-y-6">
          <div className="h-4 w-48 bg-zinc-800/50 rounded-full"></div>
          <div className="h-10 w-3/4 max-w-2xl bg-zinc-800/60 rounded-lg"></div>
          <div className="h-5 w-2/3 max-w-xl bg-zinc-800/40 rounded-lg"></div>
        </div>

        {/* Dashboard-style skeleton */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-[#0B0B0F] border border-zinc-800/50 rounded-2xl p-6 space-y-3">
              <div className="h-4 w-1/3 bg-zinc-800/50 rounded"></div>
              <div className="h-8 w-2/3 bg-zinc-800/60 rounded"></div>
              <div className="h-2 w-full bg-zinc-800/30 rounded-full"></div>
            </div>
          ))}
        </div>

        {/* Chart-like area */}
        <div className="bg-[#0B0B0F] border border-zinc-800/50 rounded-2xl p-8 space-y-4">
          <div className="h-5 w-1/4 bg-zinc-800/50 rounded"></div>
          <div className="h-48 w-full bg-zinc-800/20 rounded-xl"></div>
        </div>
      </div>
    </div>
  );
}
