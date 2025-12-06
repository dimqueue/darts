#!/bin/sh
set -e

echo "Waiting for database..."
until darts migrates-up 2>/dev/null; do
    echo "Database not ready, retrying in 2s..."
    sleep 2
done

echo "Starting server..."
exec darts run-server