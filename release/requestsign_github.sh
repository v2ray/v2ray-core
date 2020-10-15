#!/usr/bin/env bash

export SIGN_VERSION=$(cat $GITHUB_EVENT_PATH| jq -r ".release.tag_name")

echo $SIGN_VERSION

$GITHUB_WORKSPACE/release/requestsign.sh
