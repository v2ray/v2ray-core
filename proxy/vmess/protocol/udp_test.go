package protocol

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
	"github.com/v2ray/v2ray-core/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestVMessUDPReadWrite(t *testing.T) {
	assert := unit.Assert(t)

	userId, err := user.NewID("2b2966ac-16aa-4fbf-8d81-c5f172a3da51")
	assert.Error(err).IsNil()

	userSet := mocks.MockUserSet{[]user.ID{}, make(map[string]int), make(map[string]int64)}
	userSet.AddUser(user.User{userId})

	message := &VMessUDP{
		user:    userId,
		version: byte(0x01),
		token:   1234,
		address: v2net.DomainAddress("v2ray.com", 8372),
		data:    []byte("An UDP message."),
	}

	mockTime := int64(1823730)
	buffer := message.ToBytes(user.NewTimeHash(user.HMACHash{}), func(base int64, delta int) int64 { return mockTime }, nil)

	userSet.UserHashes[string(buffer[:16])] = 0
	userSet.Timestamps[string(buffer[:16])] = mockTime

	messageRestored, err := ReadVMessUDP(buffer, &userSet)
	assert.Error(err).IsNil()

	assert.String(messageRestored.user.String).Equals(message.user.String)
	assert.Byte(messageRestored.version).Equals(message.version)
	assert.Uint16(messageRestored.token).Equals(message.token)
	assert.String(messageRestored.address.String()).Equals(message.address.String())
	assert.Bytes(messageRestored.data).Equals(message.data)
}
