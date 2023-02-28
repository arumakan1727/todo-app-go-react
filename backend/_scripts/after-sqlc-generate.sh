#!/bin/bash
set -euo pipefail

top_dir=$(realpath --relative-to="$(pwd)" "$(go list -m -f '{{.Dir}}')")
dir_sql="${top_dir}/repository/pgsql"
dir_domain="${top_dir}/domain"
entity_go="${dir_domain}/entity.gen.go"

if [ -e "${dir_sql}/models.gen.go" ]; then
  mv -fv "${dir_sql}/models.gen.go" "$entity_go"
fi

TYPE_SUFFIX='Entity'

# $TYPE_SUFFIX がまだついていない構造体型名にのみ $TYPE_SUFFIX を付与
# パッケージ名も domain へ変更する
perl -i -pe 's/^type (?!\w+'"$TYPE_SUFFIX"')(\w+)/type \1'"$TYPE_SUFFIX"'/; s/^package (\w)+/package domain/' "$entity_go"

# `type {...}$TYPE_SUFFIX struct` の {...} を抽出
typenames=$(sed -nE 's/^type (\w+)'"$TYPE_SUFFIX"' struct.*/\1/p' "$entity_go")

# 改行区切りの $typenames を '|' で join して正規表現として扱う
sed -i -E 's/\b('"$(echo -n "$typenames" | tr '\n' '|')"')\b/domain.\1'"$TYPE_SUFFIX"'/g' "$dir_sql"/*sql.gen.go
echo "[OK] Renamed struct type names"

# domain パッケージのインポートを消し、さらに 'domain.' を消す
domain_pkg_regex="$(go list -m | sed 's:/:\\/:g')"'\/domain'
sed -i -E /"$domain_pkg_regex"'/d; s/\bdomain\.//g' "$entity_go"
go fmt "$entity_go"
