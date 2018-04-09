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
RED="31m"
GREEN="32m"
YELLOW="33m"
BLUE="36m"


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
        --remove)
        REMOVE="1"
        ;;
        --version)
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
    COLOR=$1
    echo -e "\033[${COLOR}${@:2}\033[0m"
}

sysArch(){
    ARCH=$(uname -m)
    if [[ "$ARCH" == "i686" ]] || [[ "$ARCH" == "i386" ]]; then
        VDIS="32"
    elif [[ "$ARCH" == *"armv7"* ]] || [[ "$ARCH" == "armv6l" ]]; then
        VDIS="arm"
    elif [[ "$ARCH" == *"armv8"* ]] || [[ "$ARCH" == "aarch64" ]]; then
        VDIS="arm64"
    elif [[ "$ARCH" == *"mips64le"* ]]; then
        VDIS="mips64le"
    elif [[ "$ARCH" == *"mips64"* ]]; then
        VDIS="mips64"
    elif [[ "$ARCH" == *"mipsle"* ]]; then
        VDIS="mipsle"
    elif [[ "$ARCH" == *"mips"* ]]; then
        VDIS="mips"
    elif [[ "$ARCH" == *"s390x"* ]]; then
        VDIS="s390x"
    fi
}

downloadV2Ray(){
    rm -rf /tmp/v2ray
    mkdir -p /tmp/v2ray
    colorEcho ${BLUE} "Downloading V2Ray."
    DOWNLOAD_LINK="https://github.com/v2ray/v2ray-core/releases/download/v${NEW_VER}/v2ray-linux-${VDIS}.zip"
    curl ${PROXY} -L -H "Cache-Control: no-cache" -o ${ZIPFILE} ${DOWNLOAD_LINK}
    if [ $? -ne 0 ];then
        colorEcho ${RED} "Failed to download! Please check your network or try again."
        return 3
    fi
}

installSoftware(){
    COMPONENT=$1
    if [[ -n `command -v $COMPONENT` ]]; then
        return
    fi

    getPMT
    if [[ $? -eq 1 ]]; then
        colorEcho $YELLOW "The system package manager tool isn't APT or YUM, please install ${COMPONENT} manually."
        return 2 
    fi
    colorEcho $GREEN "Installing $COMPONENT" 
    if [[ $SOFTWARE_UPDATED -eq 0 ]]; then
        colorEcho ${BLUE} "Updating software repo"
        $CMD_UPDATE      
        SOFTWARE_UPDATED=1
    fi

    colorEcho ${BLUE} "Installing ${COMPONENT}"
    $CMD_INSTALL $COMPONENT
    if [[ $? -ne 0 ]]; then
        colorEcho ${RED} "Install ${COMPONENT} fail, please install it manually."
        return 2
    fi
}

# return 1: not apt or yum
getPMT(){
    if [[ -n `command -v apt-get` ]];then
        CMD_INSTALL="apt-get -y -qq install"
        CMD_UPDATE="apt-get -qq update"
    elif [[ -n `command -v yum` ]]; then
        CMD_INSTALL="yum -y -q install"
        CMD_UPDATE="yum -q makecache"
    else
        return 1
    fi
    return 0
}


extract(){
    colorEcho ${BLUE}"Extracting V2Ray package to /tmp/v2ray."
    mkdir -p /tmp/v2ray
    unzip $1 -d "/tmp/v2ray/"
    if [[ $? -ne 0 ]]; then
        colorEcho ${RED} "Extracting V2Ray failed!"
        return 2
    fi
}


# 1: new V2Ray. 0: no. 2: not installed. 3: check failed.
getVersion(){
    if [[ -n "$VERSION" ]]; then
        NEW_VER="$VERSION"
        return 1
    else
        VER=`/usr/bin/v2ray/v2ray -version 2>/dev/null`
        RETVAL="$?"
        CUR_VER=`echo $VER | head -n 1 | cut -d " " -f2`
        TAG_URL="https://api.github.com/repos/v2ray/v2ray-core/releases/latest"
        NEW_VER=`curl ${PROXY} -s ${TAG_URL} --connect-timeout 10| grep 'tag_name' | cut -d\" -f4`
        if [[ $? -ne 0 ]] || [[ $NEW_VER == "" ]]; then
            colorEcho ${RED} "Network error! Please check your network or try again."
            return 3
        elif [[ $RETVAL -ne 0 ]];then
            return 2
        elif [[ "$NEW_VER" != "$CUR_VER" ]];then
            return 1
        fi
        return 0
    fi
}

