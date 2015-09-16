#!/bin/bash

VERSION=$(sed -n 's/.*Version = \"\([^"]*\)\"*/\1/p' $GOPATH/src/github.com/v2ray/v2ray-core/core.go)

REL_PATH=$GOPATH/bin/$VERSION
if [ -d "$REL_PATH" ]; then
  rm -rf "$REL_PATH"
fi

mkdir -p $REL_PATH
mkdir -p $REL_PATH/config

cp -R $GOPATH/src/github.com/v2ray/v2ray-core/release/config/* $REL_PATH/config/

function build {
  local GOOS=$1
  local GOARCH=$2
  local EXT=$3
  local TARGET=$REL_PATH/v2ray${EXT}
  GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${TARGET} -compiler gc github.com/v2ray/v2ray-core/release/server 
}

build "darwin" "amd64" "-macos"
build "windows" "amd64" "-windows-64.exe"
build "linux" "amd64" "-linux-64"
build "linux" "386" "-linux-32"

ZIP_FILE=$GOPATH/bin/v2ray-$VERSION.zip
if [ -f $ZIP_FILE ]; then
  rm -f $ZIP_FILE
fi

pushd $REL_PATH
zip -r $GOPATH/bin/v2ray-$VERSION.zip *
popd
