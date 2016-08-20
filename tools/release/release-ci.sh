#!/bin/bash

go install v2ray.com/core/tools/build

$GOPATH/bin/build --os=windows --arch=x86 --zip
$GOPATH/bin/build --os=windows --arch=x64 --zip
$GOPATH/bin/build --os=macos --arch=x64 --zip
$GOPATH/bin/build --os=linux --arch=x86 --zip
$GOPATH/bin/build --os=linux --arch=x64 --zip
$GOPATH/bin/build --os=linux --arch=arm --zip
$GOPATH/bin/build --os=linux --arch=arm64 --zip
$GOPATH/bin/build --os=linux --arch=mips64 --zip
$GOPATH/bin/build --os=freebsd --arch=x86 --zip
$GOPATH/bin/build --os=freebsd --arch=amd64 --zip