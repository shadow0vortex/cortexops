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

        {/* Resource list skeleton */}
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="bg-[#0B0B0F] border border-zinc-800/50 rounded-xl p-5 flex items-center gap-4">
              <div className="w-10 h-10 rounded-lg bg-zinc-800/50 flex-shrink-0"></div>
              <div className="flex-1 space-y-2">
                <div className="h-4 w-1/3 bg-zinc-800/50 rounded"></div>
                <div className="h-3 w-2/3 bg-zinc-800/30 rounded"></div>
              </div>
              <div className="h-8 w-20 bg-zinc-800/40 rounded-lg"></div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
