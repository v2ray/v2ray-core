package conf

import (
	"strings"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
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

func (v *ShadowsocksServerConfig) Build() (*serial.TypedMessage, error) {
	config := new(shadowsocks.ServerConfig)
	config.UdpEnabled = v.UDP

	if len(v.Password) == 0 {
		return nil, errors.New("Shadowsocks password is not specified.")
	}
	account := &shadowsocks.Account{
		Password: v.Password,
		Ota:      shadowsocks.Account_Auto,
	}
	if v.OTA != nil {
		if *v.OTA {
			account.Ota = shadowsocks.Account_Enabled
		} else {
			account.Ota = shadowsocks.Account_Disabled
		}
	}
	cipher := strings.ToLower(v.Cipher)
	switch cipher {
	case "aes-256-cfb":
		account.CipherType = shadowsocks.CipherType_AES_256_CFB
	case "aes-128-cfb":
		account.CipherType = shadowsocks.CipherType_AES_128_CFB
	case "chacha20":
		account.CipherType = shadowsocks.CipherType_CHACHA20
	case "chacha20-ietf":
		account.CipherType = shadowsocks.CipherType_CHACHA20_IETF
	default:
		return nil, errors.New("Unknown cipher method: " + cipher)
	}

	config.User = &protocol.User{
		Email:   v.Email,
		Level:   uint32(v.Level),
		Account: serial.ToTypedMessage(account),
	}

	return serial.ToTypedMessage(config), nil
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

func (v *ShadowsocksClientConfig) Build() (*serial.TypedMessage, error) {
	config := new(shadowsocks.ClientConfig)

	if len(v.Servers) == 0 {
		return nil, errors.New("0 Shadowsocks server configured.")
	}

	serverSpecs := make([]*protocol.ServerEndpoint, len(v.Servers))
	for idx, server := range v.Servers {
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
			account.CipherType = shadowsocks.CipherType_CHACHA20_IETF
		default:
			return nil, errors.New("Unknown cipher method: " + cipher)
		}

		ss := &protocol.ServerEndpoint{
			Address: server.Address.Build(),
			Port:    uint32(server.Port),
			User: []*protocol.User{
				{
					Email:   server.Email,
					Account: serial.ToTypedMessage(account),
				},
			},
		}

		serverSpecs[idx] = ss
	}

	config.Server = serverSpecs

	return serial.ToTypedMessage(config), nil
}
