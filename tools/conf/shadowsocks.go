package conf

import (
	"errors"
	"strings"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/shadowsocks"
)

type ShadowsocksServerConfig struct {
	Cipher   string `json:"method"`
	Password string `json:"password"`
	UDP      bool   `json:"udp"`
	Level    byte   `json:"level"`
	Email    string `json:"email"`
	OTA      *bool  `json:"ota"`
}

func (this *ShadowsocksServerConfig) Build() (*loader.TypedSettings, error) {
	config := new(shadowsocks.ServerConfig)
	config.UdpEnabled = this.UDP

	if len(this.Password) == 0 {
		return nil, errors.New("Shadowsocks password is not specified.")
	}
	account := &shadowsocks.Account{
		Password: this.Password,
		Ota:      shadowsocks.Account_Auto,
	}
	if this.OTA != nil {
		if *this.OTA {
			account.Ota = shadowsocks.Account_Enabled
		} else {
			account.Ota = shadowsocks.Account_Disabled
		}
	}
	cipher := strings.ToLower(this.Cipher)
	switch cipher {
	case "aes-256-cfb":
		account.CipherType = shadowsocks.CipherType_AES_256_CFB
	case "aes-128-cfb":
		account.CipherType = shadowsocks.CipherType_AES_128_CFB
	case "chacha20":
		account.CipherType = shadowsocks.CipherType_CHACHA20
	case "chacha20-ietf":
		account.CipherType = shadowsocks.CipherType_CHACHA20_IEFT
	default:
		return nil, errors.New("Unknown cipher method: " + cipher)
	}

	config.User = &protocol.User{
		Email:   this.Email,
		Level:   uint32(this.Level),
		Account: loader.NewTypedSettings(account),
	}

	return loader.NewTypedSettings(config), nil
}

type ShadowsocksServerTarget struct {
	Address  *Address `json:"address"`
	Port     uint16   `json:"port"`
	Cipher   string   `json:"method"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Ota      bool     `json:"ota"`
}

type ShadowsocksClientConfig struct {
	Servers []*ShadowsocksServerTarget `json:"servers"`
}

func (this *ShadowsocksClientConfig) Build() (*loader.TypedSettings, error) {
	config := new(shadowsocks.ClientConfig)

	if len(this.Servers) == 0 {
		return nil, errors.New("0 Shadowsocks server configured.")
	}

	serverSpecs := make([]*protocol.ServerEndpoint, len(this.Servers))
	for idx, server := range this.Servers {
		if server.Address == nil {
			return nil, errors.New("Shadowsocks server address is not set.")
		}
		if server.Port == 0 {
			return nil, errors.New("Invalid Shadowsocks port.")
		}
		if len(server.Password) == 0 {
			return nil, errors.New("Shadowsocks password is not specified.")
		}
		account := &shadowsocks.Account{
			Password: server.Password,
			Ota:      shadowsocks.Account_Enabled,
		}
		if !server.Ota {
			account.Ota = shadowsocks.Account_Disabled
		}
		cipher := strings.ToLower(server.Cipher)
		switch cipher {
		case "aes-256-cfb":
			account.CipherType = shadowsocks.CipherType_AES_256_CFB
		case "aes-128-cfb":
			account.CipherType = shadowsocks.CipherType_AES_128_CFB
		case "chacha20":
			account.CipherType = shadowsocks.CipherType_CHACHA20
		case "chacha20-ietf":
			account.CipherType = shadowsocks.CipherType_CHACHA20_IEFT
		default:
			return nil, errors.New("Unknown cipher method: " + cipher)
		}

		ss := &protocol.ServerEndpoint{
			Address: server.Address.Build(),
			Port:    uint32(server.Port),
			User: []*protocol.User{
				{
					Email:   server.Email,
					Account: loader.NewTypedSettings(account),
				},
			},
		}

		serverSpecs[idx] = ss
	}

	config.Server = serverSpecs

	return loader.NewTypedSettings(config), nil
}
