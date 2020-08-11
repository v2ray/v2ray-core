#!/bin/bash

pushd "$GOPATH/src/v2ray.com/core/" || return

# Update geoip.dat
curl -L -o release/config/geoip.dat "https://github.com/v2ray/geoip/raw/release/geoip.dat"
sleep 1

# Update geosite.dat
curl -L -o release/config/geosite.dat "https://github.com/v2ray/domain-list-community/raw/release/dlc.dat"
sleep 1

popd || return
