package protocol_test

import (
	"testing"

	"v2ray.com/core/common/predicate"
	. "v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/testing/assert"
)

func TestCmdKey(t *testing.T) {
	assert := assert.On(t)

	id := NewID(uuid.New())
	assert.Bool(predicate.BytesAll(id.CmdKey(), 0)).IsFalse()
}
