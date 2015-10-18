// Package core provides common definitions and functionalities of V2Ray.
package core

import (
	"fmt"

	"github.com/v2ray/v2ray-core/common/platform"
)

var (
	version  = "0.9"
	build    = "Custom"
	codename = "Post Apocalypse"
	intro    = "A stable and unbreakable connection for everyone."
)

func PrintVersion() {
	fmt.Printf("V2Ray %s (%s) %s%s", version, codename, build, platform.LineSeparator())
	fmt.Printf("%s%s", intro, platform.LineSeparator())
}
