package encoding_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
	. "v2ray.com/core/proxy/vmess/encoding"
	. "v2ray.com/ext/assert"
)

func TestSwitchAccount(t *testing.T) {
	assert := With(t)

	sa := &protocol.CommandSwitchAccount{
		Port:     1234,
		ID:       uuid.New(),
		AlterIds: 1024,
		Level:    128,
		ValidMin: 16,
	}

	buffer := buf.New()
	err := MarshalCommand(sa, buffer)
	assert(err, IsNil)

	cmd, err := UnmarshalCommand(1, buffer.BytesFrom(2))
	assert(err, IsNil)

	sa2, ok := cmd.(*protocol.CommandSwitchAccount)
	assert(ok, IsTrue)
	assert(sa.Host, IsNil)
	assert(sa2.Host, IsNil)
	assert(sa.Port, Equals, sa2.Port)
	assert(sa.ID.String(), Equals, sa2.ID.String())
	assert(sa.AlterIds, Equals, sa2.AlterIds)
	assert(byte(sa.Level), Equals, byte(sa2.Level))
	assert(sa.ValidMin, Equals, sa2.ValidMin)
}
