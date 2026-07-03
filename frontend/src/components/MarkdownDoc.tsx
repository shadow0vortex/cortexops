import fs from "fs";
import path from "path";
import { MarkdownRenderer } from "./MarkdownRenderer";

interface MarkdownDocProps {
  filePath: string;
}

export function MarkdownDoc({ filePath }: MarkdownDocProps) {
  const fullPath = path.join(process.cwd(), "..", filePath);
  let content = "";
  try {
    content = fs.readFileSync(fullPath, "utf-8");
  } catch {
    content = `# File not found\nCould not load markdown from \`${filePath}\``;
  }

  return (
    <div className="w-full">
      <MarkdownRenderer content={content} />
    </div>
  );
}
