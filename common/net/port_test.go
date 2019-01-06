package net_test

import (
	"testing"

	. "v2ray.com/core/common/net"
	. "v2ray.com/ext/assert"
)

func TestPortRangeContains(t *testing.T) {
	assert := With(t)

	portRange := &PortRange{
		From: 53,
		To:   53,
	}
	assert(portRange.Contains(Port(53)), IsTrue)
}
