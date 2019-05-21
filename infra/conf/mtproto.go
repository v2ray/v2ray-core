package conf

import (
	"encoding/hex"
	"encoding/json"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/mtproto"
)

type MTProtoAccount struct {
	Secret string `json:"secret"`
}

// Build implements Buildable
func (a *MTProtoAccount) Build() (*mtproto.Account, error) {
	if len(a.Secret) != 32 {
		return nil, newError("MTProto secret must have 32 chars")
	}
	secret, err := hex.DecodeString(a.Secret)
	if err != nil {
		return nil, newError("failed to decode secret: ", a.Secret).Base(err)
	}
	return &mtproto.Account{
		Secret: secret,
	}, nil
}

type MTProtoServerConfig struct {
	Users []json.RawMessage `json:"users"`
}

func (c *MTProtoServerConfig) Build() (proto.Message, error) {
	config := &mtproto.ServerConfig{}

	if len(c.Users) == 0 {
		return nil, newError("zero MTProto users configured.")
	}
	config.User = make([]*protocol.User, len(c.Users))
	for idx, rawData := range c.Users {
		user := new(protocol.User)
		if err := json.Unmarshal(rawData, user); err != nil {
			return nil, newError("invalid MTProto user").Base(err)
		}
		account := new(MTProtoAccount)
		if err := json.Unmarshal(rawData, account); err != nil {
			return nil, newError("invalid MTProto user").Base(err)
		}
		accountProto, err := account.Build()
		if err != nil {
			return nil, newError("failed to parse MTProto user").Base(err)
		}
		user.Account = serial.ToTypedMessage(accountProto)
		config.User[idx] = user
	}

	return config, nil
}

type MTProtoClientConfig struct {
}

func (c *MTProtoClientConfig) Build() (proto.Message, error) {
	config := new(mtproto.ClientConfig)
	return config, nil
}
