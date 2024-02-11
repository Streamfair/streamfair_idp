#!/bin/sh

set -e

echo "Running DB migrations"
/streamfair_identity_provider/migrate -path /streamfair_identity_provider/migration -database "$DB_SOURCE_IDP" -verbose up

echo "Starting Identity Provider"
exec "$@"