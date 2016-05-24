package net_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestPortRangeContains(t *testing.T) {
	assert := assert.On(t)

	portRange := &PortRange{
		From: Port(53),
		To:   Port(53),
	}
	assert.Bool(portRange.Contains(Port(53))).IsTrue()
}
