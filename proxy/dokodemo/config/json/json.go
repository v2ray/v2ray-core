package json

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netjson "github.com/v2ray/v2ray-core/common/net/json"
	"github.com/v2ray/v2ray-core/proxy/common/config/json"
)

type DokodemoConfig struct {
	Host    string                 `json:"address"`
	Port    int                    `json:"port"`
	Network *v2netjson.NetworkList `json:"network"`
	Timeout int                    `json:"timeout"`

	address v2net.Address
}

func init() {
	json.RegisterInboundConnectionConfig("dokodemo-door", func() interface{} {
		return new(DokodemoConfig)
	})
}
