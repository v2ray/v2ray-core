package net_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestPortRangeContains(t *testing.T) {
	v2testing.Current(t)

	portRange := &PortRange{
		From: Port(53),
		To:   Port(53),
	}
	assert.Bool(portRange.Contains(Port(53))).IsTrue()
}