stopV2ray(){
    SYSTEMCTL_CMD=$(command -v systemctl)
    SERVICE_CMD=$(command -v service)

    colorEcho ${BLUE} "Shutting down V2Ray service."
    if [[ -n "${SYSTEMCTL_CMD}" ]] || [[ -f "/lib/systemd/system/v2ray.service" ]] || [[ -f "/etc/systemd/system/v2ray.service" ]]; then
        ${SYSTEMCTL_CMD} stop v2ray
    elif [[ -n "${SERVICE_CMD}" ]] || [[ -f "/etc/init.d/v2ray" ]]; then
        ${SERVICE_CMD} v2ray stop
    fi
    if [[ $? -ne 0 ]]; then
        colorEcho ${RED} "Failed to shutdown V2Ray service."
        return 2
    fi
}

startV2ray(){
    SYSTEMCTL_CMD=$(command -v systemctl)
    SERVICE_CMD=$(command -v service)

    if [ -n "${SYSTEMCTL_CMD}" ] && [ -f "/lib/systemd/system/v2ray.service" ]; then
        ${SYSTEMCTL_CMD} start v2ray
    elif [ -n "${SYSTEMCTL_CMD}" ] && [ -f "/etc/systemd/system/v2ray.service" ]; then
        ${SYSTEMCTL_CMD} start v2ray
    elif [ -n "${SERVICE_CMD}" ] && [ -f "/etc/init.d/v2ray" ]; then
        ${SERVICE_CMD} v2ray start
    fi
    if [[ $? -ne 0 ]]; then
        colorEcho ${RED} "Failed to start V2Ray service."
        return 2
    fi
}

copyFile() {
    NAME=$1
    ERROR=`cp "/tmp/v2ray/v2ray-v${NEW_VER}-linux-${VDIS}/${NAME}" "/usr/bin/v2ray/${NAME}" 2>&1`
    if [[ $? -ne 0 ]]; then
        colorEcho ${YELLOW} "${ERROR}"
        return 2
    fi
}

makeExecutable() {
    chmod +x "/usr/bin/v2ray/$1"
}

installV2Ray(){
    # Install V2Ray binary to /usr/bin/v2ray
    mkdir -p /usr/bin/v2ray
    copyFile v2ray || return $?
    makeExecutable v2ray
    copyFile v2ctl
    makeExecutable v2ctl
    copyFile geoip.dat
    copyFile geosite.dat

    # Install V2Ray server config to /etc/v2ray
    if [[ ! -f "/etc/v2ray/config.json" ]]; then
        mkdir -p /etc/v2ray
        cp "/tmp/v2ray/v2ray-${NEW_VER}-linux-${VDIS}/vpoint_vmess_freedom.json" "/etc/v2ray/config.json"
        if [[ $? -ne 0 ]]; then
            colorEcho ${YELLOW} "Create V2Ray configuration file error, pleases create it manually."
            return 1
        fi
        let PORT=$RANDOM+10000
        UUID=$(cat /proc/sys/kernel/random/uuid)

        sed -i "s/10086/${PORT}/g" "/etc/v2ray/config.json"
        sed -i "s/23ad6b10-8d1a-40f7-8ad0-e3e35cd38297/${UUID}/g" "/etc/v2ray/config.json"

        colorEcho ${GREEN} "PORT:${PORT}"
        colorEcho ${GREEN} "UUID:${UUID}"
        mkdir -p /var/log/v2ray
    fi
}


installInitScript(){
    SYSTEMCTL_CMD=$(command -v systemctl)
    SERVICE_CMD=$(command -v service)

    if [[ -n "${SYSTEMCTL_CMD}" ]];then
        if [[ ! -f "/etc/systemd/system/v2ray.service" ]]; then
            if [[ ! -f "/lib/systemd/system/v2ray.service" ]]; then
                cp "/tmp/v2ray/v2ray-${NEW_VER}-linux-${VDIS}/systemd/v2ray.service" "/etc/systemd/system/"
                systemctl enable v2ray.service
            fi
        fi
    elif [[ -n "${SERVICE_CMD}" ]] && [[ ! -f "/etc/init.d/v2ray" ]]; then
        installSoftware "daemon" || return $?
        cp "/tmp/v2ray/v2ray-${NEW_VER}-linux-${VDIS}/systemv/v2ray" "/etc/init.d/v2ray"
        chmod +x "/etc/init.d/v2ray"
        update-rc.d v2ray defaults
    fi
}

Help(){
    echo "./install-release.sh [-h] [-c] [-p proxy] [-f] [--version vx.y.z] [-l file]"
    echo "  -h, --help            Show help"
    echo "  -p, --proxy           To download through a proxy server, use -p socks5://127.0.0.1:1080 or -p http://127.0.0.1:3128 etc"
    echo "  -f, --force           Force install"
    echo "      --version         Install a particular version"
    echo "  -l, --local           Install from a local file"
    echo "      --remove          Remove installed V2Ray"
    echo "  -c, --check           Check for update"
}

