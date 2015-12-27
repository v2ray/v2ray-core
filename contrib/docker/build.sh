#!/bin/bash

if [ ! -e server-cfg.json ]; then
        ./gen-server-config.sh
fi

docker build --rm=true --tag=$USER/v2ray ./
