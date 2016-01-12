package protocol_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	. "github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	protocoltesting "github.com/v2ray/v2ray-core/proxy/vmess/protocol/testing"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

type FakeTimestampGenerator struct {
	timestamp Timestamp
}

func (this *FakeTimestampGenerator) Next() Timestamp {
	return this.timestamp
}

type TestUser struct {
	id    *vmess.ID
	level vmess.UserLevel
}

func (u *TestUser) ID() *vmess.ID {
	return u.id
}

func (this *TestUser) Level() vmess.UserLevel {
	return this.level
}

func (this *TestUser) AlterIDs() []*vmess.ID {
	return nil
}

func (this *TestUser) AnyValidID() *vmess.ID {
	return this.id
}

func TestVMessSerialization(t *testing.T) {
	v2testing.Current(t)

	id, err := uuid.ParseString("2b2966ac-16aa-4fbf-8d81-c5f172a3da51")
	assert.Error(err).IsNil()

	userId := vmess.NewID(id)

	testUser := &TestUser{
		id: userId,
	}

	userSet := protocoltesting.MockUserSet{[]vmess.User{}, make(map[string]int), make(map[string]Timestamp)}
	userSet.AddUser(testUser)

	request := new(VMessRequest)
	request.Version = byte(0x01)
	request.User = testUser

	randBytes := make([]byte, 36)
	_, err = rand.Read(randBytes)
	assert.Error(err).IsNil()
	request.RequestIV = randBytes[:16]
	request.RequestKey = randBytes[16:32]
	request.ResponseHeader = randBytes[32:]

	request.Command = byte(0x01)
	request.Address = v2net.DomainAddress("v2ray.com")
	request.Port = v2net.Port(80)

	mockTime := Timestamp(1823730)

	buffer, err := request.ToBytes(&FakeTimestampGenerator{timestamp: mockTime}, nil)
	if err != nil {
		t.Fatal(err)
	}

	userSet.UserHashes[string(buffer.Value[:16])] = 0
	userSet.Timestamps[string(buffer.Value[:16])] = mockTime

	requestReader := NewVMessRequestReader(&userSet)
	actualRequest, err := requestReader.Read(bytes.NewReader(buffer.Value))
	if err != nil {
		t.Fatal(err)
	}

	assert.Byte(actualRequest.Version).Named("Version").Equals(byte(0x01))
	assert.String(actualRequest.User.ID()).Named("UserId").Equals(request.User.ID().String())
	assert.Bytes(actualRequest.RequestIV).Named("RequestIV").Equals(request.RequestIV[:])
	assert.Bytes(actualRequest.RequestKey).Named("RequestKey").Equals(request.RequestKey[:])
	assert.Bytes(actualRequest.ResponseHeader).Named("ResponseHeader").Equals(request.ResponseHeader[:])
	assert.Byte(actualRequest.Command).Named("Command").Equals(request.Command)
	assert.String(actualRequest.Address).Named("Address").Equals(request.Address.String())
}

func TestReadSingleByte(t *testing.T) {
	v2testing.Current(t)

	reader := NewVMessRequestReader(nil)
	_, err := reader.Read(bytes.NewReader(make([]byte, 1)))
	assert.Error(err).Equals(io.EOF)
}

func BenchmarkVMessRequestWriting(b *testing.B) {
	id, err := uuid.ParseString("2b2966ac-16aa-4fbf-8d81-c5f172a3da51")
	assert.Error(err).IsNil()

	userId := vmess.NewID(id)
	userSet := protocoltesting.MockUserSet{[]vmess.User{}, make(map[string]int), make(map[string]Timestamp)}

	testUser := &TestUser{
		id: userId,
	}
	userSet.AddUser(testUser)

	request := new(VMessRequest)
	request.Version = byte(0x01)
	request.User = testUser

	randBytes := make([]byte, 36)
	rand.Read(randBytes)
	request.RequestIV = randBytes[:16]
	request.RequestKey = randBytes[16:32]
	request.ResponseHeader = randBytes[32:]

	request.Command = byte(0x01)
	request.Address = v2net.DomainAddress("v2ray.com")
	request.Port = v2net.Port(80)

	for i := 0; i < b.N; i++ {
		request.ToBytes(NewRandomTimestampGenerator(Timestamp(time.Now().Unix()), 30), nil)
	}
}
