// +build json

package srtp

import (
	"v2ray.com/core/transport/internet"
)

func init() {
	internet.RegisterAuthenticatorConfig("srtp", func() interface{} { return &Config{} })
}
