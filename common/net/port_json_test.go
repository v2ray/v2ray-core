// +build json

package net_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestIntPort(t *testing.T) {
	v2testing.Current(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("1234"), &portRange)
	assert.Error(err).IsNil()

	assert.Uint16(portRange.From.Value()).Equals(uint16(1234))
	assert.Uint16(portRange.To.Value()).Equals(uint16(1234))
}

func TestOverRangeIntPort(t *testing.T) {
	v2testing.Current(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("70000"), &portRange)
	assert.Error(err).Equals(InvalidPortRange)

	err = json.Unmarshal([]byte("-1"), &portRange)
	assert.Error(err).Equals(InvalidPortRange)
}

func TestSingleStringPort(t *testing.T) {
	v2testing.Current(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234\""), &portRange)
	assert.Error(err).IsNil()

	assert.Uint16(portRange.From.Value()).Equals(uint16(1234))
	assert.Uint16(portRange.To.Value()).Equals(uint16(1234))
}

func TestStringPairPort(t *testing.T) {
	v2testing.Current(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234-5678\""), &portRange)
	assert.Error(err).IsNil()

	assert.Uint16(portRange.From.Value()).Equals(uint16(1234))
	assert.Uint16(portRange.To.Value()).Equals(uint16(5678))
}

func TestOverRangeStringPort(t *testing.T) {
	v2testing.Current(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"65536\""), &portRange)
	assert.Error(err).Equals(InvalidPortRange)

	err = json.Unmarshal([]byte("\"70000-80000\""), &portRange)
	assert.Error(err).Equals(InvalidPortRange)

	err = json.Unmarshal([]byte("\"1-90000\""), &portRange)
	assert.Error(err).Equals(InvalidPortRange)

	err = json.Unmarshal([]byte("\"700-600\""), &portRange)
	assert.Error(err).Equals(InvalidPortRange)
}
