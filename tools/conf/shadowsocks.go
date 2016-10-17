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
}

func (this *ShadowsocksServerConfig) Build() (*loader.TypedSettings, error) {
	config := new(shadowsocks.ServerConfig)
	config.UdpEnabled = this.UDP

	if len(this.Password) == 0 {
		return nil, errors.New("Shadowsocks password is not specified.")
	}
	account := &shadowsocks.Account{
		Password: this.Password,
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
