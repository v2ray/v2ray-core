package uuid_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestParseBytes(t *testing.T) {
	assert := assert.On(t)

	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	bytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseBytes(bytes)
	assert.Error(err).IsNil()
	assert.String(uuid.String()).Equals(str)

	_, err = ParseBytes([]byte{1, 3, 2, 4})
	assert.Error(err).Equals(ErrInvalidID)
}

func TestParseString(t *testing.T) {
	assert := assert.On(t)

	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	expectedBytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseString(str)
	assert.Error(err).IsNil()
	assert.Bytes(uuid.Bytes()).Equals(expectedBytes)

	uuid, err = ParseString("2418d087")
	assert.Error(err).Equals(ErrInvalidID)

	uuid, err = ParseString("2418d087-648k-4990-86e8-19dca1d006d3")
	assert.Error(err).IsNotNil()
}

func TestNewUUID(t *testing.T) {
	assert := assert.On(t)

	uuid := New()
	uuid2, err := ParseString(uuid.String())

	assert.Error(err).IsNil()
	assert.String(uuid.String()).Equals(uuid2.String())
	assert.Bytes(uuid.Bytes()).Equals(uuid2.Bytes())
}

func TestRandom(t *testing.T) {
	assert := assert.On(t)

	uuid := New()
	uuid2 := New()

	assert.String(uuid.String()).NotEquals(uuid2.String())
	assert.Bytes(uuid.Bytes()).NotEquals(uuid2.Bytes())
}

func TestEquals(t *testing.T) {
	assert := assert.On(t)

	var uuid *UUID = nil
	var uuid2 *UUID = nil
	assert.Bool(uuid.Equals(uuid2)).IsTrue()
	assert.Bool(uuid.Equals(New())).IsFalse()
}

func TestNext(t *testing.T) {
	assert := assert.On(t)

	uuid := New()
	uuid2 := uuid.Next()
	assert.Bool(uuid.Equals(uuid2)).IsFalse()
}
