package net_test

import (
	"testing"

	. "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestPortRangeContains(t *testing.T) {
	assert := assert.On(t)

	portRange := &PortRange{
		From: Port(53),
		To:   Port(53),
	}
	assert.Bool(portRange.Contains(Port(53))).IsTrue()
}