remove(){
    SYSTEMCTL_CMD=$(command -v systemctl)
    SERVICE_CMD=$(command -v service)
    if [[ -n "${SYSTEMCTL_CMD}" ]] && [[ -f "/etc/systemd/system/v2ray.service" ]];then
        if pgrep "v2ray" > /dev/null ; then
            stopV2ray
        fi
        systemctl disable v2ray.service
        rm -rf "/usr/bin/v2ray" "/etc/systemd/system/v2ray.service"
        if [[ $? -ne 0 ]]; then
            colorEcho ${RED} "Failed to remove V2Ray."
            return 0
        else
            colorEcho ${GREEN} "Removed V2Ray successfully."
            colorEcho ${GREEN} "If necessary, please remove configuration file and log file manually."
            return
        fi
    elif [[ -n "${SYSTEMCTL_CMD}" ]] && [[ -f "/lib/systemd/system/v2ray.service" ]];then
        if pgrep "v2ray" > /dev/null ; then
            stopV2ray
        fi
        systemctl disable v2ray.service
        rm -rf "/usr/bin/v2ray" "/lib/systemd/system/v2ray.service"
        if [[ $? -ne 0 ]]; then
            colorEcho ${RED} "Failed to remove V2Ray."
            return 0
        else
            colorEcho ${GREEN} "Removed V2Ray successfully."
            colorEcho ${GREEN} "If necessary, please remove configuration file and log file manually."
            return
        fi
    elif [[ -n "${SERVICE_CMD}" ]] && [[ -f "/etc/init.d/v2ray" ]]; then
        if pgrep "v2ray" > /dev/null ; then
            stopV2ray
        fi
        rm -rf "/usr/bin/v2ray" "/etc/init.d/v2ray"
        if [[ $? -ne 0 ]]; then
            colorEcho ${RED} "Failed to remove V2Ray."
            return 0
        else
            colorEcho ${GREEN} "Removed V2Ray successfully."
            colorEcho ${GREEN} "If necessary, please remove configuration file and log file manually."
            return
        fi       
    else
        colorEcho ${GREEN} "V2Ray not found."
        return 0
    fi
}

checkUpdate(){
    echo "Checking for update."
    VERSION=""
    getVersion
    RETVAL="$?"
    if [[ $RETVAL -eq 1 ]]; then
        colorEcho ${GREEN} "Found new version ${NEW_VER} for V2Ray.(Current version:$CUR_VER)"
    elif [[ $RETVAL -eq 0 ]]; then
        colorEcho ${GREEN} "No new version. Current version is ${NEW_VER}."
    elif [[ $RETVAL -eq 2 ]]; then
        colorEcho ${RED} "No V2Ray installed."
        colorEcho ${GREEN} "The newest version for V2Ray is ${NEW_VER}."
    fi
    return 0
}

main(){
    #helping information
    [[ "$HELP" == "1" ]] && Help && return
    [[ "$CHECK" == "1" ]] && checkUpdate && return
    [[ "$REMOVE" == "1" ]] && remove && return
    
    sysArch
    # extract local file
    if [[ $LOCAL_INSTALL -eq 1 ]]; then
        echo "Install V2Ray via local file"
        installSoftware unzip || return $?
        rm -rf /tmp/v2ray
        extract $LOCAL || return $?
        FILEVDIS=`ls /tmp/v2ray |grep v2ray-v |cut -d "-" -f4`
        SYSTEM=`ls /tmp/v2ray |grep v2ray-v |cut -d "-" -f3`
        if [[ ${SYSTEM} != "linux" ]]; then
            colorEcho $RED "The local V2Ray can not be installed in linux."
            return 1
        elif [[ ${FILEVDIS} != ${VDIS} ]]; then
            colorEcho $RED "The local V2Ray can not be installed in ${ARCH} system."
            return 1
        else
            NEW_VER=`ls /tmp/v2ray |grep v2ray-v |cut -d "-" -f2`
        fi
    else
        # download via network and extract
        installSoftware "curl" || return $?
        getVersion
        RETVAL="$?"
        if [[ $RETVAL == 0 ]] && [[ "$FORCE" != "1" ]]; then
            colorEcho ${GREEN} "Latest version ${NEW_VER} is already installed."
            return
        elif [[ $RETVAL == 3 ]]; then
            return 3
        elif [[ $RETVAL == 1 ]]; then
            colorEcho ${GREEN} "Found new version ${NEW_VER} for V2Ray.(Current version:$CUR_VER)"
        fi
    colorEcho ${BLUE} "Installing V2Ray ${NEW_VER} on ${ARCH}"
    downloadV2Ray || return $?
    installSoftware unzip || return $?
    extract ${ZIPFILE} || return $?
    fi 
    if pgrep "v2ray" > /dev/null ; then
        V2RAY_RUNNING=1
        stopV2ray
    fi
    installV2Ray || return $?
    installInitScript || return $?
    if [[ ${V2RAY_RUNNING} -eq 1 ]];then
        colorEcho ${BLUE} "Restarting V2Ray service."
        startV2ray || return $?
    fi
    colorEcho ${GREEN} "V2Ray ${NEW_VER} is installed."
    rm -rf /tmp/v2ray
}

main
