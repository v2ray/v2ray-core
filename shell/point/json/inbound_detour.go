package json

import (
	"encoding/json"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netjson "github.com/v2ray/v2ray-core/common/net/json"
	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
)

type InboundDetourConfig struct {
	ProtocolValue  string               `json:"protocol"`
	PortRangeValue *v2netjson.PortRange `json:"port"`
	SettingsValue  json.RawMessage      `json:"settings"`
	TagValue       string               `json:"tag"`
}

func (this *InboundDetourConfig) Protocol() string {
	return this.ProtocolValue
}

func (this *InboundDetourConfig) PortRange() v2net.PortRange {
	return this.PortRangeValue
}

func (this *InboundDetourConfig) Settings() interface{} {
	return loadConnectionConfig(this.SettingsValue, this.ProtocolValue, proxyconfig.TypeInbound)
}

func (this *InboundDetourConfig) Tag() string {
	return this.TagValue
}
