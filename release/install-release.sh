#!/bin/bash

VER="v1.0"

ARCH=$(uname -m)
VDIS="64"

if [ "$ARCH" == "i686" ] || [ "$ARCH" == "i386" ]; then
  VDIS="32"
fi

DOWNLOAD_LINK="https://github.com/v2ray/v2ray-core/releases/download/${VER}/v2ray-linux-${VDIS}.zip"

rm -rf /tmp/v2ray
mkdir /tmp/v2ray

CURL -L -o "/tmp/v2ray/v2ray.zip" ${DOWNLOAD_LINK}
unzip "/tmp/v2ray/v2ray.zip" -d "/tmp/v2ray/"

mkdir /usr/bin/v2ray
mkdir /etc/v2ray

cp -n "/tmp/v2ray/v2ray-${VER}-${VDIS}/*.json" "/etc/v2ray"
cp "/tmp/v2ray/v2ray-${VER}-${VDIS}/v2ray" "/usr/bin/v2ray/v2ray"
