#!/bin/bash

# This file is accessible as https://install.direct/go.sh
# Original source is located at github.com/v2ray/v2ray-core/release/install-release.sh

CUR_VER=""
NEW_VER=""
ARCH=""
VDIS="64"
ZIPFILE="/tmp/v2ray/v2ray.zip"
V2RAY_RUNNING=0

CMD_INSTALL=""
CMD_UPDATE=""
SOFTWARE_UPDATED=0

CHECK=""
FORCE=""
HELP=""

#######color code########
red="31m"
green="32m"
yellow="33m"
blue="34m"
wblue="36m"

#########################
while [[ $# > 0 ]];do
    key="$1"
    case $key in
        -p|--proxy)
        PROXY="-x ${2}"
        shift # past argument
        ;;
        -h|--help)
        HELP="1"
        ;;
        -f|--force)
        FORCE="1"
        ;;
        -c|--check)
        CHECK="1"
        ;;
        -r|--remove)
        REMOVE="1"
        ;;
        -v|--version)
        VERSION="$2"
        shift
        ;;
        -l|--local)
        LOCAL="$2"
        LOCAL_INSTALL="1"
        shift
        ;;
        *)
                # unknown option
        ;;
    esac
    shift # past argument or value
done

###############################
colorEcho(){
    color=$1
    text=$2
    echo -e "\033[${color}${@:2}\033[0m"
}

sysAcrh(){
    ARCH=$(uname -m)
    if [[ "$ARCH" == "i686" ]] || [[ "$ARCH" == "i386" ]]; then
        VDIS="32"
    elif [[ "$ARCH" == *"armv7"* ]] || [[ "$ARCH" == "armv6l" ]]; then
        VDIS="arm"
    elif [[ "$ARCH" == *"armv8"* ]]; then
        VDIS="arm64"
    fi
    return 0
}

downloadV2Ray(){
    rm -rf /tmp/v2ray
    mkdir -p /tmp/v2ray
    colorEcho ${wblue} "Donwloading V2Ray."
    DOWNLOAD_LINK="https://github.com/v2ray/v2ray-core/releases/download/${NEW_VER}/v2ray-linux-${VDIS}.zip"
    curl ${PROXY} -L -H "Cache-Control: no-cache" -o ${ZIPFILE} ${DOWNLOAD_LINK}
    if [ $? != 0 ];then
        colorEcho ${red} "Failed to download! Please check your network or try again."
        exit 1
    fi
    return 0
}

installSoftware(){
    COMPONENT=$1
    if [[ -n `command -v $COMPONENT` ]]; then
        return 0
    fi

    getPMT
    if [[ $? -eq 1 ]]; then
        colorEcho $yellow "The system package manager tool isn't APT or YUM, please install ${COMPONENT} manually."
        exit 
    fi
    colorEcho $green Installing $COMPONENT 
    if [[ $SOFTWARE_UPDATED -eq 0 ]]; then
        colorEcho ${wblue} "Updating software repo"
        $CMD_UPDATE
        if [[ $? -ne 0 ]]; then
            colorEcho ${red} "Failed update software repo, please check your source."
            exit
        fi        
        SOFTWARE_UPDATED=1
    fi

    colorEcho ${wblue} "Installing ${COMPONENT}"
    $CMD_INSTALL $COMPONENT
    if [[ $? -ne 0 ]]; then
        colorEcho ${red} Install ${COMPONENT} fail, please install it manually.
        exit
    fi
    return 0
}

# return 1: not apt or yum
getPMT(){
    if [ -n `command -v apt-get` ];then
        CMD_INSTALL="apt-get -y -qq install"
        CMD_UPDATE="apt-get -qq update"
    elif [[ -n `command -v yum` ]]; then
        CMD_INSTALL="yum -y -qq install"
        CMD_UPDATE="yum -q makecache"
    else
        return 1
    fi
    return 0
}


