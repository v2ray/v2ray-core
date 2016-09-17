package socks

import (
	"v2ray.com/core/common/protocol"
)

func (this *Account) Equals(another protocol.Account) bool {
	if account, ok := another.(*Account); ok {
		return this.Username == account.Username
	}
	return false
}

func (this *Account) AsAccount() (protocol.Account, error) {
	return this, nil
}

func NewAccount() protocol.AsAccount {
	return &Account{}
}
