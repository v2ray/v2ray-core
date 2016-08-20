// +build json

package shadowsocks

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/testing/assert"
)

func TestConfigParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "method": "aes-128-cfb",
    "password": "v2ray-password"
  }`

	config := new(Config)
	err := json.Unmarshal([]byte(rawJson), config)
	assert.Error(err).IsNil()

	assert.Int(config.Cipher.KeySize()).Equals(16)
	assert.Bytes(config.Key).Equals([]byte{160, 224, 26, 2, 22, 110, 9, 80, 65, 52, 80, 20, 38, 243, 224, 241})
}
