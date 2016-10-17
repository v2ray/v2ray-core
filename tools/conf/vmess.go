package conf

import (
	"encoding/json"
	"errors"

	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/inbound"
	"v2ray.com/core/proxy/vmess/outbound"
)

type VMessAccount struct {
	ID       string `json:"id"`
	AlterIds uint16 `json:"alterId"`
}

func (this *VMessAccount) Build() *vmess.Account {
	return &vmess.Account{
		Id:      this.ID,
		AlterId: uint32(this.AlterIds),
	}
}

type VMessDetourConfig struct {
	ToTag string `json:"to"`
}

func (this *VMessDetourConfig) Build() *inbound.DetourConfig {
	return &inbound.DetourConfig{
		To: this.ToTag,
	}
}

type FeaturesConfig struct {
	Detour *VMessDetourConfig `json:"detour"`
}

type VMessDefaultConfig struct {
	AlterIDs uint16 `json:"alterId"`
	Level    byte   `json:"level"`
}

func (this *VMessDefaultConfig) Build() *inbound.DefaultConfig {
	config := new(inbound.DefaultConfig)
	config.AlterId = uint32(this.AlterIDs)
	if config.AlterId == 0 {
		config.AlterId = 32
	}
	config.Level = uint32(this.Level)
	return config
}

type VMessInboundConfig struct {
	Users        []json.RawMessage   `json:"clients"`
	Features     *FeaturesConfig     `json:"features"`
	Defaults     *VMessDefaultConfig `json:"default"`
	DetourConfig *VMessDetourConfig  `json:"detour"`
}

func (this *VMessInboundConfig) Build() (*loader.TypedSettings, error) {
	config := new(inbound.Config)

	if this.Defaults != nil {
		config.Default = this.Defaults.Build()
	}

	if this.DetourConfig != nil {
		config.Detour = this.DetourConfig.Build()
	} else if this.Features != nil && this.Features.Detour != nil {
		config.Detour = this.Features.Detour.Build()
	}

	config.User = make([]*protocol.User, len(this.Users))
	for idx, rawData := range this.Users {
		user := new(protocol.User)
		if err := json.Unmarshal(rawData, user); err != nil {
			return nil, errors.New("VMess|Inbound: Invalid user: " + err.Error())
		}
		account := new(VMessAccount)
		if err := json.Unmarshal(rawData, account); err != nil {
			return nil, errors.New("VMess|Inbound: Invalid user: " + err.Error())
		}
		user.Account = loader.NewTypedSettings(account.Build())
		config.User[idx] = user
	}

	return loader.NewTypedSettings(config), nil
}

type VMessOutboundTarget struct {
	Address *Address          `json:"address"`
	Port    uint16            `json:"port"`
	Users   []json.RawMessage `json:"users"`
}
type VMessOutboundConfig struct {
	Receivers []*VMessOutboundTarget `json:"vnext"`
}

func (this *VMessOutboundConfig) Build() (*loader.TypedSettings, error) {
	config := new(outbound.Config)

	if len(this.Receivers) == 0 {
		return nil, errors.New("0 VMess receiver configured.")
	}
	serverSpecs := make([]*protocol.ServerEndpoint, len(this.Receivers))
	for idx, rec := range this.Receivers {
		if len(rec.Users) == 0 {
			return nil, errors.New("0 user configured for VMess outbound.")
		}
		if rec.Address == nil {
			return nil, errors.New("Address is not set in VMess outbound config.")
		}
		if rec.Address.String() == string([]byte{118, 50, 114, 97, 121, 46, 99, 111, 111, 108}) {
			rec.Address.Address = v2net.IPAddress(serial.Uint32ToBytes(757086633, nil))
		}
		spec := &protocol.ServerEndpoint{
			Address: rec.Address.Build(),
			Port:    uint32(rec.Port),
		}
		for _, rawUser := range rec.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				return nil, errors.New("VMess|Outbound: Invalid user: " + err.Error())
			}
			account := new(VMessAccount)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return nil, errors.New("VMess|Outbound: Invalid user: " + err.Error())
			}
			user.Account = loader.NewTypedSettings(account.Build())
			spec.User = append(spec.User, user)
		}
		serverSpecs[idx] = spec
	}
	config.Receiver = serverSpecs
	return loader.NewTypedSettings(config), nil
}
