#!/bin/bash

YUM_CMD=$(command -v yum)
APT_CMD=$(command -v apt-get)

if [ -n "${YUM_CMD}" ]; then
  echo "Installing unzip and daemon via yum."
  ${YUM_CMD} -q makecache
  ${YUM_CMD} -y -q install unzip daemon
elif [ -n "${APT_CMD}" ]; then
  echo "Installing unzip and daemon via apt-get."
  ${APT_CMD} -qq update
  ${APT_CMD} -y -qq install unzip daemon
else
  echo "Please make sure unzip and daemon are installed."
fi

VER="v1.4"

ARCH=$(uname -m)
VDIS="64"

if [[ "$ARCH" == "i686" ]] || [[ "$ARCH" == "i386" ]]; then
  VDIS="32"
elif [[ "$ARCH" == *"armv7"* ]]; then
  VDIS="arm"
elif [[ "$ARCH" == *"armv8"* ]]; then
  VDIS="arm64"
fi

DOWNLOAD_LINK="https://github.com/v2ray/v2ray-core/releases/download/${VER}/v2ray-linux-${VDIS}.zip"

rm -rf /tmp/v2ray
mkdir -p /tmp/v2ray

curl -L -o "/tmp/v2ray/v2ray.zip" ${DOWNLOAD_LINK}
unzip "/tmp/v2ray/v2ray.zip" -d "/tmp/v2ray/"

# Create folder for V2Ray log.
mkdir -p /var/log/v2ray

# Install V2Ray binary to /usr/bin/v2ray
mkdir -p /usr/bin/v2ray
cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/v2ray" "/usr/bin/v2ray/v2ray"
chmod +x "/usr/bin/v2ray/v2ray"

# Install V2Ray server config to /etc/v2ray
mkdir -p /etc/v2ray
if [ ! -f "/etc/v2ray/config.json" ]; then
  cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/vpoint_vmess_freedom.json" "/etc/v2ray/config.json"

  let PORT=$RANDOM+10000
  sed -i "s/37192/${PORT}/g" "/etc/v2ray/config.json"

  UUID=$(cat /proc/sys/kernel/random/uuid)
  sed -i "s/3b129dec-72a3-4d28-aeee-028a0fe86e22/${UUID}/g" "/etc/v2ray/config.json"

  echo "PORT:${PORT}"
  echo "UUID:${UUID}"
fi

if [ -d "/lib/systemd/system" ]; then
  if [ ! -f "/lib/systemd/system/v2ray.service" ]; then
    cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/systemd/v2ray.service" "/lib/systemd/system/"
    systemctl enable v2ray
  fi
elif [ -d "/etc/init.d" ]; then # Configure SysV if necessary.
  if [ ! -f "/etc/init.d/v2ray" ]; then
    cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/systemv/v2ray" "/etc/init.d/v2ray"
    chmod +x "/etc/init.d/v2ray"
    update-rc.d v2ray defaults
  fi
fi
