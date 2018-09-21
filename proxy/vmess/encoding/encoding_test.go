package encoding_test

import (
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/vmess"
	. "v2ray.com/core/proxy/vmess/encoding"
	. "v2ray.com/ext/assert"
)

func toAccount(a *vmess.Account) protocol.Account {
	account, err := a.AsAccount()
	common.Must(err)
	return account
}

func TestRequestSerialization(t *testing.T) {
	assert := With(t)

	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@v2ray.com",
	}
	id := uuid.New()
	account := &vmess.Account{
		Id:      id.String(),
		AlterId: 0,
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommandTCP,
		Address:  net.DomainAddress("www.v2ray.com"),
		Port:     net.Port(443),
		Security: protocol.SecurityType_AES128_GCM,
	}

	buffer := buf.New()
	client := NewClientSession(protocol.DefaultIDHash)
	common.Must(client.EncodeRequestHeader(expectedRequest, buffer))

	buffer2 := buf.New()
	buffer2.Write(buffer.Bytes())

	sessionHistory := NewSessionHistory()
	defer common.Close(sessionHistory)

	userValidator := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)
	defer common.Close(userValidator)

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
}

func TestInvalidRequest(t *testing.T) {
	assert := With(t)

	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@v2ray.com",
	}
	id := uuid.New()
	account := &vmess.Account{
		Id:      id.String(),
		AlterId: 0,
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommand(100),
		Address:  net.DomainAddress("www.v2ray.com"),
		Port:     net.Port(443),
		Security: protocol.SecurityType_AES128_GCM,
	}

	buffer := buf.New()
	client := NewClientSession(protocol.DefaultIDHash)
	common.Must(client.EncodeRequestHeader(expectedRequest, buffer))

	buffer2 := buf.New()
	buffer2.Write(buffer.Bytes())

	sessionHistory := NewSessionHistory()
	defer common.Close(sessionHistory)

	userValidator := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)
	defer common.Close(userValidator)

	server := NewServerSession(userValidator, sessionHistory)
	_, err := server.DecodeRequestHeader(buffer)
	assert(err, IsNotNil)
}

func TestMuxRequest(t *testing.T) {
	assert := With(t)

	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@v2ray.com",
	}
	id := uuid.New()
	account := &vmess.Account{
		Id:      id.String(),
		AlterId: 0,
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommandMux,
		Security: protocol.SecurityType_AES128_GCM,
	}

	buffer := buf.New()
	client := NewClientSession(protocol.DefaultIDHash)
	common.Must(client.EncodeRequestHeader(expectedRequest, buffer))

	buffer2 := buf.New()
	buffer2.Write(buffer.Bytes())

	sessionHistory := NewSessionHistory()
	defer common.Close(sessionHistory)

	userValidator := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)
	defer common.Close(userValidator)

	server := NewServerSession(userValidator, sessionHistory)
	actualRequest, err := server.DecodeRequestHeader(buffer)
	assert(err, IsNil)

	assert(expectedRequest.Version, Equals, actualRequest.Version)
	assert(byte(expectedRequest.Command), Equals, byte(actualRequest.Command))
	assert(byte(expectedRequest.Option), Equals, byte(actualRequest.Option))
	assert(byte(expectedRequest.Security), Equals, byte(actualRequest.Security))
}
