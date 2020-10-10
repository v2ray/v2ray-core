#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; exit 1' ERR

NOW=$(date '+%Y%m%d-%H%M%S')
TMP=$(mktemp -d)
SRCDIR=$(pwd)

CODENAME="user"
BUILDNAME=$NOW

cleanup() { rm -rf "$TMP"; }
trap cleanup INT TERM ERR

get_source() {
	echo ">>> Clone v2fly/v2ray-core repo..."
	git clone https://github.com/v2fly/v2ray-core.git
	cd v2ray-core
	go mod download
}

build_v2() {
	if [[ $nosource != 1 ]]; then
		cd ${SRCDIR}/v2ray-core
		local VERSIONTAG=$(git describe --abbrev=0 --tags)
	else
		echo ">>> Use current directory as WORKDIR"
		local VERSIONTAG=$(git describe --abbrev=0 --tags)
	fi

	LDFLAGS="-s -w -buildid= -X v2ray.codename=${CODENAME} -X v2ray.build=${BUILDNAME} -X v2ray.version=${VERSIONTAG}"

	echo ">>> Compile v2ray ..."
	env CGO_ENABLED=0 go build -o "$TMP"/v2ray"${EXESUFFIX}" -ldflags "$LDFLAGS" ./main
	if [[ $GOOS == "windows" ]]; then
		env CGO_ENABLED=0 go build -o "$TMP"/wv2ray"${EXESUFFIX}" -ldflags "-H windowsgui $LDFLAGS" ./main
	fi

	echo ">>> Compile v2ctl ..."
	env CGO_ENABLED=0 go build -o "$TMP"/v2ctl"${EXESUFFIX}" -tags confonly -ldflags "$LDFLAGS" ./infra/control/main
}

build_dat() {
	echo ">>> Download latest geoip..."
	curl -s -L -o "$TMP"/geoip.dat "https://github.com/v2fly/geoip/raw/release/geoip.dat"

	echo ">>> Download latest geosite..."
	curl -s -L -o "$TMP"/geosite.dat "https://github.com/v2fly/domain-list-community/raw/release/dlc.dat"
}

copyconf() {
	echo ">>> Copying config..."
	cd ./release/config
	if [[ $GOOS == "linux" ]]; then
		tar c --exclude "*.dat" . | tar x -C "$TMP"
	else
		tar c --exclude "*.dat" --exclude "systemd/**" . | tar x -C "$TMP"
	fi
}

packzip() {
	echo ">>> Generating zip package"
	cd "$TMP"
	local PKG=${SRCDIR}/v2ray-custom-${GOARCH}-${GOOS}-${PKGSUFFIX}${NOW}.zip
	zip -r "$PKG" .
	echo ">>> Generated: $(basename "$PKG") at $(dirname "$PKG")"
}

packtgz() {
	echo ">>> Generating tgz package"
	cd "$TMP"
	local PKG=${SRCDIR}/v2ray-custom-${GOARCH}-${GOOS}-${PKGSUFFIX}${NOW}.tar.gz
	tar cvfz "$PKG" .
	echo ">>> Generated: $(basename "$PKG") at $(dirname "$PKG")"
}

packtgzAbPath() {
	local ABPATH="$1"
	echo ">>> Generating tgz package at $ABPATH"
	cd "$TMP"
	tar cvfz "$ABPATH" .
	echo ">>> Generated: $ABPATH"
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
	386 | arm* | mips* | ppc64* | riscv64 | s390x)
		GOARCH=$arg
		;;
	windows)
		GOOS=$arg
		EXESUFFIX=.exe
		;;
	darwin | dragonfly | freebsd | openbsd)
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
	abpathtgz=*)
		pkg=${arg##abpathtgz=}
		;;
	codename=*)
		CODENAME=${arg##codename=}
		;;
	buildname=*)
		BUILDNAME=${arg##buildname=}
		;;
	esac
done

if [[ $nosource != 1 ]]; then
	get_source
fi

export GOOS GOARCH
echo "Build ARGS: GOOS=${GOOS} GOARCH=${GOARCH} CODENAME=${CODENAME} BUILDNAME=${BUILDNAME}"
echo "PKG ARGS: pkg=${pkg}"
build_v2

if [[ $nodat != 1 ]]; then
	build_dat
fi

if [[ $noconf != 1 ]]; then
	copyconf
fi

if [[ $pkg == "zip" ]]; then
	packzip
elif [[ $pkg == "tgz" ]]; then
	packtgz
else
	packtgzAbPath "$pkg"
fi

cleanup
