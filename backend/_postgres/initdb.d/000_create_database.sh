#!/bin/bash
set -eu

psql -U "$POSTGRES_USER" "$POSTGRES_DB" << EOT
  CREATE DATABASE ${POSTGRES_DB}__test;
EOT

set +eu
