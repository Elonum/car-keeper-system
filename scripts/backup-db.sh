#!/usr/bin/env sh
# Logical backup of CarKeeper PostgreSQL. Requires pg_dump and a reachable database.
# Usage: PGHOST=localhost PGUSER=postgres PGDATABASE=carkeeper ./scripts/backup-db.sh

set -eu
STAMP=$(date +%Y%m%d-%H%M%S)
OUT="${BACKUP_DIR:-.}/carkeeper-${STAMP}.sql.gz"
mkdir -p "$(dirname "$OUT")"
pg_dump --no-owner --format=plain "${PGDATABASE:-carkeeper}" | gzip -c >"$OUT"
echo "Wrote $OUT"
