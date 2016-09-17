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

	"github.com/golang/protobuf/ptypes"
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

	this.UdpEnabled = jsonConfig.UDP
	jsonConfig.Cipher = strings.ToLower(jsonConfig.Cipher)
	switch jsonConfig.Cipher {
	case "aes-256-cfb":
		this.Cipher = Config_AES_256_CFB
	case "aes-128-cfb":
		this.Cipher = Config_AES_128_CFB
	case "chacha20":
		this.Cipher = Config_CHACHA20
	case "chacha20-ietf":
		this.Cipher = Config_CHACHA20_IEFT
	default:
		log.Error("Shadowsocks: Unknown cipher method: ", jsonConfig.Cipher)
		return common.ErrBadConfiguration
	}

	if len(jsonConfig.Password) == 0 {
		log.Error("Shadowsocks: Password is not specified.")
		return common.ErrBadConfiguration
	}
	account, err := ptypes.MarshalAny(&Account{
		Password: jsonConfig.Password,
	})
	if err != nil {
		log.Error("Shadowsocks: Failed to create account: ", err)
		return common.ErrBadConfiguration
	}
	this.User = &protocol.User{
		Email:   jsonConfig.Email,
		Level:   uint32(jsonConfig.Level),
		Account: account,
	}

	return nil
}

func init() {
	registry.RegisterInboundConfig("shadowsocks", func() interface{} { return new(Config) })
}
