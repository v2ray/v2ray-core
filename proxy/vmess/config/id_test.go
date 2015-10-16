package config

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestUUIDToID(t *testing.T) {
	assert := unit.Assert(t)

	uuid := "2418d087-648d-4990-86e8-19dca1d006d3"
	expectedBytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	actualBytes, _ := NewID(uuid)
	assert.Bytes(actualBytes.Bytes[:]).Named("UUID").Equals(expectedBytes)
}
