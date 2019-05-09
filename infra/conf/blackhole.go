package conf

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/blackhole"
)

type NoneResponse struct{}

func (*NoneResponse) Build() (proto.Message, error) {
	return new(blackhole.NoneResponse), nil
}

type HttpResponse struct{}

func (*HttpResponse) Build() (proto.Message, error) {
	return new(blackhole.HTTPResponse), nil
}

type BlackholeConfig struct {
	Response json.RawMessage `json:"response"`
}

func (v *BlackholeConfig) Build() (proto.Message, error) {
	config := new(blackhole.Config)
	if v.Response != nil {
		response, _, err := configLoader.Load(v.Response)
		if err != nil {
			return nil, newError("Config: Failed to parse Blackhole response config.").Base(err)
		}
		responseSettings, err := response.(Buildable).Build()
		if err != nil {
			return nil, err
		}
		config.Response = serial.ToTypedMessage(responseSettings)
	}

	return config, nil
}

var (
	configLoader = NewJSONConfigLoader(
		ConfigCreatorCache{
			"none": func() interface{} { return new(NoneResponse) },
			"http": func() interface{} { return new(HttpResponse) },
		},
		"type",
		"")
)
