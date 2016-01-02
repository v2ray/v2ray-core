package json

import (
	"encoding/json"
)

type ConnectionConfig struct {
	ProtocolString  string          `json:"protocol"`
	SettingsMessage json.RawMessage `json:"settings"`
}

func (c *ConnectionConfig) Protocol() string {
	return c.ProtocolString
}

func (c *ConnectionConfig) Settings() []byte {
	return []byte(c.SettingsMessage)
}
