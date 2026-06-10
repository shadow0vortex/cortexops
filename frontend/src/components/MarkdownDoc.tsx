import fs from "fs";
import path from "path";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface MarkdownDocProps {
  filePath: string;
}

export function MarkdownDoc({ filePath }: MarkdownDocProps) {
  const fullPath = path.join(process.cwd(), "..", filePath);
  let content = "";
  try {
    content = fs.readFileSync(fullPath, "utf-8");
  } catch (e) {
    content = `# File not found\nCould not load markdown from \`${filePath}\``;
  }

  return (
    <div className="w-full">
      <ReactMarkdown 
        remarkPlugins={[remarkGfm]}
        components={{
          h1: ({node, ...props}) => <h1 className="text-4xl font-bold text-white mb-4" {...props} />,
          h2: ({node, ...props}) => <h2 className="text-2xl font-semibold text-white mt-12 mb-6 border-b border-zinc-800 pb-2" {...props} />,
          h3: ({node, ...props}) => <h3 className="text-lg font-medium text-white mt-8 mb-3" {...props} />,
          h4: ({node, ...props}) => <h4 className="text-base font-medium text-zinc-300 mt-6 mb-2" {...props} />,
          p: ({node, ...props}) => <p className="text-zinc-300 mb-6 leading-relaxed" {...props} />,
          a: ({node, ...props}) => <a className="text-cortex-400 hover:text-cortex-300 transition-colors" {...props} />,
          code: ({node, ...props}) => {
            const isInline = !node?.position?.start.column || node.position.start.column > 1; // Basic inline heuristic if ReactMarkdown doesn't pass 'inline' properly here
            return (
              <code className="text-cortex-400 font-mono bg-zinc-900/50 border border-zinc-800 px-1.5 py-0.5 rounded text-sm" {...props} />
            );
          },
          pre: ({node, ...props}) => (
            <pre className="bg-black border border-zinc-800 rounded-lg p-4 overflow-x-auto text-sm text-zinc-300 mb-6 font-mono" {...props} />
          ),
          ul: ({node, ...props}) => <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2 mb-6" {...props} />,
          ol: ({node, ...props}) => <ol className="list-decimal pl-5 text-sm text-zinc-300 space-y-2 mb-6" {...props} />,
          li: ({node, ...props}) => <li {...props} />,
          blockquote: ({node, ...props}) => (
            <blockquote className="bg-zinc-900/50 border-l-4 border-cortex-500 p-6 rounded-r-xl mb-6 text-zinc-300 italic" {...props} />
          ),
          table: ({node, ...props}) => (
            <div className="bg-zinc-900/50 border border-zinc-800 rounded-xl overflow-hidden mb-6">
              <table className="w-full text-left text-sm text-zinc-300" {...props} />
            </div>
          ),
          thead: ({node, ...props}) => <thead className="bg-zinc-800 text-zinc-100" {...props} />,
          th: ({node, ...props}) => <th className="px-6 py-3 font-medium" {...props} />,
          td: ({node, ...props}) => <td className="px-6 py-4 border-t border-zinc-800" {...props} />,
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
