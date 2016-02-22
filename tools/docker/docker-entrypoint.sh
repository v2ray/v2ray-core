#!/bin/bash
set -e

if [ ! -e "server-cfg.json" ]; then
    ./gen-server-cfg.sh
fi

exec "$@"
