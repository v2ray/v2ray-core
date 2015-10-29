package json

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/config"
)

type ConnectionConfig struct {
	ProtocolString  string          `json:"protocol"`
	SettingsMessage json.RawMessage `json:"settings"`
	Type            config.Type     `json:"-"`
}

func (c *ConnectionConfig) Protocol() string {
	return c.ProtocolString
}

func (c *ConnectionConfig) Settings() interface{} {
	creator, found := configCache[getConfigKey(c.Protocol(), c.Type)]
	if !found {
		panic("Unknown protocol " + c.Protocol())
	}
	configObj := creator()
	err := json.Unmarshal(c.SettingsMessage, configObj)
	if err != nil {
		log.Error("Unable to parse connection config: %v", err)
		panic("Failed to parse connection config.")
	}
	return configObj
}
