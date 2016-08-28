package socks

import (
	v2net "v2ray.com/core/common/net"
)

func (this *ServerConfig) HasAccount(username, password string) bool {
	if this.Accounts == nil {
		return false
	}
	storedPassed, found := this.Accounts[username]
	if !found {
		return false
	}
	return storedPassed == password
}

func (this *ServerConfig) GetNetAddress() v2net.Address {
	if this.Address == nil {
		return v2net.LocalHostIP
	}
	return this.Address.AsAddress()
}
