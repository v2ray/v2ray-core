// +build json

package http

import (
	"encoding/json"
	"errors"

	"github.com/v2ray/v2ray-core/proxy/registry"
)

// UnmarshalJSON implements json.Unmarshaler
func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Timeout int `json:"timeout"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("HTTP: Failed to parse config: " + err.Error())
	}
	this.Timeout = jsonConfig.Timeout

	return nil
}

func init() {
	registry.RegisterInboundConfig("http", func() interface{} { return new(Config) })
}
