package vmess

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/net"
	"github.com/v2ray/v2ray-core/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestVMessSerialization(t *testing.T) {
	assert := unit.Assert(t)

	userId, err := core.NewID("2b2966ac-16aa-4fbf-8d81-c5f172a3da51")
	if err != nil {
		t.Fatal(err)
	}

	userSet := mocks.MockUserSet{[]core.ID{}, make(map[string]int)}
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

	userSet.UserHashes[string(buffer.Bytes()[:16])] = 0

	requestReader := NewVMessRequestReader(&userSet)
	actualRequest, err := requestReader.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	assert.Byte(actualRequest.Version).Named("Version").Equals(byte(0x01))
	assert.String(actualRequest.UserId.String).Named("UserId").Equals(request.UserId.String)
	assert.Bytes(actualRequest.RequestIV[:]).Named("RequestIV").Equals(request.RequestIV[:])
	assert.Bytes(actualRequest.RequestKey[:]).Named("RequestKey").Equals(request.RequestKey[:])
	assert.Bytes(actualRequest.ResponseHeader[:]).Named("ResponseHeader").Equals(request.ResponseHeader[:])
	assert.Byte(actualRequest.Command).Named("Command").Equals(request.Command)
	assert.String(actualRequest.Address.String()).Named("Address").Equals(request.Address.String())
}

func BenchmarkVMessRequestWriting(b *testing.B) {
	userId, _ := core.NewID("2b2966ac-16aa-4fbf-8d81-c5f172a3da51")
	userSet := mocks.MockUserSet{[]core.ID{}, make(map[string]int)}
	userSet.AddUser(core.User{userId})

	request := new(VMessRequest)
	request.Version = byte(0x01)
	request.UserId = userId

	rand.Read(request.RequestIV[:])
	rand.Read(request.RequestKey[:])
	rand.Read(request.ResponseHeader[:])

	request.Command = byte(0x01)
	request.Address = v2net.DomainAddress("v2ray.com", 80)

	requestWriter := NewVMessRequestWriter()
	for i := 0; i < b.N; i++ {
		requestWriter.Write(ioutil.Discard, request)
	}
}
