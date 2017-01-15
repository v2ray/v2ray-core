package main

import (
	"strings"
)

type GoOS string

const (
	Windows   = GoOS("windows")
	MacOS     = GoOS("darwin")
	Linux     = GoOS("linux")
	FreeBSD   = GoOS("freebsd")
	OpenBSD   = GoOS("openbsd")
	UnknownOS = GoOS("unknown")
)

type GoArch string

const (
	X86         = GoArch("386")
	Amd64       = GoArch("amd64")
	Arm         = GoArch("arm")
	Arm64       = GoArch("arm64")
	Mips64      = GoArch("mips64")
	Mips        = GoArch("mips")
	MipsLE      = GoArch("mipsle")
	UnknownArch = GoArch("unknown")
)

func parseOS(rawOS string) GoOS {
	osStr := strings.ToLower(rawOS)
	switch osStr {
	case "windows", "win":
		return Windows
	case "darwin", "mac", "macos", "osx":
		return MacOS
	case "linux", "debian", "ubuntu", "redhat", "centos":
		return Linux
	case "freebsd":
		return FreeBSD
	case "openbsd":
		return OpenBSD
	default:
		return UnknownOS
	}
}

func parseArch(rawArch string) GoArch {
	archStr := strings.ToLower(rawArch)
	switch archStr {
	case "x86", "386", "i386":
		return X86
	case "amd64", "x86-64", "x64":
		return Amd64
	case "arm":
		return Arm
	case "arm64":
		return Arm64
	case "mips":
		return Mips
	case "mipsle":
		return MipsLE
	case "mips64":
		return Mips64
	default:
		return UnknownArch
	}
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
		case Mips64:
			suffix = "-linux-mips64"
		case Mips:
			suffix = "-linux-mips"
		case MipsLE:
			suffix = "-linux-mipsle"
		}
	case FreeBSD:
		switch arch {
		case X86:
			suffix = "-freebsd-32"
		case Amd64:
			suffix = "-freebsd-64"
		case Arm:
			suffix = "-freebsd-arm"
		}
	case OpenBSD:
		switch arch {
		case X86:
			suffix = "-openbsd-32"
		case Amd64:
			suffix = "-openbsd-64"
		}
	}

	return suffix
}
