package noop

import (
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/transport/internet"
)

type NoOpAuthenticator struct{}

func (this NoOpAuthenticator) Overhead() int {
	return 0
}
func (this NoOpAuthenticator) Open(payload *alloc.Buffer) bool {
	return true
}
func (this NoOpAuthenticator) Seal(payload *alloc.Buffer) {}

type NoOpAuthenticatorFactory struct{}

func (this NoOpAuthenticatorFactory) Create(config internet.AuthenticatorConfig) internet.Authenticator {
	return NoOpAuthenticator{}
}

type NoOpAuthenticatorConfig struct{}

func init() {
	internet.RegisterAuthenticator("none", NoOpAuthenticatorFactory{}, func() interface{} { return &NoOpAuthenticatorConfig{} })
}
