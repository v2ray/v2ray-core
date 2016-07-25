package encoding_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	. "github.com/v2ray/v2ray-core/proxy/vmess/encoding"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestRequestSerialization(t *testing.T) {
	assert := assert.On(t)

	user := protocol.NewUser(
		protocol.UserLevelUntrusted,
		"test@v2ray.com")
	user.Account = &vmess.Account{
		ID:       protocol.NewID(uuid.New()),
		AlterIDs: nil,
	}

	expectedRequest := &protocol.RequestHeader{
		Version: 1,
		User:    user,
		Command: protocol.RequestCommandTCP,
		Option:  protocol.RequestOption(0),
		Address: v2net.DomainAddress("www.v2ray.com"),
		Port:    v2net.Port(443),
	}

	buffer := alloc.NewBuffer().Clear()
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
}