extra(){
    colorEcho ${wblue}"Extracting V2Ray package to /tmp/v2ray."
    mkdir -p /tmp/v2ray
    unzip $1 -d "/tmp/v2ray/"
    if [[ $? -ne 0 ]]; then
        colorEcho ${red} "Extracting V2Ray faile!"
        exit
    fi
    return 0
}


# 1: new V2Ray. 0: no
getVersion(){
    if [[ -n "$VERSION" ]]; then
        NEW_VER="$VERSION"
        return 1
    else
        CUR_VER=`/usr/bin/v2ray/v2ray -version 2>/dev/null | head -n 1 | cut -d " " -f2`
        TAG_URL="https://api.github.com/repos/v2ray/v2ray-core/releases/latest"
        NEW_VER=`curl ${PROXY} -s ${TAG_URL} --connect-timeout 10| grep 'tag_name' | cut -d\" -f4`

        if [[ $? -ne 0 ]] || [[ $NEW_VER == "" ]]; then
            colorEcho ${red} "Network error! Please check your network or try again."
            exit
        elif [[ "$NEW_VER" != "$CUR_VER" ]];then
                return 1
        fi
        return 0
    fi
}

stopV2ray(){
    SYSTEMCTL_CMD=$(command -v systemctl)
    SERVICE_CMD=$(command -v service)

    colorEcho ${wblue} "Shutting down V2Ray service."
    if [[ -n "${SYSTEMCTL_CMD}" ]] || [[ -f "/lib/systemd/system/v2ray.service" ]]; then
        ${SYSTEMCTL_CMD} stop v2ray
    elif [[ -n "${SERVICE_CMD}" ]] || [[ -f "/etc/init.d/v2ray" ]]; then
        ${SERVICE_CMD} v2ray stop
    fi
    return 0
}

startV2ray(){
    SYSTEMCTL_CMD=$(command -v systemctl)
    SERVICE_CMD=$(command -v service)

    if [ -n "${SYSTEMCTL_CMD}" ] && [ -f "/lib/systemd/system/v2ray.service" ]; then
        ${SYSTEMCTL_CMD} start v2ray
    elif [ -n "${SERVICE_CMD}" ] && [ -f "/etc/init.d/v2ray" ]; then
        ${SERVICE_CMD} v2ray start
    fi
    return 0
}

installV2Ray(){
    # Install V2Ray binary to /usr/bin/v2ray
    mkdir -p /usr/bin/v2ray
    ERROR=`cp "/tmp/v2ray/v2ray-${NEW_VER}-linux-${VDIS}/v2ray" "/usr/bin/v2ray/v2ray"`
    if [[ $? -ne 0 ]]; then
          colorEcho ${yellow} "${ERROR}"
          exit
    fi
    chmod +x "/usr/bin/v2ray/v2ray"

    # Install V2Ray server config to /etc/v2ray
    mkdir -p /etc/v2ray
    if [[ ! -f "/etc/v2ray/config.json" ]]; then
      cp "/tmp/v2ray/v2ray-${NEW_VER}-linux-${VDIS}/vpoint_vmess_freedom.json" "/etc/v2ray/config.json"
      if [[ $? -ne 0 ]]; then
          colorEcho ${yellow} "Create V2Ray configuration file error, pleases create it manually."
          return 1
      fi
      let PORT=$RANDOM+10000
      UUID=$(cat /proc/sys/kernel/random/uuid)

      sed -i "s/10086/${PORT}/g" "/etc/v2ray/config.json"
      sed -i "s/23ad6b10-8d1a-40f7-8ad0-e3e35cd38297/${UUID}/g" "/etc/v2ray/config.json"

      colorEcho ${green} "PORT:${PORT}"
      colorEcho ${green} "UUID:${UUID}"
      mkdir -p /var/log/v2ray
    fi
    return 0
}


