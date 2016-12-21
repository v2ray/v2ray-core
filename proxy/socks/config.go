package socks

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

func (v *Account) Equals(another protocol.Account) bool {
	if account, ok := another.(*Account); ok {
		return v.Username == account.Username
	}
	return false
}

func (v *Account) AsAccount() (protocol.Account, error) {
	return v, nil
}

func NewAccount() protocol.AsAccount {
	return &Account{}
}

func (v *ServerConfig) HasAccount(username, password string) bool {
	if v.Accounts == nil {
		return false
	}
	storedPassed, found := v.Accounts[username]
	if !found {
		return false
	}
	return storedPassed == password
}

func (v *ServerConfig) GetNetAddress() v2net.Address {
	if v.Address == nil {
		return v2net.LocalHostIP
	}
	return v.Address.AsAddress()
}
