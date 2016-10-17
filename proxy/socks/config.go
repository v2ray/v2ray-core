package socks

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"

	"github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
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

func (this *Account) AsAny() (*google_protobuf.Any, error) {
	return ptypes.MarshalAny(this)
}

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
