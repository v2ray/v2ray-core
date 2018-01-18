package encoding_test

import (
	"context"
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/vmess"
	. "v2ray.com/core/proxy/vmess/encoding"
	. "v2ray.com/ext/assert"
)

func TestRequestSerialization(t *testing.T) {
	assert := With(t)

	user := &protocol.User{
		Level: 0,
		Email: "test@v2ray.com",
	}
	id := uuid.New()
	account := &vmess.Account{
		Id:      id.String(),
		AlterId: 0,
	}
	user.Account = serial.ToTypedMessage(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommandTCP,
		Address:  net.DomainAddress("www.v2ray.com"),
		Port:     net.Port(443),
		Security: protocol.Security(protocol.SecurityType_AES128_GCM),
	}

	buffer := buf.New()
	client := NewClientSession(protocol.DefaultIDHash)
	common.Must(client.EncodeRequestHeader(expectedRequest, buffer))

	buffer2 := buf.New()
	buffer2.Append(buffer.Bytes())

	ctx, cancel := context.WithCancel(context.Background())
	sessionHistory := NewSessionHistory(ctx)

	userValidator := vmess.NewTimedUserValidator(ctx, protocol.DefaultIDHash)
	userValidator.Add(user)

	server := NewServerSession(userValidator, sessionHistory)
	actualRequest, err := server.DecodeRequestHeader(buffer)
	assert(err, IsNil)

	assert(expectedRequest.Version, Equals, actualRequest.Version)
	assert(byte(expectedRequest.Command), Equals, byte(actualRequest.Command))
	assert(byte(expectedRequest.Option), Equals, byte(actualRequest.Option))
	assert(expectedRequest.Address, Equals, actualRequest.Address)
	assert(expectedRequest.Port, Equals, actualRequest.Port)
	assert(byte(expectedRequest.Security), Equals, byte(actualRequest.Security))

	_, err = server.DecodeRequestHeader(buffer2)
	// anti replay attack
	assert(err, IsNotNil)

	cancel()
}
