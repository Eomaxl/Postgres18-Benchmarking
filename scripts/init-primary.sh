#!/bin/bash
set -e

# Initialize primary database with replication user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create replication user
    CREATE ROLE replicator WITH REPLICATION PASSWORD 'replicator_password' LOGIN;
    
    -- Grant necessary permissions
    GRANT CONNECT ON DATABASE $POSTGRES_DB TO replicator;
    
    -- Enable UUID extension
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    
    -- Enable pg_stat_statements for query performance monitoring
    CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
    
    -- Create users table with hash partitioning
    -- Note: Using gen_random_uuid() for now, will implement UUIDv7 in application layer
    -- Note: Email uniqueness will be enforced at application level due to partitioning constraints
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        metadata JSONB,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    ) PARTITION BY HASH (id);
    
    -- Create 4 hash partitions
    CREATE TABLE users_p0 PARTITION OF users
        FOR VALUES WITH (MODULUS 4, REMAINDER 0);
    
    CREATE TABLE users_p1 PARTITION OF users
        FOR VALUES WITH (MODULUS 4, REMAINDER 1);
    
    CREATE TABLE users_p2 PARTITION OF users
        FOR VALUES WITH (MODULUS 4, REMAINDER 2);
    
    CREATE TABLE users_p3 PARTITION OF users
        FOR VALUES WITH (MODULUS 4, REMAINDER 3);
    
    -- Create indexes
    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
    CREATE INDEX IF NOT EXISTS idx_users_metadata ON users USING GIN(metadata);
    
    -- Multicolumn index for skip scan optimization (PostgreSQL 18 feature)
    CREATE INDEX IF NOT EXISTS idx_users_name_email ON users(name, email);
EOSQL

echo "Primary database initialized successfully"
