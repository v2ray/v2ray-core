#!/bin/bash

pushd $GOPATH/src
for DIR in $(find ./v2ray.com/core -type d -not -path "*.git*"); do
  TEST_FILES=($DIR/*.proto)
  #echo ${TEST_FILES}
  if [ -f ${TEST_FILES[0]} ]; then
    protoc --proto_path=. --go_out=. $DIR/*.proto
  fi
done
popd