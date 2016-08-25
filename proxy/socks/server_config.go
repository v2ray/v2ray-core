package socks

import (
	v2net "v2ray.com/core/common/net"
)

const (
	AuthTypeNoAuth   = byte(0)
	AuthTypePassword = byte(1)
)

type Config struct {
	AuthType   byte
	Accounts   map[string]string
	Address    v2net.Address
	UDPEnabled bool
	Timeout    uint32
}

func (this *Config) HasAccount(username, password string) bool {
	if this.Accounts == nil {
		return false
	}
	storedPassed, found := this.Accounts[username]
	if !found {
		return false
	}
	return storedPassed == password
}
