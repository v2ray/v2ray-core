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
$GOPATH/bin/build --os=openbsd --arch=x86 --zip
$GOPATH/bin/build --os=openbsd --arch=amd64 --zip

INSTALL_DIR=_install

git clone "https://github.com/v2ray/install.git" ${INSTALL_DIR}

RELEASE_DIR=${INSTALL_DIR}/releases/${TRAVIS_TAG}
mkdir -p ${RELEASE_DIR}/
cp $GOPATH/bin/metadata.txt ${RELEASE_DIR}/
cp $GOPATH/bin/v2ray-*.zip ${RELEASE_DIR}/
echo ${TRAVIS_TAG} > ${INSTALL_DIR}/releases/latest.txt

cp $GOPATH/bin/v2ray-${TRAVIS_TAG}-linux-64/v2ray ${INSTALL_DIR}/docker/official/

pushd ${INSTALL_DIR}
git config user.name "V2Ray Auto Build"
git config user.email "admin@v2ray.com"
git add -A
git commit -m "Update for ${TRAVIS_TAG}"
git push "https://${GIT_KEY_INSTALL}@github.com/v2ray/install.git" master
popd

DOCKER_HUB_API=https://registry.hub.docker.com/u/v2ray/official/trigger/${DOCKER_HUB_KEY}/
curl -H "Content-Type: application/json" --data '{"build": true}' -X POST "${DOCKER_HUB_API}"
