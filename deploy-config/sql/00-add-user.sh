#!/bin/bash
set -euo pipefail  # Strict mode: fail on errors, unset vars, pipe failures

echo "Initializing application role and database..."

# Create the app role idempotently
psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'vuecom') THEN
            CREATE ROLE vuecom WITH LOGIN PASSWORD '${APP_PG_PASSWORD}';
        END IF;
    END
    \$\$;

    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO vuecom;
EOSQL

echo "Initialization complete."