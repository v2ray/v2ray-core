#!/bin/bash

function git_not_installed {
  git --version 2>&1 >/dev/null
  GIT_IS_AVAILABLE=$?
  return $GIT_IS_AVAILABLE
}

if [ git_not_installed ]; then
  apt-get install git -y
fi


if [ -z "$GOPATH" ]; then
  curl -o go_latest.tar.gz https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
  tar -C /usr/local -xzf go_latest.tar.gz
  rm go_latest.tar.gz
  export PATH=$PATH:/usr/local/go/bin
  
  mkdir /v2ray
  export GOPATH=/v2ray
fi

go get github.com/v2ray/v2ray-core
go build -o $GOPATH/bin/v2ray -compiler gc github.com/v2ray/v2ray-core/release/server