installInitScrip(){
    SYSTEMCTL_CMD=$(command -v systemctl)
    SERVICE_CMD=$(command -v service)

    if [[ -n "${SYSTEMCTL_CMD}" ]];then
        if [[ ! -f "/lib/systemd/system/v2ray.service" ]]; then
            cp "/tmp/v2ray/v2ray-${NEW_VER}-linux-${VDIS}/systemd/v2ray.service" "/lib/systemd/system/"
            systemctl enable v2ray
        fi
        return
    elif [[ -n "${SERVICE_CMD}" ]] && [[ ! -f "/etc/init.d/v2ray" ]]; then
        installSoftware "daemon"
        cp "/tmp/v2ray/v2ray-${NEW_VER}-linux-${VDIS}/systemv/v2ray" "/etc/init.d/v2ray"
        chmod +x "/etc/init.d/v2ray"
        update-rc.d v2ray defaults
    fi
    return
}

Help(){
    echo "./install-release.sh [-h] [-c] [-p proxy] [-f] [-v vx.y.z] [-l file]"
    echo "  -h, --help            Show help"
    echo "  -p, --proxy           To download through a proxy server, use -p socks5://127.0.0.1:1080 or -p http://127.0.0.1:3128 etc"
    echo "  -f, --force           Force install"
    echo "  -v, --version         Install a particular version"
    echo "  -l, --local           Install from a local file"
    echo "  -r, --remove          Remove installed V2Ray"
    echo "  -c, --check           Check for update"
    exit  
}

remove(){
    if pgrep "v2ray" > /dev/null ; then
        stopV2ray
    fi
    rm -rf "/usr/lib/v2ray" "/lib/systemd/system/v2ray.service"
    if [[ $? -ne 0 ]]; then
        colorEcho ${red} "Failed to remove V2Ray."
        exit
    else
        colorEcho ${green} "Removed V2Ray successfully."
        colorEcho ${green} "If necessary, please remove configuration file and log file manually."
        exit
    fi
}

checkUpdate(){
        echo "Checking for update."
        getVersion
        if [[ $? -eq 1 ]]; then
            colorEcho ${green} "Found new version ${NEW_VER} for V2Ray."
            exit 
        else 
            colorEcho ${green} "No new version."
            exit
        fi
}

main(){
    #helping information
    [[ "$HELP" == "1" ]] && Help
    [[ "$CHECK" == "1" ]] && checkUpdate
    [[ "$REMOVE" == "1" ]] && remove
    
    sysAcrh
    # extra local file
    if [[ $LOCAL_INSTALL -eq 1 ]]; then
        echo "Install V2Ray via local file"
        installSoftware unzip
        rm -rf /tmp/v2ray
        extra $LOCAL
        FILEVDIS=`ls /tmp/v2ray |grep v2ray-v |cut -d "-" -f4`
        SYSTEM=`ls /tmp/v2ray |grep v2ray-v |cut -d "-" -f3`
        if [[ ${SYSTEM} != "linux" ]]; then
            colorEcho $red "The local V2Ray can not be installed in linux."
            exit
        elif [[ ${FILEVDIS} != ${VDIS} ]]; then
            colorEcho $red "The local V2Ray can not be installed in ${ARCH} system."
            exit
        else
            NEW_VER=`ls /tmp/v2ray |grep v2ray-v |cut -d "-" -f2`
        fi
    else
        # dowload via network and extra
        installSoftware "curl"
        getVersion
        if [[ $? == 0 ]] && [[ "$FORCE" != "1" ]]; then
            colorEcho ${green} "Lastest version ${NEW_VER} is already installed."
            exit
        else
            colorEcho ${wblue} "Installing V2Ray ${NEW_VER} on ${ARCH}"
            downloadV2Ray
            installSoftware unzip
            extra ${ZIPFILE}
        fi
    fi 
    if pgrep "v2ray" > /dev/null ; then
        V2RAY_RUNNING=1
        stopV2ray
    fi
    installV2Ray
    installInitScrip
    if [[ ${V2RAY_RUNNING} -eq 1 ]];then
        colorEcho ${wblue} "Restarting V2Ray service."
        startV2ray

    fi
    colorEcho ${green} "V2Ray ${NEW_VER} is installed."
    rm -rf /tmp/v2ray
    return 0
}

main
