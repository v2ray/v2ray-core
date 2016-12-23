#!/bin/bash

function detect_protoc() {
  SYS_LOC=$(which protoc)
  if [ -n "${SYS_LOC}" ]; then
    echo ${SYS_LOC}
    return
  fi

  if [[ "$OSTYPE" == "linux"* ]]; then
    echo $GOPATH/src/v2ray.com/core/.dev/protoc/linux/protoc
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo $GOPATH/src/v2ray.com/core/.dev/protoc/linux/protoc
  fi
}

PROTOC=$(detect_protoc)

# Update Golang proto compiler
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

pushd $GOPATH/src
for DIR in $(find ./v2ray.com/core -type d -not -path "*.git*"); do
  TEST_FILES=($DIR/*.proto)
  #echo ${TEST_FILES}
  if [ -f ${TEST_FILES[0]} ]; then
    ${PROTOC} --proto_path=. --go_out=. $DIR/*.proto
  fi
done
popd