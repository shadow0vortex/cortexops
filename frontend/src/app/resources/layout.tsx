export default function ResourcesLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="pt-32 pb-24 px-6 max-w-4xl mx-auto min-h-screen relative">
      <main className="glass-panel rounded-3xl p-8 md:p-12 prose prose-invert max-w-none">
        {children}
      </main>
    </div>
  );
}
