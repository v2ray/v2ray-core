package encoding_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/vmess"
	. "v2ray.com/core/proxy/vmess/encoding"
	"v2ray.com/core/testing/assert"
)

func TestRequestSerialization(t *testing.T) {
	assert := assert.On(t)

	user := &protocol.User{
		Level: 0,
		Email: "test@v2ray.com",
	}
	account := &vmess.Account{
		Id:      uuid.New().String(),
		AlterId: 0,
	}
	user.Account = loader.NewTypedSettings(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommandTCP,
		Option:   protocol.RequestOptionConnectionReuse,
		Address:  v2net.DomainAddress("www.v2ray.com"),
		Port:     v2net.Port(443),
		Security: protocol.Security(protocol.SecurityType_AES128_GCM),
	}

	buffer := buf.New()
	client := NewClientSession(protocol.DefaultIDHash)
	client.EncodeRequestHeader(expectedRequest, buffer)

	userValidator := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)

	server := NewServerSession(userValidator)
	actualRequest, err := server.DecodeRequestHeader(buffer)
	assert.Error(err).IsNil()

	assert.Byte(expectedRequest.Version).Equals(actualRequest.Version)
	assert.Byte(byte(expectedRequest.Command)).Equals(byte(actualRequest.Command))
	assert.Byte(byte(expectedRequest.Option)).Equals(byte(actualRequest.Option))
	assert.Address(expectedRequest.Address).Equals(actualRequest.Address)
	assert.Port(expectedRequest.Port).Equals(actualRequest.Port)
	assert.Byte(byte(expectedRequest.Security)).Equals(byte(actualRequest.Security))
}
