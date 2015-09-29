#!/bin/bash

GIT_REV=$(git rev-parse HEAD)
VERSION=$(git name-rev --tags --name-only $GIT_REV)

LD_FLAGS="-s"

if [ "$VERSION" != "undefined" ]; then
  VERSION=${VERSION%^0}
  TODAY="$(date -u +%Y%m%d)"
  LD_FLAGS="${LD_FLAGS} -X github.com/v2ray/v2ray-core.version=${VERSION} -X github.com/v2ray/v2ray-core.build=${TODAY}"
else
  VERSION="custom"
fi

BIN_PATH=$GOPATH/bin
mkdir -p $BIN_PATH

function build {
  local GOOS=$1
  local GOARCH=$2
  local SUFFIX=$3
  local EXT=$4
  
  local REL_PATH=$BIN_PATH/v2ray_${VERSION}${SUFFIX}
  local TARGET=$REL_PATH/v2ray${EXT}
  if [ -d "$REL_PATH" ]; then
    rm -rf "$REL_PATH"
  fi
  mkdir -p $REL_PATH/config
  cp -R $GOPATH/src/github.com/v2ray/v2ray-core/release/config/* $REL_PATH/config/
  GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${TARGET} -compiler gc -ldflags "${LD_FLAGS}" github.com/v2ray/v2ray-core/release/server
  
  ZIP_FILE=$BIN_PATH/v2ray${SUFFIX}.zip
  if [ -f $ZIP_FILE ]; then
    rm -f $ZIP_FILE
  fi
  
  pushd $BIN_PATH
  zip -r $ZIP_FILE ./v2ray_${VERSION}${SUFFIX}/*
  popd
}

build "darwin" "amd64" "-macos" "-macos"
build "windows" "amd64" "-windows-64" "-windows-64.exe"
build "windows" "386" "-windows-32" "-windows-32.exe"
build "linux" "amd64" "-linux-64" "-linux-64"
build "linux" "386" "-linux-32" "-linux-32"
build "linux" "arm" "-armv6" "-armv6"
