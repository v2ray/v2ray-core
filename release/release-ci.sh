#!/bin/bash

set -x

apt-get update
apt-get -y install jq git file p7zip-full

function getattr() {
  curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/$2/attributes/$1
}

GITHUB_TOKEN=$(getattr "github_token" "project")
RELEASE_TAG=$(getattr "release_tag" "instance")
PRERELEASE=$(getattr "prerelease" "instance")
DOCKER_HUB_KEY=$(getattr "docker_hub_key" "project")
SIGN_KEY_PATH=$(getattr "sign_key_path" "project")
SIGN_KEY_PASS=$(getattr "sign_key_pass" "project")
VUSER=$(getattr "b_user" "project")

mkdir -p /v2ray/build

gsutil cp ${SIGN_KEY_PATH} /v2ray/build/sign_key.asc
echo ${SIGN_KEY_PASS} | gpg --passphrase-fd 0 --batch --import /v2ray/build/sign_key.asc

curl -L -o /v2ray/build/releases https://api.github.com/repos/v2ray/v2ray-core/releases

GO_INSTALL=golang.tar.gz
curl -L -o ${GO_INSTALL} https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
tar -C /usr/local -xzf ${GO_INSTALL}
export PATH=$PATH:/usr/local/go/bin

mkdir -p /v2ray/src
export GOPATH=/v2ray

go get -u v2ray.com/core/...
go get -u v2ray.com/ext/...

pushd $GOPATH/src/v2ray.com/core/
git checkout tags/${RELEASE_TAG}
popd

go install v2ray.com/ext/tools/build/vbuild

export TRAVIS_TAG=${RELEASE_TAG}
export GPG_SIGN_PASS=${SIGN_KEY_PASS}
export V_USER=${VUSER}

$GOPATH/bin/vbuild --os=windows --arch=x86 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=windows --arch=x64 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=macos --arch=x64 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=x86 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=x64 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=arm --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=arm64 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=mips64 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=mips64le --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=mips --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=mipsle --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=linux --arch=s390x --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=freebsd --arch=x86 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=freebsd --arch=amd64 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=openbsd --arch=x86 --zip --sign #--encrypt
$GOPATH/bin/vbuild --os=openbsd --arch=amd64 --zip --sign #--encrypt

RELBODY="https://www.v2ray.com/chapter_00/01_versions.html"
JSON_DATA=$(echo "{}" | jq -c ".tag_name=\"${RELEASE_TAG}\"")
JSON_DATA=$(echo ${JSON_DATA} | jq -c ".prerelease=${PRERELEASE}")
JSON_DATA=$(echo ${JSON_DATA} | jq -c ".body=${RELBODY}")
RELEASE_ID=$(curl --data "${JSON_DATA}" -H "Authorization: token ${GITHUB_TOKEN}" -X POST https://api.github.com/repos/v2ray/v2ray-core/releases | jq ".id")

function upload() {
  FILE=$1
  CTYPE=$(file -b --mime-type $FILE)
  curl -H "Authorization: token ${GITHUB_TOKEN}" -H "Content-Type: ${CTYPE}" --data-binary @$FILE "https://uploads.github.com/repos/v2ray/v2ray-core/releases/${RELEASE_ID}/assets?name=$(basename $FILE)"
}

upload $GOPATH/bin/v2ray-macos.zip
upload $GOPATH/bin/v2ray-windows-64.zip
upload $GOPATH/bin/v2ray-windows-32.zip
upload $GOPATH/bin/v2ray-linux-64.zip
upload $GOPATH/bin/v2ray-linux-32.zip
upload $GOPATH/bin/v2ray-linux-arm.zip
upload $GOPATH/bin/v2ray-linux-arm64.zip
upload $GOPATH/bin/v2ray-linux-mips64.zip
upload $GOPATH/bin/v2ray-linux-mips64le.zip
upload $GOPATH/bin/v2ray-linux-mips.zip
upload $GOPATH/bin/v2ray-linux-mipsle.zip
upload $GOPATH/bin/v2ray-linux-s390x.zip
upload $GOPATH/bin/v2ray-freebsd-64.zip
upload $GOPATH/bin/v2ray-freebsd-32.zip
upload $GOPATH/bin/v2ray-openbsd-64.zip
upload $GOPATH/bin/v2ray-openbsd-32.zip
upload $GOPATH/bin/metadata.txt

if [[ "${PRERELEASE}" == "false" ]]; then

gsutil cp $GOPATH/bin/v2ray-${RELEASE_TAG}-linux-64/v2ray gs://v2ray-docker/
gsutil cp $GOPATH/bin/v2ray-${RELEASE_TAG}-linux-64/v2ctl gs://v2ray-docker/
gsutil cp $GOPATH/bin/v2ray-${RELEASE_TAG}-linux-64/geoip.dat gs://v2ray-docker/
gsutil cp $GOPATH/bin/v2ray-${RELEASE_TAG}-linux-64/geosite.dat gs://v2ray-docker/

DOCKER_HUB_API=https://registry.hub.docker.com/u/v2ray/official/trigger/${DOCKER_HUB_KEY}/
curl -H "Content-Type: application/json" --data '{"build": true}' -X POST "${DOCKER_HUB_API}"

fi

shutdown -h now
