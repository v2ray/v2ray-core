package json

import (
	"net"

	"github.com/v2ray/v2ray-core/proxy/common/config/json"
)

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

type SocksAccount struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

type SocksConfig struct {
	AuthMethod string         `json:"auth"`
	Accounts   []SocksAccount `json:"accounts"`
	UDPEnabled bool           `json:"udp"`
	HostIP     string         `json:"ip"`

	accountMap map[string]string
	ip         net.IP
}

func (sc *SocksConfig) Initialize() {
	sc.accountMap = make(map[string]string)
	for _, account := range sc.Accounts {
		sc.accountMap[account.Username] = account.Password
	}

	if len(sc.HostIP) > 0 {
		sc.ip = net.ParseIP(sc.HostIP)
		if sc.ip == nil {
			sc.ip = net.IPv4(127, 0, 0, 1)
		}
	}
}

func (sc *SocksConfig) IsNoAuth() bool {
	return sc.AuthMethod == AuthMethodNoAuth
}

func (sc *SocksConfig) IsPassword() bool {
	return sc.AuthMethod == AuthMethodUserPass
}

func (sc *SocksConfig) HasAccount(user, pass string) bool {
	if actualPass, found := sc.accountMap[user]; found {
		return actualPass == pass
	}
	return false
}

func (sc *SocksConfig) IP() net.IP {
	return sc.ip
}

func init() {
	json.RegisterInboundConnectionConfig("socks", func() interface{} {
		return new(SocksConfig)
	})
}
