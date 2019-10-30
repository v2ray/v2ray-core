#!/bin/bash

CONST_refs="refs"

TRIGGER_REASON_A=${TRIGGER_REASON:0:${#CONST_refs}}

if [ $TRIGGER_REASON_A != $CONST_refs ]
then
  echo "not a tag: $TRIGGER_REASON_A"
  exit
fi

CONST_refsB="refs/tags/"

TRIGGER_REASON_B=${TRIGGER_REASON:0:${#CONST_refsB}}

if [ $TRIGGER_REASON_B != $CONST_refsB ]
then
  echo "not a tag (B)"
  exit
fi


GITHUB_RELEASE_TAG=${TRIGGER_REASON:${#CONST_refsB}:10}

echo ${GITHUB_RELEASE_TAG}


RELEASE_DATA=$(curl -H "Authorization: token ${GITHUB_TOKEN}" -X GET https://api.github.com/repos/v2fly/v2ray-core/releases/tags/${GITHUB_RELEASE_TAG})
echo $RELEASE_DATA
RELEASE_ID=$(echo $RELEASE_DATA| jq ".id")

function uploadfile() {
  FILE=$1
  CTYPE=$(file -b --mime-type $FILE)

  sleep 1
  curl -H "Authorization: token ${GITHUB_TOKEN}" -H "Content-Type: ${CTYPE}" --data-binary @$FILE "https://uploads.github.com/repos/v2fly/v2ray-core/releases/${RELEASE_ID}/assets?name=$(basename $FILE)"
  sleep 1
}

function upload() {
  FILE=$1
  DGST=$1.dgst
  openssl dgst -md5 $FILE | sed 's/([^)]*)//g' >> $DGST
  openssl dgst -sha1 $FILE | sed 's/([^)]*)//g' >> $DGST
  openssl dgst -sha256 $FILE | sed 's/([^)]*)//g' >> $DGST
  openssl dgst -sha512 $FILE | sed 's/([^)]*)//g' >> $DGST
  uploadfile $FILE
  uploadfile $DGST
}

ART_ROOT=$GOPATH/src/v2ray.com/core/bazel-bin/release

upload ${ART_ROOT}/v2ray-macos.zip
upload ${ART_ROOT}/v2ray-windows-64.zip
upload ${ART_ROOT}/v2ray-windows-32.zip
upload ${ART_ROOT}/v2ray-linux-64.zip
upload ${ART_ROOT}/v2ray-linux-32.zip
upload ${ART_ROOT}/v2ray-linux-arm.zip
upload ${ART_ROOT}/v2ray-linux-arm64.zip
upload ${ART_ROOT}/v2ray-linux-mips64.zip
upload ${ART_ROOT}/v2ray-linux-mips64le.zip
upload ${ART_ROOT}/v2ray-linux-mips.zip
upload ${ART_ROOT}/v2ray-linux-mipsle.zip
upload ${ART_ROOT}/v2ray-linux-ppc64.zip
upload ${ART_ROOT}/v2ray-linux-ppc64le.zip
upload ${ART_ROOT}/v2ray-linux-s390x.zip
upload ${ART_ROOT}/v2ray-freebsd-64.zip
upload ${ART_ROOT}/v2ray-freebsd-32.zip
upload ${ART_ROOT}/v2ray-openbsd-64.zip
upload ${ART_ROOT}/v2ray-openbsd-32.zip
upload ${ART_ROOT}/v2ray-dragonfly-64.zip
