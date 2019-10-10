#!/usr/bin/env bash
# Bash3 Boilerplate. Copyright (c) 2014, kvz.io

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; exit 1' ERR

# Set magic variables for current file & dir
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"
__root="$(cd "$(dirname "${__dir}")" && pwd)" # <-- change this as it depends on your app


NOW=$(date '+%Y%m%d-%H%M%S')
TMP=$(mktemp -d)

CODENAME="user"
BUILDNAME=$NOW
GOPATH=$(go env GOPATH)

cleanup () { rm -rf $TMP; }
trap cleanup INT TERM ERR

get_source() {
	echo ">>> Getting v2ray sources ..."
	go get -insecure -v -t v2ray.com/core/...
}

build_v2() {
	pushd $GOPATH/src/v2ray.com/core
	echo ">>> Update source code name ..."
	sed -i "s/^[ \t]\+codename.\+$/\tcodename = \"${CODENAME}\"/;s/^[ \t]\+build.\+$/\tbuild = \"${BUILDNAME}\"/;" core.go

	echo ">>> Compile v2ray ..."
	pushd $GOPATH/src/v2ray.com/core/main
	env CGO_ENABLED=0 go build -o $TMP/v2ray${EXESUFFIX} -ldflags "-s -w"
	if [[ $GOOS == "windows" ]];then
	  env CGO_ENABLED=0 go build -o $TMP/wv2ray${EXESUFFIX} -ldflags "-s -w -H windowsgui"
	fi
	popd

	git checkout -- core.go
	popd

	echo ">>> Compile v2ctl ..."
	pushd $GOPATH/src/v2ray.com/core/infra/control/main
	env CGO_ENABLED=0 go build -o $TMP/v2ctl${EXESUFFIX} -tags confonly -ldflags "-s -w"
	popd
}

build_dat() {
	echo ">>> Downloading newest geoip ..."
	wget -qO - https://api.github.com/repos/v2ray/geoip/releases/latest \
	| grep browser_download_url | cut -d '"' -f 4 \
	| wget -i - -O $TMP/geoip.dat

	echo ">>> Downloading newest geosite ..."
	wget -qO - https://api.github.com/repos/v2ray/domain-list-community/releases/latest \
	| grep browser_download_url | cut -d '"' -f 4 \
	| wget -i - -O $TMP/geosite.dat
}

copyconf() {
	echo ">>> Copying config..."
	pushd $GOPATH/src/v2ray.com/core/release/config
	tar c --exclude "*.dat" . | tar x -C $TMP
}

packzip() {
	echo ">>> Generating zip package"
	pushd $TMP
	local PKG=${__dir}/v2ray-custom-${GOARCH}-${GOOS}-${PKGSUFFIX}${NOW}.zip
	zip -r $PKG .
	echo ">>> Generated: $(basename $PKG)"
}

packtgz() {
	echo ">>> Generating tgz package"
	pushd $TMP
	local PKG=${__dir}/v2ray-custom-${GOARCH}-${GOOS}-${PKGSUFFIX}${NOW}.tar.gz
	tar cvfz $PKG .
	echo ">>> Generated: $(basename $PKG)"
}


pkg=zip
nosource=0
nodat=0
noconf=0
GOOS=linux
GOARCH=amd64
EXESUFFIX=
PKGSUFFIX=

for arg in "$@"; do
case $arg in
	arm*)
		GOARCH=$arg
		;;
	mips*)
		GOARCH=$arg
		;;
	386)
		GOARCH=386
		;;
	windows)
		GOOS=windows
		EXESUFFIX=.exe
		;;
	darwin)
		GOOS=$arg
		;;
	nodat)
		nodat=1
		PKGSUFFIX=${PKGSUFFIX}nodat-
		;;
	noconf)
		noconf=1
		;;
	nosource)
		nosource=1
		;;
	tgz)
		pkg=tgz
		;;
esac
done

if [[ $nosource != 1 ]]; then
  get_source	
fi

export GOOS GOARCH
build_v2

if [[ $nodat != 1 ]]; then
  build_dat
fi

if [[ $noconf != 1 ]]; then
  copyconf 
fi

if [[ $pkg == "zip" ]]; then
  packzip
fi

if [[ $pkg == "tgz" ]]; then
  packtgz
fi

cleanup

