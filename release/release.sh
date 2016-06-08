#!/bin/bash

go install github.com/v2ray/v2ray-core/tools/build

$GOPATH/bin/build --os=windows --arch=x86 --zip
$GOPATH/bin/build --os=windows --arch=x64 --zip
$GOPATH/bin/build --os=macos --arch=x64 --zip
$GOPATH/bin/build --os=linux --arch=x86 --zip
$GOPATH/bin/build --os=linux --arch=x64 --zip
$GOPATH/bin/build --os=linux --arch=arm --zip
$GOPATH/bin/build --os=linux --arch=arm64 --zip
$GOPATH/bin/build --os=linux --arch=mips64 --zip
