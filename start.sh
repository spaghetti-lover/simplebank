#!/bin/sh

set -e

echo "run db migration"
source /app/app.env
source app.env
/app/migrate -path /app/migration -database "postgresql://root:Narutora12345@bank-system.cfcwm8aiq0n8.ap-southeast-2.rds.amazonaws.com:5432/postgres" -verbose up

echo "start the app"
exec "$@"