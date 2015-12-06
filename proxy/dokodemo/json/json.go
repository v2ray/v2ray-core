package json

import (
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netjson "github.com/v2ray/v2ray-core/common/net/json"
	"github.com/v2ray/v2ray-core/proxy/common/config/json"
)

type DokodemoConfig struct {
	Host         string                 `json:"address"`
	Port         v2net.Port             `json:"port"`
	NetworkList  *v2netjson.NetworkList `json:"network"`
	TimeoutValue int                    `json:"timeout"`
}

func (this *DokodemoConfig) Address() v2net.Address {
	ip := net.ParseIP(this.Host)
	if ip != nil {
		return v2net.IPAddress(ip, this.Port)
	} else {
		return v2net.DomainAddress(this.Host, this.Port)
	}
}

func (this *DokodemoConfig) Network() v2net.NetworkList {
	return this.NetworkList
}

func (this *DokodemoConfig) Timeout() int {
	return this.TimeoutValue
}

func init() {
	json.RegisterInboundConnectionConfig("dokodemo-door", func() interface{} {
		return new(DokodemoConfig)
	})
}
