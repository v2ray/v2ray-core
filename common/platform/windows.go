// +build windows

package platform

import (
	"os"
)

type windowsEnvironment struct {
}

var environmentInstance = &windowsEnvironment{}

func (e *windowsEnvironment) ExpandEnv(s string) string {
	// TODO
	return s
}
