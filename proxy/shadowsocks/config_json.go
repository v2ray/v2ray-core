// +build json

package shadowsocks

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Cipher   serial.StringLiteral `json:"method"`
		Password serial.StringLiteral `json:"password"`
		UDP      bool                 `json:"udp"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	if len(jsonConfig.Password) == 0 {
		log.Error("Shadowsocks: Password is not specified.")
		return internal.ErrorBadConfiguration
	}
	this.UDP = jsonConfig.UDP
	this.Password = jsonConfig.Password.String()
	if this.Cipher == nil {
		log.Error("Shadowsocks: Cipher method is not specified.")
		return internal.ErrorBadConfiguration
	}
	jsonConfig.Cipher = jsonConfig.Cipher.ToLower()
	switch jsonConfig.Cipher.String() {
	case "aes-256-cfb":
		this.Cipher = &AesCfb{
			KeyBytes: 32,
		}
	case "aes-128-cfb":
		this.Cipher = &AesCfb{
			KeyBytes: 32,
		}
	default:
		log.Error("Shadowsocks: Unknown cipher method: ", jsonConfig.Cipher)
		return internal.ErrorBadConfiguration
	}
	return nil
}
