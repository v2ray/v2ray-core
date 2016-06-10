// +build json

package blackhole

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/loader"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Response json.RawMessage `json:"response"`
	}
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}

	this.Response = new(NoneResponse)
	if jsonConfig.Response == nil {
		loader := loader.NewJSONConfigLoader("type", "")
		loader.RegisterCreator("none", func() interface{} { return new(NoneResponse) })
		loader.RegisterCreator("http", func() interface{} { return new(HTTPResponse) })
		response, err := loader.Load(jsonConfig.Response)
		if err != nil {
			return err
		}
		this.Response = response.(Response)
	}

	return nil
}

func init() {
	internal.RegisterOutboundConfig("blackhole", func() interface{} { return new(Config) })
}
