// +build !confonly

package socks

import "v2ray.com/core/common/protocol"

func (a *Account) Equals(another protocol.Account) bool {
	if account, ok := another.(*Account); ok {
		return a.Username == account.Username
	}
	return false
}

func (a *Account) AsAccount() (protocol.Account, error) {
	return a, nil
}

func (c *ServerConfig) HasAccount(username, password string) bool {
	if c.Accounts == nil {
		return false
	}
	storedPassed, found := c.Accounts[username]
	if !found {
		return false
	}
	return storedPassed == password
}
