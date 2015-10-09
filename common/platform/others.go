// +build !windows

package platform

import (
	"os"
)

type otherPlatformEnvironment struct {
}

var environmentInstance = &otherPlatformEnvironment{}

func (e *otherPlatformEnvironment) ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}
