package json

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/app/point/config"
	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
)

type InboundDetourConfig struct {
	ProtocolValue  string
	PortRangeValue *PortRange
	SettingsValue  json.RawMessage
}

func (this *InboundDetourConfig) Protocol() string {
	return this.ProtocolValue
}

func (this *InboundDetourConfig) PortRange() config.PortRange {
	return this.PortRangeValue
}

func (this *InboundDetourConfig) Settings() interface{} {
	return loadConnectionConfig(this.SettingsValue, this.ProtocolValue, proxyconfig.TypeInbound)
}
