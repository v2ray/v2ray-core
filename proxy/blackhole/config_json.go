// +build json

package blackhole

import (
	"encoding/json"
	"errors"

	"strings"
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

	if jsonConfig.Response != nil {
		response, rType, err := configLoader.Load(jsonConfig.Response)
		if err != nil {
			return errors.New("Blackhole: Failed to parse response config: " + err.Error())
		}
		this.Response = new(Response)
		switch rType {
		case strings.ToLower(Response_Type_name[int32(Response_None)]):
			this.Response.Type = Response_None
		case strings.ToLower(Response_Type_name[int32(Response_HTTP)]):
			this.Response.Type = Response_HTTP
		}
		this.Response.Settings = response.(ResponseConfig).AsAny()
	}

	return nil
}

var (
	configLoader = loader.NewJSONConfigLoader(cache, "type", "")
)

func init() {
	registry.RegisterOutboundConfig("blackhole", func() interface{} { return new(Config) })
}
