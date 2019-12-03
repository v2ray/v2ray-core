#!/usr/bin/env bash
case "$1" in

client)
    ./v2ray --config="./config/config-client.json"
;;

server)
    ./v2ray --config="./config/config-server.json"
;;

*)
echo "usage: ./run.sh client/server"

esac