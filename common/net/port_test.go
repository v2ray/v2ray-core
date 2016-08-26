package net_test

import (
	"testing"

	. "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestPortRangeContains(t *testing.T) {
	assert := assert.On(t)

	portRange := &PortRange{
		From: 53,
		To:   53,
	}
	assert.Bool(portRange.Contains(Port(53))).IsTrue()
}
