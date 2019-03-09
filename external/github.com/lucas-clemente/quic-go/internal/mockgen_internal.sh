#!/bin/bash

# Mockgen refuses to generate mocks for internal packages.
# This script copies the internal directory and renames it to internalpackage.
# It also creages a public alias for private types.
# It then creates a mock for this public (alias) type.
# Afterwards, it corrects the import paths (replaces internalpackage back to internal).

TEMP_DIR=$(mktemp -d)
mkdir -p $TEMP_DIR/src/github.com/lucas-clemente/quic-go/internalpackage

# uppercase the name of the interface (only has an effect for private interfaces)
INTERFACE_NAME="$(tr '[:lower:]' '[:upper:]' <<< ${4:0:1})${4:1}"
PACKAGE_NAME=`echo $3 | sed 's/.*\///'`

cp -r $GOPATH/src/github.com/lucas-clemente/quic-go/internal/* $TEMP_DIR/src/github.com/lucas-clemente/quic-go/internalpackage
find $TEMP_DIR -type f -name "*.go" -exec sed -i '' 's/internal/internalpackage/g' {} \;

export GOPATH="$TEMP_DIR:$GOPATH"
PACKAGE_PATH=${3/internal/internalpackage}

# if we're mocking a private interface, we need to add a public alias
if [ "$INTERFACE_NAME" != "$4" ]; then
  # create a public alias for the interface, so that mockgen can process it
  echo -e "package $PACKAGE_NAME\n" > $TEMP_DIR/src/$PACKAGE_PATH/mockgen_interface.go
  echo "type $INTERFACE_NAME = $4" >> $TEMP_DIR/src/$PACKAGE_PATH/mockgen_interface.go
fi

mockgen -package $1 -self_package $1 -destination $2 $PACKAGE_PATH $INTERFACE_NAME
sed -i '' 's/internalpackage/internal/g' $2

# mockgen imports the package we're generating a mock for
sed -i '' "s/$1\.//g" $2
goimports -w $2

rm -r "$TEMP_DIR"
