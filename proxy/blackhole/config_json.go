// +build json

package blackhole

import (
	"encoding/json"
	"errors"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common/loader"
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
		response, _, err := configLoader.Load(jsonConfig.Response)
		if err != nil {
			return errors.New("Blackhole: Failed to parse response config: " + err.Error())
		}
		this.Response = loader.NewTypedSettings(response.(proto.Message))
	}

	return nil
}

var (
	configLoader = loader.NewJSONConfigLoader(
		loader.NamedTypeMap{
			"none": loader.GetType(new(NoneResponse)),
			"http": loader.GetType(new(HTTPResponse)),
		},
		"type",
		"")
)
