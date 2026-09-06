#!/usr/bin/env bash
# ==============================================================================
# MI-Tech (StarTech Clone) — PostgreSQL Automated Backup Script
# Retention: 7 days
# Cron setup: 0 3 * * * /path/to/startech-clone/scripts/backup-db.sh >> /var/log/db-backup.log 2>&1
# ==============================================================================

set -e

BACKUP_DIR="${BACKUP_DIR:-/var/backups/mitech-db}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
FILENAME="mitech_backup_${TIMESTAMP}.sql.gz"

CONTAINER_NAME="mitech-postgres"
DB_NAME="${POSTGRES_DB:-start}"
DB_USER="${POSTGRES_USER:-postgres}"

mkdir -p "$BACKUP_DIR"

echo "[$(date)] 📦 Starting database backup of '${DB_NAME}'..."

# Dump and gzip directly from the postgres container
docker exec -t "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "${BACKUP_DIR}/${FILENAME}"

FILESIZE=$(du -h "${BACKUP_DIR}/${FILENAME}" | cut -f1)
echo "[$(date)] ✅ Backup completed: ${BACKUP_DIR}/${FILENAME} (${FILESIZE})"

# Remove backups older than 7 days
echo "[$(date)] 🧹 Removing backups older than 7 days..."
find "$BACKUP_DIR" -type f -name "mitech_backup_*.sql.gz" -mtime +7 -delete

echo "[$(date)] 🎉 Backup routine finished."
