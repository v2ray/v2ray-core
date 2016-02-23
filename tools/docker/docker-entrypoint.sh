#!/bin/bash
set -e

if [ "$1" = 'v2ray' ]; then
  if [ ! -e "server-cfg.json" ]; then
      ./gen-server-cfg.sh
  fi
fi

exec "$@"
