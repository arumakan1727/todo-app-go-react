#!/bin/bash
set -euo pipefail
read -rp "Enter migration name: " migration_name

if [ -z "$migration_name" ]; then
  echo "Error: migration name is empty. Bye."
  exit 1
fi

migration_file="$(dbmate new "$migration_name" | sed -E 's/^\s*Creating migration:\s*(\S+)$/\1/i')"
echo "OK: created $migration_file"

read -rp "Edit file? [y/N]" yn
case "$yn" in
  [yY]* )
    open_cmd='open'
    command -v xdg-open >/dev/null && open_cmd='xdg-open'
    command -v "$EDITOR" >/dev/null && open_cmd="$EDITOR"
    $open_cmd "$migration_file"
    ;;

  * )
    echo 'Bye!'
    ;;
esac
