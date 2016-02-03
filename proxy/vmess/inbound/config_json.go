// +build json

package inbound

import (
	"encoding/json"

	proto "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

func (this *DetourConfig) UnmarshalJSON(data []byte) error {
	type JsonDetourConfig struct {
		ToTag string `json:"to"`
	}
	jsonConfig := new(JsonDetourConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.ToTag = jsonConfig.ToTag
	return nil
}

func (this *FeaturesConfig) UnmarshalJSON(data []byte) error {
	type JsonFeaturesConfig struct {
		Detour *DetourConfig `json:"detour"`
	}
	jsonConfig := new(JsonFeaturesConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Detour = jsonConfig.Detour
	return nil
}

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Users    []*proto.User   `json:"clients"`
		Features *FeaturesConfig `json:"features"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.AllowedUsers = jsonConfig.Users
	this.Features = jsonConfig.Features
	return nil
}

func init() {
	config.RegisterInboundConfig("vmess",
		func(data []byte) (interface{}, error) {
			config := new(Config)
			err := json.Unmarshal(data, config)
			return config, err
		})
}
