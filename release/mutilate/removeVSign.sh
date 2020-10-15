#!/usr/bin/env bash

export VROOT=$(dirname "${BASH_SOURCE[0]}")/../../

rm $VROOT/infra/control/verify.go

sed -i '/VSign/d' $VROOT/go.mod
