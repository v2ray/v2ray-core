#!/bin/bash

GO_AMD64=https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
GO_X86=https://storage.googleapis.com/golang/go1.5.1.linux-386.tar.gz
ARCH=$(uname -m)
GO_CUR=${GO_AMD64}

if [ "$ARCH" == "i686" ] || [ "$ARCH" == "i386" ]; then
  GO_CUR=${GO_X86}
fi

function git_not_installed {
  git --version 2>&1 >/dev/null
  GIT_IS_AVAILABLE=$?
  return $GIT_IS_AVAILABLE
}

if [ git_not_installed ]; then
  apt-get install git -y
fi


if [ -z "$GOPATH" ]; then
  curl -o go_latest.tar.gz ${GO_CUR}
  tar -C /usr/local -xzf go_latest.tar.gz
  rm go_latest.tar.gz
  export PATH=$PATH:/usr/local/go/bin
  
  mkdir /v2ray
  export GOPATH=/v2ray
fi

go get -u github.com/v2ray/v2ray-core
rm $GOPATH/bin/build
go install github.com/v2ray/v2ray-core/tools/build
$GOPATH/bin/build
