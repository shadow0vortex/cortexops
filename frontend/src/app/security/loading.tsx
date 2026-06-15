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

        {/* Security cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="bg-[#0B0B0F] border border-zinc-800/50 rounded-2xl p-6 space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-zinc-800/50"></div>
                <div className="h-5 w-1/2 bg-zinc-800/50 rounded"></div>
              </div>
              <div className="space-y-2">
                <div className="h-3 w-full bg-zinc-800/30 rounded"></div>
                <div className="h-3 w-4/5 bg-zinc-800/30 rounded"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
