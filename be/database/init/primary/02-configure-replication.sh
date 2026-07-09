#!/bin/sh
set -eu

# Allow the replica to authenticate for streaming replication.
if ! grep -q "host replication replicator 0.0.0.0/0 scram-sha-256" "$PGDATA/pg_hba.conf"; then
  echo "host replication replicator 0.0.0.0/0 scram-sha-256" >> "$PGDATA/pg_hba.conf"
fi
