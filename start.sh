#!/bin/sh

set -e

echo "Run db migration"
/app/migration/

echo "start the app"
exec "$@"
