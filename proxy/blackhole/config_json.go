// +build json

package blackhole

import (
	"encoding/json"
	"errors"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/proxy/registry"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Response json.RawMessage `json:"response"`
	}
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Blackhole: Failed to parse config: " + err.Error())
	}

	this.Response = new(NoneResponse)
	if jsonConfig.Response != nil {
		loader := loader.NewJSONConfigLoader(loader.ConfigCreatorCache{}, "type", "")
		loader.RegisterCreator("none", func() interface{} { return new(NoneResponse) })
		loader.RegisterCreator("http", func() interface{} { return new(HTTPResponse) })
		response, _, err := loader.Load(jsonConfig.Response)
		if err != nil {
			return errors.New("Blackhole: Failed to parse response config: " + err.Error())
		}
		this.Response = response.(Response)
	}

	return nil
}

func init() {
	registry.RegisterOutboundConfig("blackhole", func() interface{} { return new(Config) })
}
