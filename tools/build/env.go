package main

import (
	"strings"
)

type GoOS string

const (
	Windows   = GoOS("windows")
	MacOS     = GoOS("darwin")
	Linux     = GoOS("linux")
	UnknownOS = GoOS("unknown")
)

type GoArch string

const (
	X86         = GoArch("386")
	Amd64       = GoArch("amd64")
	Arm         = GoArch("arm")
	Arm64       = GoArch("arm64")
	UnknownArch = GoArch("unknown")
)

func parseOS(rawOS string) GoOS {
	osStr := strings.ToLower(rawOS)
	if osStr == "windows" || osStr == "win" {
		return Windows
	}
	if osStr == "darwin" || osStr == "mac" || osStr == "macos" || osStr == "osx" {
		return MacOS
	}
	if osStr == "linux" || osStr == "debian" || osStr == "ubuntu" || osStr == "redhat" || osStr == "centos" {
		return Linux
	}
	return UnknownOS
}

func parseArch(rawArch string) GoArch {
	archStr := strings.ToLower(rawArch)
	if archStr == "x86" || archStr == "386" || archStr == "i386" {
		return X86
	}
	if archStr == "amd64" || archStr == "x86-64" || archStr == "x64" {
		return Amd64
	}
	if archStr == "arm" {
		return Arm
	}
	if archStr == "arm64" {
		return Arm64
	}
	return UnknownArch
}

func getSuffix(os GoOS, arch GoArch) string {
	suffix := "-custom"
	switch os {
	case Windows:
		switch arch {
		case X86:
			suffix = "-windows-32"
		case Amd64:
			suffix = "-windows-64"
		}
	case MacOS:
		suffix = "-macos"
	case Linux:
		switch arch {
		case X86:
			suffix = "-linux-32"
		case Amd64:
			suffix = "-linux-64"
		case Arm:
			suffix = "-linux-arm"
		case Arm64:
			suffix = "-linux-arm64"
		}

	}
	return suffix
}
