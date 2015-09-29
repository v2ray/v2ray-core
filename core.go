// Package core provides common definitions and functionalities of V2Ray.
package core

var (
	version = "0.6.2"
	build   = "Custom"
)

const (
	Codename = "Post Apocalypse"
	Intro    = "A stable and unbreakable connection for everyone."
)

func Version() string {
	return version
}

func Build() string {
	return build
}
