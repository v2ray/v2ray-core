package json

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
	proxyjson "github.com/v2ray/v2ray-core/proxy/common/config/json"
)

type ConnectionConfig struct {
	ProtocolString  string           `json:"protocol"`
	SettingsMessage json.RawMessage  `json:"settings"`
	Type            proxyconfig.Type `json:"-"`
}

func (c *ConnectionConfig) Protocol() string {
	return c.ProtocolString
}

func (c *ConnectionConfig) Settings() interface{} {
	configObj := proxyjson.CreateConfig(c.Protocol(), c.Type)
	if configObj == nil {
		panic("Unknown protocol " + c.Protocol())
	}
	err := json.Unmarshal(c.SettingsMessage, configObj)
	if err != nil {
		log.Error("Unable to parse connection config: %v", err)
		panic("Failed to parse connection config.")
	}
	return configObj
}
