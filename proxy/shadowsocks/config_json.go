// +build json

package shadowsocks

import (
	"encoding/json"
	"errors"
	"strings"

	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/registry"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Cipher   string `json:"method"`
		Password string `json:"password"`
		UDP      bool   `json:"udp"`
		Level    byte   `json:"level"`
		Email    string `json:"email"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Shadowsocks: Failed to parse config: " + err.Error())
	}

	this.UDP = jsonConfig.UDP
	jsonConfig.Cipher = strings.ToLower(jsonConfig.Cipher)
	switch jsonConfig.Cipher {
	case "aes-256-cfb":
		this.Cipher = &AesCfb{
			KeyBytes: 32,
		}
	case "aes-128-cfb":
		this.Cipher = &AesCfb{
			KeyBytes: 16,
		}
	case "chacha20":
		this.Cipher = &ChaCha20{
			IVBytes: 8,
		}
	case "chacha20-ietf":
		this.Cipher = &ChaCha20{
			IVBytes: 12,
		}
	default:
		log.Error("Shadowsocks: Unknown cipher method: ", jsonConfig.Cipher)
		return common.ErrBadConfiguration
	}

	if len(jsonConfig.Password) == 0 {
		log.Error("Shadowsocks: Password is not specified.")
		return common.ErrBadConfiguration
	}
	this.Key = PasswordToCipherKey(jsonConfig.Password, this.Cipher.KeySize())

	this.Level = protocol.UserLevel(jsonConfig.Level)
	this.Email = jsonConfig.Email

	return nil
}

func init() {
	registry.RegisterInboundConfig("shadowsocks", func() interface{} { return new(Config) })
}
