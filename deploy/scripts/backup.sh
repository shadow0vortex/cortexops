#!/bin/sh
set -e

BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "[$(date)] Starting CortexOps backup routine..."

# PostgreSQL Backup
PG_BACKUP_FILE="${BACKUP_DIR}/postgres_${TIMESTAMP}.sql.gz"
echo "[$(date)] Taking PostgreSQL pg_dump..."
pg_dump -h "${POSTGRES_HOST}" -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -F c | gzip > "${PG_BACKUP_FILE}"
echo "[$(date)] PostgreSQL backup saved to ${PG_BACKUP_FILE}"

# Qdrant Snapshot
echo "[$(date)] Taking Qdrant snapshot..."
curl -s -X POST "${QDRANT_URL}/collections/cortexops/snapshots"
echo "[$(date)] Qdrant snapshot triggered."

# Cleanup older than 7 days
find "${BACKUP_DIR}" -type f -name "*.sql.gz" -mtime +7 -exec rm {} \;

echo "[$(date)] Backup routine complete."
