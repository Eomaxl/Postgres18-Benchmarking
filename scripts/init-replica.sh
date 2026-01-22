#!/bin/bash
set -e

# Get the PostgreSQL data directory
PGDATA="/var/lib/postgresql/data"

# Wait for primary to be ready
echo "Waiting for primary database to be ready..."
until PGPASSWORD=postgres psql -h postgres-primary -U postgres -d userdb -c '\q' 2>/dev/null; do
  echo "Primary database is unavailable - sleeping"
  sleep 2
done

echo "Primary database is ready"

# Stop PostgreSQL if running
pg_ctl -D "$PGDATA" stop || true

# Remove existing data directory
rm -rf "$PGDATA"/*

# Create base backup from primary
echo "Creating base backup from primary..."
PGPASSWORD=replicator_password pg_basebackup -h postgres-primary -D "$PGDATA" -U replicator -v -P -W -R

# Create standby.signal file to indicate this is a replica
touch "$PGDATA/standby.signal"

# Configure replication connection
cat >> "$PGDATA/postgresql.auto.conf" <<EOF
primary_conninfo = 'host=postgres-primary port=5432 user=replicator password=replicator_password'
primary_slot_name = 'replica_slot_$(hostname)'
EOF

echo "Replica initialized successfully"
