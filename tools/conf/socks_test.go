package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/socks"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/tools/conf"
)

func TestSocksInboundConfig(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "auth": "password",
    "accounts": [
      {
        "user": "my-username",
        "pass": "my-password"
      }
    ],
    "udp": false,
    "ip": "127.0.0.1",
    "timeout": 5
  }`

	config := new(SocksServerConfig)
	err := json.Unmarshal([]byte(rawJson), &config)
	assert.Error(err).IsNil()

	message, err := config.Build()
	assert.Error(err).IsNil()

	iConfig, err := message.GetInstance()
	assert.Error(err).IsNil()

	socksConfig := iConfig.(*socks.ServerConfig)
	assert.Bool(socksConfig.AuthType == socks.AuthType_PASSWORD).IsTrue()
	assert.Int(len(socksConfig.Accounts)).Equals(1)
	assert.String(socksConfig.Accounts["my-username"]).Equals("my-password")
	assert.Bool(socksConfig.UdpEnabled).IsFalse()
	assert.Address(socksConfig.Address.AsAddress()).Equals(net.LocalHostIP)
	assert.Uint32(socksConfig.Timeout).Equals(5)
}

func TestSocksOutboundConfig(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "servers": [{
      "address": "127.0.0.1",
      "port": 1234,
      "users": [
        {"user": "test user", "pass": "test pass", "email": "test@email.com"}
      ]
    }]
  }`

	config := new(SocksClientConfig)
	err := json.Unmarshal([]byte(rawJson), &config)
	assert.Error(err).IsNil()

	message, err := config.Build()
	assert.Error(err).IsNil()

	iConfig, err := message.GetInstance()
	assert.Error(err).IsNil()

	socksConfig := iConfig.(*socks.ClientConfig)
	assert.Int(len(socksConfig.Server)).Equals(1)

	ss := protocol.NewServerSpecFromPB(*socksConfig.Server[0])
	assert.Destination(ss.Destination()).EqualsString("tcp:127.0.0.1:1234")

	user := ss.PickUser()
	assert.String(user.Email).Equals("test@email.com")

	account, err := user.GetTypedAccount()
	assert.Error(err).IsNil()

	socksAccount := account.(*socks.Account)
	assert.String(socksAccount.Username).Equals("test user")
	assert.String(socksAccount.Password).Equals("test pass")
}
