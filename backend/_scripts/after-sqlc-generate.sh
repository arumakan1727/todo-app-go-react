#!/bin/bash
set -euo pipefail

top_dir=$(realpath --relative-to="$(pwd)" "$(go list -m -f '{{.Dir}}')")
dir_sqlc="${top_dir}/repository/pgsql/sqlcgen"
dir_domain="${top_dir}/domain"
entity_go="${dir_domain}/entity.gen.go"

if [ -e "${dir_sqlc}/models.gen.go" ]; then
  mv -fv "${dir_sqlc}/models.gen.go" "$entity_go"
fi

# パッケージ名を domain へ変更する
sed -i -E 's/^package (\w)+/package domain/' "$entity_go"

# `type {...}$TYPE_SUFFIX struct` の {...} を抽出
typenames=$(sed -nE 's/^type (\w+) struct.*/\1/p' "$entity_go")

# 改行区切りの $typenames を '|' で join して正規表現として扱う
sed -i -E 's/\b('"$(echo -n "$typenames" | tr '\n' '|')"')\b/domain.\1/g' "$dir_sqlc"/*sql.gen.go
echo "[OK] Renamed struct type names"

# domain パッケージのインポートを消し、さらに 'domain.' を消す
domain_pkg_regex="$(go list -m | sed 's:/:\\/:g')"'\/domain'
sed -i -E /"$domain_pkg_regex"'/d; s/\bdomain\.//g' "$entity_go"
go fmt "$entity_go"
