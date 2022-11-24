#!/usr/bin/env bash

CONST_refs="refs"

TRIGGER_REASON_A=${TRIGGER_REASON:0:${#CONST_refs}}

if [ $TRIGGER_REASON_A != $CONST_refs ]; then
  echo "not a tag: $TRIGGER_REASON_A"
  exit
fi

CONST_refsB="refs/tags/"

TRIGGER_REASON_B=${TRIGGER_REASON:0:${#CONST_refsB}}

if [ $TRIGGER_REASON_B != $CONST_refsB ]; then
  echo "not a tag (B)"
  exit
fi

GITHUB_RELEASE_TAG=${TRIGGER_REASON:${#CONST_refsB}:25}

echo ${GITHUB_RELEASE_TAG}

RELEASE_DATA=$(curl -H "Authorization: token ${GITHUB_TOKEN}" -X GET https://api.github.com/repos/v2fly/v2ray-core/releases/tags/${GITHUB_RELEASE_TAG})
echo $RELEASE_DATA
RELEASE_ID=$(echo $RELEASE_DATA | jq ".id")

echo $RELEASE_ID

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
  openssl dgst -md5 $FILE | sed 's/([^)]*)//g' >>$DGST
  openssl dgst -sha1 $FILE | sed 's/([^)]*)//g' >>$DGST
  openssl dgst -sha256 $FILE | sed 's/([^)]*)//g' >>$DGST
  openssl dgst -sha512 $FILE | sed 's/([^)]*)//g' >>$DGST
  uploadfile $FILE
  uploadfile $DGST
}

ART_ROOT=${WORKDIR}/bazel-bin/release

pushd ${ART_ROOT}
{
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen version "${GITHUB_RELEASE_TAG}"
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen project "v2fly"
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-macos-64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-windows-64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-windows-32.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-windows-arm32-v7a.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-32.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-arm64-v8a.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-arm32-v7a.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-arm32-v6.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-arm32-v5.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-mips64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-mips64le.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-mips32.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-mips32le.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-ppc64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-ppc64le.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-riscv64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-linux-s390x.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-freebsd-64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-freebsd-32.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-openbsd-64.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-openbsd-32.zip
  go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen file v2ray-dragonfly-64.zip
} >Release.unsigned.unsorted
go run github.com/xiaokangwang/V2BuildAssist/v2buildutil gen sort <Release.unsigned.unsorted >Release.unsigned

popd

upload ${ART_ROOT}/v2ray-macos-64.zip
upload ${ART_ROOT}/v2ray-windows-64.zip
upload ${ART_ROOT}/v2ray-windows-32.zip
upload ${ART_ROOT}/v2ray-windows-arm32-v7a.zip
upload ${ART_ROOT}/v2ray-linux-64.zip
upload ${ART_ROOT}/v2ray-linux-32.zip
upload ${ART_ROOT}/v2ray-linux-arm64-v8a.zip
upload ${ART_ROOT}/v2ray-linux-arm32-v7a.zip
upload ${ART_ROOT}/v2ray-linux-arm32-v6.zip
upload ${ART_ROOT}/v2ray-linux-arm32-v5.zip
upload ${ART_ROOT}/v2ray-linux-mips64.zip
upload ${ART_ROOT}/v2ray-linux-mips64le.zip
upload ${ART_ROOT}/v2ray-linux-mips32.zip
upload ${ART_ROOT}/v2ray-linux-mips32le.zip
upload ${ART_ROOT}/v2ray-linux-ppc64.zip
upload ${ART_ROOT}/v2ray-linux-ppc64le.zip
upload ${ART_ROOT}/v2ray-linux-riscv64.zip
upload ${ART_ROOT}/v2ray-linux-s390x.zip
upload ${ART_ROOT}/v2ray-freebsd-64.zip
upload ${ART_ROOT}/v2ray-freebsd-32.zip
upload ${ART_ROOT}/v2ray-openbsd-64.zip
upload ${ART_ROOT}/v2ray-openbsd-32.zip
upload ${ART_ROOT}/v2ray-dragonfly-64.zip
upload ${ART_ROOT}/Release.unsigned
