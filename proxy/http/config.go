package http

import (
	"v2ray.com/core/common/protocol"
)

func (a *Account) Equals(another protocol.Account) bool {
	if account, ok := another.(*Account); ok {
		return a.Username == account.Username
	}
	return false
}

func (a *Account) AsAccount() (protocol.Account, error) {
	return a, nil
}

func (sc *ServerConfig) HasAccount(username, password string) bool {
	if sc.Accounts == nil {
		return false
	}

	p, found := sc.Accounts[username]
	if !found {
		return false
	}
	return p == password
}
