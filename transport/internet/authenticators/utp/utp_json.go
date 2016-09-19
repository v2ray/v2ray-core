// +build json

package utp

import (
	"v2ray.com/core/transport/internet"
)

func init() {
	internet.RegisterAuthenticatorConfig("utp", func() interface{} { return &Config{} })
}
