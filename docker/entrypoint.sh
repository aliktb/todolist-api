#!/bin/sh
set -eu

if [ "${RUN_MIGRATIONS:-true}" = "true" ]; then
  echo "Running database migrations"
  todolist-migrate up
else
  echo "Skipping database migrations"
fi

echo "Starting API server"
exec "$@"
