#!/bin/bash

VER=$1
MSG=$2

if [ -z "$MSG" ]; then
  MSG="Weekly Release"
fi

echo Creating a new release: $VER: $MSG

IFS="." read -a PARTS <<< "$VER"
MAJOR=${PARTS[0]}
MINOR=${PARTS[1]}
MINOR=$((MINOR+1))
VERN=${MAJOR}.${MINOR}

pushd $GOPATH/src/github.com/v2ray/v2ray-core
echo "Adding a new tag: " "v$VER"
git tag -s -a "v$VER" -m "$MSG"
sed -i '' "s/\(version *= *\"\).*\(\"\)/\1$VERN\2/g" core.go
echo "Commiting core.go (may not necessary)"
git commit core.go -S -m "Update version"
echo "Pushing changes"
git push --follow-tags
popd