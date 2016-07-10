// +build json

package http

import (
	"encoding/json"
	"errors"

	"github.com/v2ray/v2ray-core/proxy/internal"
)

// UnmarshalJSON implements json.Unmarshaler
func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("HTTP: Failed to parse config: " + err.Error())
	}

	return nil
}

func init() {
	internal.RegisterInboundConfig("http", func() interface{} { return new(Config) })
}
