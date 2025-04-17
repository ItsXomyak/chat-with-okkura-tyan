#!/bin/sh

set -e

host="$DB_HOST"
port="$DB_PORT"

echo "⏳ Ожидаем подключение к Postgres на $host:$port..."

until nc -z "$host" "$port"; do
  sleep 1
done

echo "✅ PostgreSQL доступен!"

exec /app/donate-service

