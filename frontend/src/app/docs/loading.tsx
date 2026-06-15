export default function DocsLoading() {
  return (
    <div className="animate-pulse space-y-8">
      {/* Title skeleton */}
      <div className="space-y-4">
        <div className="h-10 w-3/4 bg-zinc-800/60 rounded-lg"></div>
        <div className="h-5 w-full max-w-md bg-zinc-800/40 rounded-lg"></div>
      </div>

      {/* Section divider */}
      <div className="h-px w-full bg-zinc-800/50"></div>

      {/* Content skeleton */}
      <div className="space-y-6">
        {/* Paragraph block */}
        <div className="space-y-3">
          <div className="h-4 w-full bg-zinc-800/40 rounded"></div>
          <div className="h-4 w-11/12 bg-zinc-800/40 rounded"></div>
          <div className="h-4 w-4/5 bg-zinc-800/40 rounded"></div>
        </div>

        {/* Code block skeleton */}
        <div className="bg-black border border-zinc-800 rounded-lg p-4 space-y-2">
          <div className="h-3 w-2/3 bg-zinc-800/50 rounded"></div>
          <div className="h-3 w-1/2 bg-zinc-800/50 rounded"></div>
          <div className="h-3 w-3/4 bg-zinc-800/50 rounded"></div>
        </div>

        {/* Another paragraph block */}
        <div className="space-y-3">
          <div className="h-4 w-full bg-zinc-800/40 rounded"></div>
          <div className="h-4 w-5/6 bg-zinc-800/40 rounded"></div>
        </div>

        {/* Info card skeleton */}
        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl space-y-3">
          <div className="h-5 w-1/3 bg-zinc-800/50 rounded"></div>
          <div className="h-4 w-full bg-zinc-800/30 rounded"></div>
          <div className="h-4 w-4/5 bg-zinc-800/30 rounded"></div>
        </div>

        {/* List items skeleton */}
        <div className="space-y-3 pl-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="flex items-center gap-3">
              <div className="w-2 h-2 rounded-full bg-zinc-800/60 flex-shrink-0"></div>
              <div className="h-4 bg-zinc-800/40 rounded" style={{ width: `${60 + Math.random() * 30}%` }}></div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
