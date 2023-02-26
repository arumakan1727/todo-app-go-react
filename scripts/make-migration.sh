#!/bin/bash
set -euo pipefail

f="/tmp/todoapp-sqldef.sql"

trap cleanup SIGINT SIGTERM ERR EXIT

cleanup() {
  trap - SIGINT SIGTERM ERR EXIT
  rm -f "$f"
}

psqldef -Utodoapp -p25432 -h 127.0.0.1 --dry-run --config=.psqldef.yml todoapp < ./db/schema.sql > "$f"

if [ "$(head -1 "$f")" = "-- Nothing is modified --" ]; then
  echo "Nothing is modified. Bye!"
  exit
fi

bat "$f" || cat "$f"

echo ""
read -rp "Migration name: " migration_name

if [ -z "$migration_name" ]; then
  echo "Error: migration name is empty. exit."
  exit 1
fi

migration_file="$(dbmate new "$migration_name" | sed -E 's/^\s*Creating migration:\s*(\S+)$/\1/i')"

( echo '-- migrate:up'; tail -n +2 "$f"; echo -e '\n-- migrate: down'; ) > "$migration_file"

echo "OK: created $migration_file"
