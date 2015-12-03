package config

import (
	"net"
)

type SocksConfig interface {
	IsNoAuth() bool
	IsPassword() bool
	HasAccount(username, password string) bool
	IP() net.IP
	UDPEnabled() bool
}
