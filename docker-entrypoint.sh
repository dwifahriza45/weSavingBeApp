#!/bin/sh
set -e

wait_for() {
  host="$1"
  port="$2"
  retries="${3:-30}"

  until nc -z "$host" "$port"; do
    retries=$((retries-1))
    if [ "$retries" -le 0 ]; then
      echo "timed out waiting for $host:$port"
      exit 1
    fi
    echo "waiting for $host:$port..."
    sleep 1
  done

  echo "$host:$port is available"
}

# ==========================
# WAIT FOR DATABASE (OPTIONAL)
# ==========================
if [ "${WAIT_FOR_DB:-false}" = "true" ]; then
  : "${DB_HOST:?DB_HOST is required when WAIT_FOR_DB=true}"
  DB_PORT="${DB_PORT:-5432}"
  wait_for "$DB_HOST" "$DB_PORT"
fi

# ==========================
# WAIT FOR MINIO (OPTIONAL)
# ==========================
if [ "${WAIT_FOR_MINIO:-false}" = "true" ]; then
  : "${MINIO_HOST:?MINIO_HOST is required when WAIT_FOR_MINIO=true}"
  MINIO_PORT="${MINIO_PORT:-9000}"
  wait_for "$MINIO_HOST" "$MINIO_PORT"
fi

# ==========================
# START APP
# ==========================
echo "Starting application..."
exec "$@"
