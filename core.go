// Package core provides common definitions and functionalities of V2Ray.
package core

import (
	"fmt"
)

var (
	version  = "0.8"
	build    = "Custom"
	codename = "Post Apocalypse"
	intro    = "A stable and unbreakable connection for everyone."
)

func PrintVersion() {
	fmt.Printf("V2Ray %s (%s) %s", version, codename, build)
	fmt.Println()
	fmt.Println(intro)
}
