package conf

import (
	"encoding/json"
	"strings"

	"v2ray.com/core/common/errors"
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
	Security string `json:"security"`
}

func (v *VMessAccount) Build() *vmess.Account {
	var st protocol.SecurityType
	switch strings.ToLower(v.Security) {
	case "aes-128-gcm":
		st = protocol.SecurityType_AES128_GCM
	case "chacha20-poly1305":
		st = protocol.SecurityType_CHACHA20_POLY1305
	case "auto":
		st = protocol.SecurityType_AUTO
	case "none":
		st = protocol.SecurityType_NONE
	default:
		st = protocol.SecurityType_LEGACY
	}
	return &vmess.Account{
		Id:      v.ID,
		AlterId: uint32(v.AlterIds),
		SecuritySettings: &protocol.SecurityConfig{
			Type: st,
		},
	}
}

type VMessDetourConfig struct {
	ToTag string `json:"to"`
}

func (v *VMessDetourConfig) Build() *inbound.DetourConfig {
	return &inbound.DetourConfig{
		To: v.ToTag,
	}
}

type FeaturesConfig struct {
	Detour *VMessDetourConfig `json:"detour"`
}

type VMessDefaultConfig struct {
	AlterIDs uint16 `json:"alterId"`
	Level    byte   `json:"level"`
}

func (v *VMessDefaultConfig) Build() *inbound.DefaultConfig {
	config := new(inbound.DefaultConfig)
	config.AlterId = uint32(v.AlterIDs)
	if config.AlterId == 0 {
		config.AlterId = 32
	}
	config.Level = uint32(v.Level)
	return config
}

type VMessInboundConfig struct {
	Users        []json.RawMessage   `json:"clients"`
	Features     *FeaturesConfig     `json:"features"`
	Defaults     *VMessDefaultConfig `json:"default"`
	DetourConfig *VMessDetourConfig  `json:"detour"`
}

func (v *VMessInboundConfig) Build() (*loader.TypedSettings, error) {
	config := new(inbound.Config)

	if v.Defaults != nil {
		config.Default = v.Defaults.Build()
	}

	if v.DetourConfig != nil {
		config.Detour = v.DetourConfig.Build()
	} else if v.Features != nil && v.Features.Detour != nil {
		config.Detour = v.Features.Detour.Build()
	}

	config.User = make([]*protocol.User, len(v.Users))
	for idx, rawData := range v.Users {
		user := new(protocol.User)
		if err := json.Unmarshal(rawData, user); err != nil {
			return nil, errors.Base(err).Message("Invalid VMess user.")
		}
		account := new(VMessAccount)
		if err := json.Unmarshal(rawData, account); err != nil {
			return nil, errors.Base(err).Message("Invalid VMess user.")
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

func (v *VMessOutboundConfig) Build() (*loader.TypedSettings, error) {
	config := new(outbound.Config)

	if len(v.Receivers) == 0 {
		return nil, errors.New("0 VMess receiver configured.")
	}
	serverSpecs := make([]*protocol.ServerEndpoint, len(v.Receivers))
	for idx, rec := range v.Receivers {
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
				return nil, errors.Base(err).Message("Invalid VMess user.")
			}
			account := new(VMessAccount)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return nil, errors.Base(err).Message("Invalid VMess user.")
			}
			user.Account = loader.NewTypedSettings(account.Build())
			spec.User = append(spec.User, user)
		}
		serverSpecs[idx] = spec
	}
	config.Receiver = serverSpecs
	return loader.NewTypedSettings(config), nil
}
