package command_test

import (
	"bytes"
	"testing"

	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	"github.com/v2ray/v2ray-core/common/uuid"
	. "github.com/v2ray/v2ray-core/proxy/vmess/command"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSwitchAccount(t *testing.T) {
	v2testing.Current(t)

	sa := &SwitchAccount{
		Port:     1234,
		ID:       uuid.New(),
		AlterIds: 1024,
		Level:    128,
		ValidMin: 16,
	}

	cmd, err := CreateResponseCommand(1)
	assert.Error(err).IsNil()

	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	sa.Marshal(buffer)

	cmd.Unmarshal(buffer.Bytes())
	sa2, ok := cmd.(*SwitchAccount)
	assert.Bool(ok).IsTrue()
	netassert.Port(sa.Port).Equals(sa2.Port)
	assert.String(sa.ID).Equals(sa2.ID.String())
	assert.Uint16(sa.AlterIds.Value()).Equals(sa2.AlterIds.Value())
	assert.Byte(byte(sa.Level)).Equals(byte(sa2.Level))
	assert.Byte(sa.ValidMin).Equals(sa2.ValidMin)
}
