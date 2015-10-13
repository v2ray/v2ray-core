// +build !windows

package platform

import (
	"os"
)

func ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

func LineSeparator() string {
	return "\n"
}
