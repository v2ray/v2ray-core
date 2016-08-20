package socks

import (
	"v2ray.com/core/common/protocol"
)

type Account struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (this *Account) Equals(another protocol.Account) bool {
	socksAccount, ok := another.(*Account)
	if !ok {
		return false
	}
	return this.Username == socksAccount.Username
}
