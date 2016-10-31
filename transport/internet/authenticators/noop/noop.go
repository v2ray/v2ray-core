package noop

import (
	"io"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
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

func (this NoOpAuthenticatorFactory) Create(config interface{}) internet.Authenticator {
	return NoOpAuthenticator{}
}

type NoOpConnectionAuthenticator struct{}

func (NoOpConnectionAuthenticator) Open(reader io.Reader) (bool, io.Reader) {
	return true, reader
}

func (NoOpConnectionAuthenticator) Seal(writer io.Writer) io.Writer {
	return writer
}

type NoOpConnectionAuthenticatorFactory struct{}

func (NoOpConnectionAuthenticatorFactory) Create(config interface{}) internet.ConnectionAuthenticator {
	return NoOpConnectionAuthenticator{}
}

func init() {
	internet.RegisterAuthenticator(loader.GetType(new(Config)), NoOpAuthenticatorFactory{})
	internet.RegisterConnectionAuthenticator(loader.GetType(new(Config)), NoOpConnectionAuthenticatorFactory{})
}
