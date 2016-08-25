// +build json

package noop

import (
	"v2ray.com/core/transport/internet"
)

func init() {
	internet.RegisterAuthenticatorConfig("none", func() interface{} { return &NoOpAuthenticatorConfig{} })
}
