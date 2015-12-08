#!/bin/bash

VER="v1.1"

ARCH=$(uname -m)
VDIS="64"

if [ "$ARCH" == "i686" ] || [ "$ARCH" == "i386" ]; then
  VDIS="32"
fi

DOWNLOAD_LINK="https://github.com/v2ray/v2ray-core/releases/download/${VER}/v2ray-linux-${VDIS}.zip"

rm -rf /tmp/v2ray
mkdir -p /tmp/v2ray

curl -L -o "/tmp/v2ray/v2ray.zip" ${DOWNLOAD_LINK}
unzip "/tmp/v2ray/v2ray.zip" -d "/tmp/v2ray/"

mkdir -p /usr/bin/v2ray
mkdir -p /etc/v2ray
mkdir -p /var/log/v2ray

cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/v2ray" "/usr/bin/v2ray/v2ray"

if [ ! -f "/etc/v2ray/config.json" ]; then
  cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/vpoint_vmess_freedom.json" "/etc/v2ray/config.json"
  
  #PORT=$(expr $RANDOM + 10000)
  #sed -i "s/37192/${PORT}/g" "/etc/v2ray/config.json"

  UUID=$(cat /proc/sys/kernel/random/uuid)
  sed -i "s/3b129dec-72a3-4d28-aeee-028a0fe86e22/${UUID}/g" "/etc/v2ray/config.json"

  #echo "PORT:${PORT}"
  echo "UUID:${UUID}"
fi
