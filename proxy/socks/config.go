package socks

import (
	"encoding/json"
	"errors"

	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
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

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

func (this *ServerConfig) UnmarshalJSON(data []byte) error {
	type SocksConfig struct {
		AuthMethod string            `json:"auth"`
		Accounts   []*Account        `json:"accounts"`
		UDP        bool              `json:"udp"`
		Host       *v2net.IPOrDomain `json:"ip"`
		Timeout    uint32            `json:"timeout"`
	}

	rawConfig := new(SocksConfig)
	if err := json.Unmarshal(data, rawConfig); err != nil {
		return errors.New("Socks: Failed to parse config: " + err.Error())
	}
	if rawConfig.AuthMethod == AuthMethodNoAuth {
		this.AuthType = AuthType_NO_AUTH
	} else if rawConfig.AuthMethod == AuthMethodUserPass {
		this.AuthType = AuthType_PASSWORD
	} else {
		log.Error("Socks: Unknown auth method: ", rawConfig.AuthMethod)
		return common.ErrBadConfiguration
	}

	if len(rawConfig.Accounts) > 0 {
		this.Accounts = make(map[string]string, len(rawConfig.Accounts))
		for _, account := range rawConfig.Accounts {
			this.Accounts[account.Username] = account.Password
		}
	}

	this.UdpEnabled = rawConfig.UDP
	if rawConfig.Host != nil {
		this.Address = rawConfig.Host
	}

	if rawConfig.Timeout >= 0 {
		this.Timeout = rawConfig.Timeout
	}
	return nil
}
