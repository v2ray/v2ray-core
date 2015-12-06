package socks

import (
	"net"
)

type Config interface {
	IsNoAuth() bool
	IsPassword() bool
	HasAccount(username, password string) bool
	IP() net.IP
	UDPEnabled() bool
}
