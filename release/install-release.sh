#!/bin/bash

while [[ $# > 0 ]]
do
key="$1"

case $key in
    -p|--proxy)
    PROXY="$2"
    shift # past argument
    ;;
    -h|--help)
    HELP="1"
    ;;
    -f|--force)
    FORCE="1"
    ;;
    --version)
    VERSION="$2"
    shift
    ;;
    --local)
    LOCAL="$2"
    shift
    ;;
    *)
            # unknown option
    ;;
esac
shift # past argument or value
done

if [[ "$HELP" == "1" ]]; then
  echo "./install-release.sh [-p proxy] [-h] [-f] [--version vx.y.z] [--local file]"
  echo "-p: To download through a proxy server, use -p socks5://127.0.0.1:1080 or -p http://127.0.0.1:3128 etc"
  echo "-h: Show help"
  echo "-f: Force install"
  echo "--version: Install a particular version"
  echo "--local: Install from a local file"
  exit
fi

YUM_CMD=$(command -v yum)
APT_CMD=$(command -v apt-get)

SOFTWARE_UPDATED=0

function update_software() {
  if [ ${SOFTWARE_UPDATED} -eq 1 ]; then
    return
  fi
  if [ -n "${YUM_CMD}" ]; then
    echo "Updating software repo via yum."
    ${YUM_CMD} -q makecache
  elif [ -n "${APT_CMD}" ]; then
    echo "Updating software repo via apt-get."
    ${APT_CMD} -qq update
  fi
  SOFTWARE_UPDATED=1
}

function install_component() {
  local COMPONENT=$1
  COMPONENT_CMD=$(command -v $COMPONENT)
  if [ -n "${COMPONENT_CMD}" ]; then
    return
  fi

  update_software
  if [ -n "${YUM_CMD}" ]; then
    echo "Installing ${COMPONENT} via yum."
    ${YUM_CMD} -y -q install $COMPONENT
  elif [ -n "${APT_CMD}" ]; then
    echo "Installing ${COMPONENT} via apt-get."
    ${APT_CMD} -y -qq install $COMPONENT
  fi
}

V2RAY_RUNNING=0
if pgrep "v2ray" > /dev/null ; then
  V2RAY_RUNNING=1
fi

if [ -n "$VERSION" ]; then
  VER="$VERSION"
else
  VER="$(curl -s https://api.github.com/repos/v2ray/v2ray-core/releases/latest | grep 'tag_name' | cut -d\" -f4)"
  CUR_VER="$(/usr/bin/v2ray/v2ray -version | head -n 1 | cut -d " " -f2)"

  if [[ "$VER" == "$CUR_VER" ]] && [[ "$FORCE" != "1" ]]; then
    echo "Lastest version $VER is already installed. Exiting..."
    exit
  fi
fi

ARCH=$(uname -m)
VDIS="64"

if [[ "$ARCH" == "i686" ]] || [[ "$ARCH" == "i386" ]]; then
  VDIS="32"
elif [[ "$ARCH" == *"armv7"* ]] || [[ "$ARCH" == "armv6l" ]]; then
  VDIS="arm"
elif [[ "$ARCH" == *"armv8"* ]]; then
  VDIS="arm64"
fi

rm -rf /tmp/v2ray
mkdir -p /tmp/v2ray

echo "Installing V2Ray ${VER} on ${ARCH}"

if [ -n "$LOCAL" ]; then
  cp "$LOCAL" "/tmp/v2ray/v2ray.zip"
else
  DOWNLOAD_LINK="https://github.com/v2ray/v2ray-core/releases/download/${VER}/v2ray-linux-${VDIS}.zip"

  install_component "curl"

  if [ -n "${PROXY}" ]; then
    echo "Downloading ${DOWNLOAD_LINK} via proxy ${PROXY}."
    curl -x ${PROXY} -L -H "Cache-Control: no-cache" -o "/tmp/v2ray/v2ray.zip" ${DOWNLOAD_LINK}
  else
    echo "Downloading ${DOWNLOAD_LINK} directly."
    curl -L -H "Cache-Control: no-cache" -o "/tmp/v2ray/v2ray.zip" ${DOWNLOAD_LINK}
  fi
fi

echo "Extracting V2Ray package to /tmp/v2ray."
install_component "unzip"
unzip "/tmp/v2ray/v2ray.zip" -d "/tmp/v2ray/"

# Create folder for V2Ray log.
mkdir -p /var/log/v2ray

# Stop v2ray daemon if necessary.
SYSTEMCTL_CMD=$(command -v systemctl)
SERVICE_CMD=$(command -v service)

if [ ${V2RAY_RUNNING} -eq 1 ]; then
  echo "Shutting down V2Ray service."
  if [ -n "${SYSTEMCTL_CMD}" ]; then
    if [ -f "/lib/systemd/system/v2ray.service" ]; then
      ${SYSTEMCTL_CMD} stop v2ray
    fi
  elif [ -n "${SERVICE_CMD}" ]; then
    if [ -f "/etc/init.d/v2ray" ]; then
      ${SERVICE_CMD} v2ray stop
    fi
  fi
fi

# Install V2Ray binary to /usr/bin/v2ray
mkdir -p /usr/bin/v2ray
cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/v2ray" "/usr/bin/v2ray/v2ray"
chmod +x "/usr/bin/v2ray/v2ray"

# Install V2Ray server config to /etc/v2ray
mkdir -p /etc/v2ray
if [ ! -f "/etc/v2ray/config.json" ]; then
  cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/vpoint_vmess_freedom.json" "/etc/v2ray/config.json"

  let PORT=$RANDOM+10000
  sed -i "s/10086/${PORT}/g" "/etc/v2ray/config.json"

  UUID=$(cat /proc/sys/kernel/random/uuid)
  sed -i "s/23ad6b10-8d1a-40f7-8ad0-e3e35cd38297/${UUID}/g" "/etc/v2ray/config.json"

  echo "PORT:${PORT}"
  echo "UUID:${UUID}"
fi

if [ -n "${SYSTEMCTL_CMD}" ]; then
  if [ ! -f "/lib/systemd/system/v2ray.service" ]; then
    cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/systemd/v2ray.service" "/lib/systemd/system/"
    systemctl enable v2ray
  else
    if [ ${V2RAY_RUNNING} -eq 1 ]; then
      echo "Restarting V2Ray service."
      ${SYSTEMCTL_CMD} start v2ray
    fi
  fi
elif [ -n "${SERVICE_CMD}" ]; then # Configure SysV if necessary.
  if [ ! -f "/etc/init.d/v2ray" ]; then
    install_component "daemon"
    cp "/tmp/v2ray/v2ray-${VER}-linux-${VDIS}/systemv/v2ray" "/etc/init.d/v2ray"
    chmod +x "/etc/init.d/v2ray"
    update-rc.d v2ray defaults
  else
    if [ ${V2RAY_RUNNING} -eq 1 ]; then
      echo "Restarting V2Ray service."
      ${SERVICE_CMD} v2ray start
    fi
  fi
fi

echo "V2Ray ${VER} is installed."