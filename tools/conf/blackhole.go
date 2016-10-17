package conf

import (
	"encoding/json"
	"errors"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/proxy/blackhole"
)

type NoneResponse struct{}

func (*NoneResponse) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(blackhole.NoneResponse)), nil
}

type HttpResponse struct{}

func (*HttpResponse) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(blackhole.HTTPResponse)), nil
}

type BlackholeConfig struct {
	Response json.RawMessage `json:"response"`
}

func (this *BlackholeConfig) Build() (*loader.TypedSettings, error) {
	config := new(blackhole.Config)
	if this.Response != nil {
		response, _, err := configLoader.Load(this.Response)
		if err != nil {
			return nil, errors.New("Blackhole: Failed to parse response config: " + err.Error())
		}
		responseSettings, err := response.(Buildable).Build()
		if err != nil {
			return nil, err
		}
		config.Response = responseSettings
	}

	return loader.NewTypedSettings(config), nil
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
