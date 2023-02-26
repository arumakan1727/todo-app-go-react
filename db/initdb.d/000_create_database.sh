#!/bin/bash
set -eu

psql -U "$POSTGRES_USER" << EOT
  CREATE DATABASE ${POSTGRES_DB}__test;
EOT

set +eu
