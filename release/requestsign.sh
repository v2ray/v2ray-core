#!/usr/bin/env bash

RELEASE_DATA=$(curl --data "version=${SIGN_VERSION}" --data "password=${SIGN_SERVICE_PASSWORD}" -X POST "${SIGN_SERIVCE_URL}" )
echo $RELEASE_DATA
RELEASE_ID=$(echo $RELEASE_DATA| jq -r ".id")

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

curl "https://raw.githubusercontent.com/v2fly/Release/master/v2fly/${SIGN_VERSION}.Release" > Release
upload Release
