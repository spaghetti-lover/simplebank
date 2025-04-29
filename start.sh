#!/bin/sh

set -e

<<<<<<< HEAD
echo "run db migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

=======
>>>>>>> d4d0e58 (refactor)
echo "start the app"
exec "$@"
