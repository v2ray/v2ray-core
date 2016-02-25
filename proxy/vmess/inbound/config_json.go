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

func (this *DefaultConfig) UnmarshalJSON(data []byte) error {
	type JsonDefaultConfig struct {
		AlterIDs uint16 `json:"alterId"`
		Level    byte   `json:"level"`
	}
	jsonConfig := new(JsonDefaultConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.AlterIDs = jsonConfig.AlterIDs
	if this.AlterIDs == 0 {
		this.AlterIDs = 32
	}
	this.Level = proto.UserLevel(jsonConfig.Level)
	return nil
}

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Users    []*proto.User   `json:"clients"`
		Features *FeaturesConfig `json:"features"`
		Defaults *DefaultConfig  `json:"default"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.AllowedUsers = jsonConfig.Users
	this.Features = jsonConfig.Features
	this.Defaults = jsonConfig.Defaults
	if this.Defaults == nil {
		this.Defaults = &DefaultConfig{
			Level:    proto.UserLevel(0),
			AlterIDs: 32,
		}
	}
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
