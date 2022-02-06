#!/bin/sh

set -e

echo "Run db migrations..."
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting the app..."
exec "$@"