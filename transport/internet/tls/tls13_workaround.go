// +build !confonly

package tls

import (
	"os"
	"strings"
)

func init() {
	// opt-in TLS 1.3 for Go1.12
	// TODO: remove this line when Go1.13 is released.
	if !strings.Contains(os.Getenv("GODEBUG"), "tls13") {
		_ = os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
	}
}
