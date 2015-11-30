package json

import (
	"encoding/json"

	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
)

type OutboundDetourConfig struct {
	ProtocolValue string          `json:"protocol"`
	TagValue      string          `json:"tag"`
	SettingsValue json.RawMessage `json:"settings"`
}

func (this *OutboundDetourConfig) Protocol() string {
	return this.ProtocolValue
}

func (this *OutboundDetourConfig) Tag() string {
	return this.TagValue
}

func (this *OutboundDetourConfig) Settings() interface{} {
	return loadConnectionConfig(this.SettingsValue, this.ProtocolValue, proxyconfig.TypeOutbound)
}
