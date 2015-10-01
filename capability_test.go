package core

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestCapabilities(t *testing.T) {
	assert := unit.Assert(t)

	caps := NewCapabilities()
	assert.Bool(caps.HasCapability(TCPConnection)).IsFalse()

	caps.AddCapability(TCPConnection)
	assert.Bool(caps.HasCapability(TCPConnection)).IsTrue()

	caps.AddCapability(UDPConnection)
	assert.Bool(caps.HasCapability(TCPConnection)).IsTrue()
	assert.Bool(caps.HasCapability(UDPConnection)).IsTrue()
}
