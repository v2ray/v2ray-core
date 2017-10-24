package uuid_test

import (
	"testing"

	. "v2ray.com/core/common/uuid"
	. "v2ray.com/ext/assert"
)

func TestParseBytes(t *testing.T) {
	assert := With(t)

	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	bytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseBytes(bytes)
	assert(err, IsNil)
	assert(uuid.String(), Equals, str)

	_, err = ParseBytes([]byte{1, 3, 2, 4})
	assert(err, IsNotNil)
}

func TestParseString(t *testing.T) {
	assert := With(t)

	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	expectedBytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseString(str)
	assert(err, IsNil)
	assert(uuid.Bytes(), Equals, expectedBytes)

	uuid, err = ParseString("2418d087")
	assert(err, IsNotNil)

	uuid, err = ParseString("2418d087-648k-4990-86e8-19dca1d006d3")
	assert(err, IsNotNil)
}

func TestNewUUID(t *testing.T) {
	assert := With(t)

	uuid := New()
	uuid2, err := ParseString(uuid.String())

	assert(err, IsNil)
	assert(uuid.String(), Equals, uuid2.String())
	assert(uuid.Bytes(), Equals, uuid2.Bytes())
}

func TestRandom(t *testing.T) {
	assert := With(t)

	uuid := New()
	uuid2 := New()

	assert(uuid.String(), NotEquals, uuid2.String())
	assert(uuid.Bytes(), NotEquals, uuid2.Bytes())
}

func TestEquals(t *testing.T) {
	assert := With(t)

	var uuid *UUID = nil
	var uuid2 *UUID = nil
	assert(uuid.Equals(uuid2), IsTrue)
	assert(uuid.Equals(New()), IsFalse)
}

func TestNext(t *testing.T) {
	assert := With(t)

	uuid := New()
	uuid2 := uuid.Next()
	assert(uuid.Equals(uuid2), IsFalse)
}
