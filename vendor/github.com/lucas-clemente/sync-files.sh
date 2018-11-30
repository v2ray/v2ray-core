#!/bin/bash

rsync -rv "$GOPATH/src/github.com/lucas-clemente/quic-go/" "$GOPATH/src/v2ray.com/core/vendor/github.com/lucas-clemente/quic-go/"
find . -name "*_test.go" -delete
rm -rf ./quic-go/\.*
rm -rf ./quic-go/benchmark
rm -rf ./quic-go/docs
rm -rf ./quic-go/example
rm -rf ./quic-go/h2quic
rm -rf ./quic-go/integrationtests
rm -rf ./quic-go/vendor/golang\.org/
rm ./quic-go/vendor/vendor.json
