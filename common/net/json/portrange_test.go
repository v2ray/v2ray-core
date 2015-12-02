package json

import (
	"encoding/json"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestIntPort(t *testing.T) {
	assert := unit.Assert(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("1234"), &portRange)
	assert.Error(err).IsNil()

	assert.Uint16(portRange.from.Value()).Equals(uint16(1234))
	assert.Uint16(portRange.to.Value()).Equals(uint16(1234))
}

func TestOverRangeIntPort(t *testing.T) {
	assert := unit.Assert(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("70000"), &portRange)
	assert.Error(err).Equals(InvalidPortRange)

	err = json.Unmarshal([]byte("-1"), &portRange)
	assert.Error(err).Equals(InvalidPortRange)
}

func TestSingleStringPort(t *testing.T) {
	assert := unit.Assert(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234\""), &portRange)
	assert.Error(err).IsNil()

	assert.Uint16(portRange.from.Value()).Equals(uint16(1234))
	assert.Uint16(portRange.to.Value()).Equals(uint16(1234))
}

func TestStringPairPort(t *testing.T) {
	assert := unit.Assert(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234-5678\""), &portRange)
	assert.Error(err).IsNil()

	assert.Uint16(portRange.from.Value()).Equals(uint16(1234))
	assert.Uint16(portRange.to.Value()).Equals(uint16(5678))
}

func TestOverRangeStringPort(t *testing.T) {
	assert := unit.Assert(t)

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
