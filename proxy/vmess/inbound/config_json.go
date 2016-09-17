// +build json

package inbound

import (
	"encoding/json"
	"errors"

	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/proxy/vmess"

	"github.com/golang/protobuf/ptypes"
)

func (this *DetourConfig) UnmarshalJSON(data []byte) error {
	type JsonDetourConfig struct {
		ToTag string `json:"to"`
	}
	jsonConfig := new(JsonDetourConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("VMessIn: Failed to parse detour config: " + err.Error())
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
		return errors.New("VMessIn: Failed to parse features config: " + err.Error())
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
		return errors.New("VMessIn: Failed to parse default config: " + err.Error())
	}
	this.AlterIDs = jsonConfig.AlterIDs
	if this.AlterIDs == 0 {
		this.AlterIDs = 32
	}
	this.Level = uint32(jsonConfig.Level)
	return nil
}

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Users        []json.RawMessage `json:"clients"`
		Features     *FeaturesConfig   `json:"features"`
		Defaults     *DefaultConfig    `json:"default"`
		DetourConfig *DetourConfig     `json:"detour"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("VMessIn: Failed to parse config: " + err.Error())
	}
	this.Features = jsonConfig.Features // Backward compatibility
	this.Defaults = jsonConfig.Defaults
	if this.Defaults == nil {
		this.Defaults = &DefaultConfig{
			Level:    0,
			AlterIDs: 32,
		}
	}
	this.DetourConfig = jsonConfig.DetourConfig
	// Backward compatibility
	if this.Features != nil && this.DetourConfig == nil {
		this.DetourConfig = this.Features.Detour
	}
	this.AllowedUsers = make([]*protocol.User, len(jsonConfig.Users))
	for idx, rawData := range jsonConfig.Users {
		user := new(protocol.User)
		if err := json.Unmarshal(rawData, user); err != nil {
			return errors.New("VMess|Inbound: Invalid user: " + err.Error())
		}
		account := new(vmess.AccountPB)
		if err := json.Unmarshal(rawData, account); err != nil {
			return errors.New("VMess|Inbound: Invalid user: " + err.Error())
		}
		anyAccount, err := ptypes.MarshalAny(account)
		if err != nil {
			log.Error("VMess|Inbound: Failed to create account: ", err)
			return common.ErrBadConfiguration
		}
		user.Account = anyAccount
		this.AllowedUsers[idx] = user
	}

	return nil
}

func init() {
	registry.RegisterInboundConfig("vmess", func() interface{} { return new(Config) })
}
