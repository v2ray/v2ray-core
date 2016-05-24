package protocol_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestCmdKey(t *testing.T) {
	assert := assert.On(t)

	id := NewID(uuid.New())
	assert.Bool(serial.BytesT(id.CmdKey()).All(0)).IsFalse()
}
