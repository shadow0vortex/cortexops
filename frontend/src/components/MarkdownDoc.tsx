import fs from "fs";
import path from "path";
import { MarkdownRenderer } from "./MarkdownRenderer";

interface MarkdownDocProps {
  filePath: string;
}

export function MarkdownDoc({ filePath }: MarkdownDocProps) {
  const fileName = path.basename(filePath);
  
  // Vercel deployment: read from the copied directory so Next.js static analysis includes it
  const localCopyPath = path.join(process.cwd(), "playbooks", fileName);
  // Fallback to parent directory for local dev if script hasn't run
  const originalPath = path.join(process.cwd(), "..", filePath);

  let content = "";
  try {
    if (fs.existsSync(localCopyPath)) {
      content = fs.readFileSync(localCopyPath, "utf-8");
    } else {
      content = fs.readFileSync(originalPath, "utf-8");
    }
  } catch {
    content = `# File not found\nCould not load markdown from \`${filePath}\``;
  }

  return (
    <div className="w-full">
      <MarkdownRenderer content={content} />
    </div>
  );
}
