package noop

import (
	"net"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/transport/internet"
)

type NoOpAuthenticator struct{}

func (v NoOpAuthenticator) Overhead() int {
	return 0
}
func (v NoOpAuthenticator) Open(payload *alloc.Buffer) bool {
	return true
}
func (v NoOpAuthenticator) Seal(payload *alloc.Buffer) {}

type NoOpAuthenticatorFactory struct{}

func (v NoOpAuthenticatorFactory) Create(config interface{}) internet.Authenticator {
	return NoOpAuthenticator{}
}

type NoOpConnectionAuthenticator struct{}

func (NoOpConnectionAuthenticator) Client(conn net.Conn) net.Conn {
	return conn
}

func (NoOpConnectionAuthenticator) Server(conn net.Conn) net.Conn {
	return conn
}

type NoOpConnectionAuthenticatorFactory struct{}

func (NoOpConnectionAuthenticatorFactory) Create(config interface{}) internet.ConnectionAuthenticator {
	return NoOpConnectionAuthenticator{}
}

func init() {
	internet.RegisterAuthenticator(loader.GetType(new(Config)), NoOpAuthenticatorFactory{})
	internet.RegisterConnectionAuthenticator(loader.GetType(new(Config)), NoOpConnectionAuthenticatorFactory{})
}
