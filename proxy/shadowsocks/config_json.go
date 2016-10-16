// +build json

package shadowsocks

import (
	"encoding/json"
	"errors"
	"strings"

	"v2ray.com/core/common"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/protocol"
)

func (this *ServerConfig) UnmarshalJSON(data []byte) error {
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

	if len(jsonConfig.Password) == 0 {
		log.Error("Shadowsocks: Password is not specified.")
		return common.ErrBadConfiguration
	}
	account := &Account{
		Password: jsonConfig.Password,
	}
	jsonConfig.Cipher = strings.ToLower(jsonConfig.Cipher)
	switch jsonConfig.Cipher {
	case "aes-256-cfb":
		account.CipherType = CipherType_AES_256_CFB
	case "aes-128-cfb":
		account.CipherType = CipherType_AES_128_CFB
	case "chacha20":
		account.CipherType = CipherType_CHACHA20
	case "chacha20-ietf":
		account.CipherType = CipherType_CHACHA20_IEFT
	default:
		log.Error("Shadowsocks: Unknown cipher method: ", jsonConfig.Cipher)
		return common.ErrBadConfiguration
	}

	this.User = &protocol.User{
		Email:   jsonConfig.Email,
		Level:   uint32(jsonConfig.Level),
		Account: loader.NewTypedSettings(account),
	}

	return nil
}
