import { MarkdownRenderer } from "@/components/MarkdownRenderer";

const content = `# Playbook: Backup & Restore

To maintain historical continuity, key state must be backed up.

## 1. Incident Audit Logs (Postgres)
The \`AuditStore\` (in a real system) resides in Postgres.
- **Backup**: \`pg_dump -U cortex cortexops > backup.sql\`
- **Restore**: \`psql -U cortex cortexops < backup.sql\`

## 2. Historical Incident Memory (Qdrant)
Embeddings stored in Qdrant power the RCA engine.
- **Backup**: Utilize Qdrant snapshots.
- **Location**: \`/qdrant/storage/snapshots/\`

## 3. Workflow History (Temporal)
Temporal workflow histories are stored in Postgres. Backing up the Postgres instance covers both audits and active workflow state.
`;

export default function BackupRestorePage() {
  return <MarkdownRenderer content={content} />;
}
