#!/bin/bash

# This file is accessible as https://install.direct/go.sh
# Original source is located at github.com/v2fly/v2ray-core/release/install-release.sh

# If not specify, default meaning of return value:
# 0: Success
# 1: System error
# 2: Application error
# 3: Network error

#######color code########
RED="31m"      # Error message
YELLOW="33m"   # Warning message
colorEcho(){
    echo -e "\033[${1}${@:2}\033[0m" 1>& 2
}

colorEcho ${RED} "WARN: This script will be obsolete soon, please switch to the fhs-install-v2ray project in time."
colorEcho ${YELLOW} "URL: https://github.com/v2fly/fhs-install-v2ray"
