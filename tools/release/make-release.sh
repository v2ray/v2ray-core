#!/bin/bash

VER=$1
PRE=$2
PROJECT=$3

if [ -z "$PRE" ]; then
  PRE="true"
fi

if [ -z "$PROJECT" ]; then
  echo "Project not specified. Exiting..."
  exit 0
fi

echo Creating a new release: $VER: $MSG

IFS="." read -a PARTS <<< "$VER"
MAJOR=${PARTS[0]}
MINOR=${PARTS[1]}
MINOR=$((MINOR+1))
VERN=${MAJOR}.${MINOR}

pushd $GOPATH/src/v2ray.com/core
echo "Adding a new tag: " "v$VER"
git tag -s -a "v$VER" -m "Version ${VER}"
sed -i '' "s/\(version *= *\"\).*\(\"\)/\1$VERN\2/g" core.go
echo "Commiting core.go (may not necessary)"
git commit core.go -S -m "Update version"
echo "Pushing changes"
git push --follow-tags
popd

echo "Launching build machine."
DIR="$(dirname "$0")"
gcloud compute instances create "build-upload" \
    --machine-type=n1-highcpu-2 \
    --metadata=release_tag=v${VER},prerelease=${PRE} \
    --metadata-from-file=startup-script=${DIR}/release-ci.sh \
    --zone=us-west1-a \
    --project ${PROJECT}
