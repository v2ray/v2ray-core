#!/bin/bash

PORT=27183

rand_str () {
        cat /dev/urandom | tr -dc 'a-f0-9' | fold -w $1 | head -n 1
}

ID="$(rand_str 8)-$(rand_str 4)-$(rand_str 4)-$(rand_str 4)-$(rand_str 12)"
echo "Generated client ID: $ID"

cat <<EOF > server-cfg.json
{
  "port": $PORT,
  "log" : {
    "access": "/go/access.log"
  },
  "inbound": {
    "protocol": "vmess",
    "settings": {
      "clients": [
        {"id": "$ID"}
      ]
    }
  },
  "outbound": {
    "protocol": "freedom",
    "settings": {}
  }
}
EOF
