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

	config := new(ServerConfig)
	err := json.Unmarshal([]byte(rawJson), config)
	assert.Error(err).IsNil()

	account := new(Account)
	_, err = config.User.GetTypedAccount(account)
	assert.Error(err).IsNil()

	cipher, err := account.GetCipher()
	assert.Error(err).IsNil()
	assert.Int(cipher.KeySize()).Equals(16)
	assert.Bytes(account.GetCipherKey()).Equals([]byte{160, 224, 26, 2, 22, 110, 9, 80, 65, 52, 80, 20, 38, 243, 224, 241})
}
