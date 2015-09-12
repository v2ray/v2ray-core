package vmess

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/net"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestVMessSerialization(t *testing.T) {
	assert := unit.Assert(t)

	userId, err := core.UUIDToID("2b2966ac-16aa-4fbf-8d81-c5f172a3da51")
	if err != nil {
		t.Fatal(err)
	}

	userSet := core.NewUserSet()
	userSet.AddUser(core.User{userId})

	request := new(VMessRequest)
	request.Version = byte(0x01)
	request.UserId = userId

	_, err = rand.Read(request.RequestIV[:])
	if err != nil {
		t.Fatal(err)
	}

	_, err = rand.Read(request.RequestKey[:])
	if err != nil {
		t.Fatal(err)
	}

	_, err = rand.Read(request.ResponseHeader[:])
	if err != nil {
		t.Fatal(err)
	}

	request.Command = byte(0x01)
	request.Address = v2net.DomainAddress("v2ray.com", 80)

	buffer := bytes.NewBuffer(make([]byte, 0, 300))
	requestWriter := NewVMessRequestWriter()
	err = requestWriter.Write(buffer, request)
	if err != nil {
		t.Fatal(err)
	}

	requestReader := NewVMessRequestReader(userSet)
	actualRequest, err := requestReader.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	assert.Byte(actualRequest.Version).Named("Version").Equals(byte(0x01))
	assert.Bytes(actualRequest.UserId[:]).Named("UserId").Equals(request.UserId[:])
	assert.Bytes(actualRequest.RequestIV[:]).Named("RequestIV").Equals(request.RequestIV[:])
	assert.Bytes(actualRequest.RequestKey[:]).Named("RequestKey").Equals(request.RequestKey[:])
	assert.Bytes(actualRequest.ResponseHeader[:]).Named("ResponseHeader").Equals(request.ResponseHeader[:])
	assert.Byte(actualRequest.Command).Named("Command").Equals(request.Command)
	assert.String(actualRequest.Address.String()).Named("Address").Equals(request.Address.String())
}
