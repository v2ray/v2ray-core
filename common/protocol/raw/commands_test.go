package raw_test

import (
	"bytes"
	"testing"

	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	"github.com/v2ray/v2ray-core/common/protocol"
	. "github.com/v2ray/v2ray-core/common/protocol/raw"
	"github.com/v2ray/v2ray-core/common/uuid"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSwitchAccount(t *testing.T) {
	v2testing.Current(t)

	sa := &protocol.CommandSwitchAccount{
		Port:     1234,
		ID:       uuid.New(),
		AlterIds: 1024,
		Level:    128,
		ValidMin: 16,
	}

	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	err := MarshalCommand(sa, buffer)
	assert.Error(err).IsNil()

	cmd, err := UnmarshalCommand(1, buffer.Bytes())
	assert.Error(err).IsNil()

	sa2, ok := cmd.(*protocol.CommandSwitchAccount)
	assert.Bool(ok).IsTrue()
	assert.Pointer(sa.Host).IsNil()
	assert.Pointer(sa2.Host).IsNil()
	netassert.Port(sa.Port).Equals(sa2.Port)
	assert.String(sa.ID).Equals(sa2.ID.String())
	assert.Uint16(sa.AlterIds.Value()).Equals(sa2.AlterIds.Value())
	assert.Byte(byte(sa.Level)).Equals(byte(sa2.Level))
	assert.Byte(sa.ValidMin).Equals(sa2.ValidMin)
}
