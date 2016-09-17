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

	assert.Int(config.GetCipher().KeySize()).Equals(16)
	account, err := config.User.GetTypedAccount(&Account{})
	assert.Error(err).IsNil()
	assert.Bytes(account.(*Account).GetCipherKey(config.GetCipher().KeySize())).Equals([]byte{160, 224, 26, 2, 22, 110, 9, 80, 65, 52, 80, 20, 38, 243, 224, 241})
}
