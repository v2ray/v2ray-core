package json

import (
	"encoding/json"
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

func (this *OutboundDetourConfig) Settings() []byte {
	return []byte(this.SettingsValue)
}
