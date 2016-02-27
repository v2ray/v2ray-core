package raw_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	"github.com/v2ray/v2ray-core/common/protocol"
	. "github.com/v2ray/v2ray-core/common/protocol/raw"
	"github.com/v2ray/v2ray-core/common/uuid"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestRequestSerialization(t *testing.T) {
	v2testing.Current(t)

	user := protocol.NewUser(
		protocol.NewID(uuid.New()),
		protocol.UserLevelUntrusted,
		0,
		"test@v2ray.com")

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

	userValidator := protocol.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)

	server := NewServerSession(userValidator)
	actualRequest, err := server.DecodeRequestHeader(buffer)
	assert.Error(err).IsNil()

	assert.Byte(expectedRequest.Version).Equals(actualRequest.Version)
	assert.Byte(byte(expectedRequest.Command)).Equals(byte(actualRequest.Command))
	assert.Byte(byte(expectedRequest.Option)).Equals(byte(actualRequest.Option))
	netassert.Address(expectedRequest.Address).Equals(actualRequest.Address)
	netassert.Port(expectedRequest.Port).Equals(actualRequest.Port)
}
