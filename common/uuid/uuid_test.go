package uuid_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/uuid"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestParseBytes(t *testing.T) {
	v2testing.Current(t)

	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	bytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseBytes(bytes)
	assert.Error(err).IsNil()
	assert.String(uuid).Equals(str)
}

func TestParseString(t *testing.T) {
	v2testing.Current(t)

	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	expectedBytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseString(str)
	assert.Error(err).IsNil()
	assert.Bytes(uuid.Bytes()).Equals(expectedBytes)
}

func TestNewUUID(t *testing.T) {
	v2testing.Current(t)

	uuid := New()
	uuid2, err := ParseString(uuid.String())

	assert.Error(err).IsNil()
	assert.StringLiteral(uuid.String()).Equals(uuid2.String())
	assert.Bytes(uuid.Bytes()).Equals(uuid2.Bytes())
}

func TestRandom(t *testing.T) {
	v2testing.Current(t)

	uuid := New()
	uuid2 := New()

	assert.StringLiteral(uuid.String()).NotEquals(uuid2.String())
	assert.Bytes(uuid.Bytes()).NotEquals(uuid2.Bytes())
}

func TestEquals(t *testing.T) {
	v2testing.Current(t)

	var uuid *UUID = nil
	var uuid2 *UUID = nil
	assert.Bool(uuid.Equals(uuid2)).IsTrue()
	assert.Bool(uuid.Equals(New())).IsFalse()
}
