#!/bin/sh
set -eu

PRIMARY_HOST="${PRIMARY_HOST:-db-primary}"
REPLICATION_USER="${REPLICATION_USER:-replicator}"
REPLICATION_PASSWORD="${REPLICATION_PASSWORD:-replicator_password}"

# Keep retrying until the primary is reachable for base backup.
until pg_isready -h "$PRIMARY_HOST" -U postgres -d postgres; do
  echo "Waiting for primary database at $PRIMARY_HOST..."
  sleep 2
done

mkdir -p "$PGDATA"
chmod 700 "$PGDATA"

if [ -z "$(ls -A "$PGDATA" 2>/dev/null)" ]; then
  echo "Initializing replica from primary..."
  PGPASSWORD="$REPLICATION_PASSWORD" pg_basebackup \
    -h "$PRIMARY_HOST" \
    -D "$PGDATA" \
    -U "$REPLICATION_USER" \
    -vP \
    -R
fi

  chown -R postgres:postgres "$PGDATA"
  exec gosu postgres postgres
