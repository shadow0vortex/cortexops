import fs from 'fs';
import path from 'path';

// This script copies the playbooks from the root docs directory into the frontend directory
// so that Vercel's Next.js build can trace and include them in the Serverless Functions.

const sourceDir = path.join(process.cwd(), '../docs/playbooks');
const targetDir = path.join(process.cwd(), 'playbooks');

if (!fs.existsSync(targetDir)) {
  fs.mkdirSync(targetDir, { recursive: true });
}

if (fs.existsSync(sourceDir)) {
  const files = fs.readdirSync(sourceDir);
  for (const file of files) {
    if (file.endsWith('.md')) {
      fs.copyFileSync(path.join(sourceDir, file), path.join(targetDir, file));
    }
  }
  console.log('Successfully copied playbooks to frontend/playbooks for Vercel build.');
} else {
  console.warn('Source playbooks directory not found at:', sourceDir);
}
