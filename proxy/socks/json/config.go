package json

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/v2ray/v2ray-core/proxy/internal/config"
	jsonconfig "github.com/v2ray/v2ray-core/proxy/internal/config/json"
)

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

type SocksAccount struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

type SocksAccountMap map[string]string

func (this *SocksAccountMap) UnmarshalJSON(data []byte) error {
	var accounts []SocksAccount
	err := json.Unmarshal(data, &accounts)
	if err != nil {
		return err
	}
	*this = make(map[string]string)
	for _, account := range accounts {
		(*this)[account.Username] = account.Password
	}
	return nil
}

func (this *SocksAccountMap) HasAccount(user, pass string) bool {
	if actualPass, found := (*this)[user]; found {
		return actualPass == pass
	}
	return false
}

type IPAddress net.IP

func (this *IPAddress) UnmarshalJSON(data []byte) error {
	var ipStr string
	err := json.Unmarshal(data, &ipStr)
	if err != nil {
		return err
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return errors.New("Unknown IP format: " + ipStr)
	}
	*this = IPAddress(ip)
	return nil
}

type SocksConfig struct {
	AuthMethod string          `json:"auth"`
	Accounts   SocksAccountMap `json:"accounts"`
	UDP        bool            `json:"udp"`
	HostIP     IPAddress       `json:"ip"`
}

func (sc *SocksConfig) IsNoAuth() bool {
	return sc.AuthMethod == AuthMethodNoAuth
}

func (sc *SocksConfig) IsPassword() bool {
	return sc.AuthMethod == AuthMethodUserPass
}

func (sc *SocksConfig) HasAccount(user, pass string) bool {
	return sc.Accounts.HasAccount(user, pass)
}

func (sc *SocksConfig) IP() net.IP {
	return net.IP(sc.HostIP)
}

func (this *SocksConfig) UDPEnabled() bool {
	return this.UDP
}

func init() {
	config.RegisterInboundConnectionConfig("socks", jsonconfig.JsonConfigLoader(func() interface{} {
		return &SocksConfig{
			HostIP: IPAddress(net.IPv4(127, 0, 0, 1)),
		}
	}))
}
