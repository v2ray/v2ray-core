package json

import (
	"encoding/json"

	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
	"github.com/v2ray/v2ray-core/shell/point/config"
)

type OutboundDetourConfig struct {
	ProtocolValue string          `json:"protocol"`
	TagValue      string          `json:"tag"`
	SettingsValue json.RawMessage `json:"settings"`
}

func (this *OutboundDetourConfig) Protocol() string {
	return this.ProtocolValue
}

func (this *OutboundDetourConfig) Tag() config.DetourTag {
	return config.DetourTag(this.TagValue)
}

func (this *OutboundDetourConfig) Settings() interface{} {
	return loadConnectionConfig(this.SettingsValue, this.ProtocolValue, proxyconfig.TypeOutbound)
}
