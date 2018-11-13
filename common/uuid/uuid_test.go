package uuid_test

import (
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/compare"
	. "v2ray.com/core/common/uuid"
	. "v2ray.com/ext/assert"
)

func TestParseBytes(t *testing.T) {
	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	bytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseBytes(bytes)
	common.Must(err)
	if err := compare.StringEqualWithDetail(uuid.String(), str); err != nil {
		t.Fatal(err)
	}

	_, err = ParseBytes([]byte{1, 3, 2, 4})
	if err == nil {
		t.Fatal("Expect error but nil")
	}
}

func TestParseString(t *testing.T) {
	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	expectedBytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseString(str)
	common.Must(err)
	if err := compare.BytesEqualWithDetail(expectedBytes, uuid.Bytes()); err != nil {
		t.Fatal(err)
	}

	_, err = ParseString("2418d087")
	if err == nil {
		t.Fatal("Expect error but nil")
	}

	_, err = ParseString("2418d087-648k-4990-86e8-19dca1d006d3")
	if err == nil {
		t.Fatal("Expect error but nil")
	}
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

	var uuid *UUID
	var uuid2 *UUID
	assert(uuid.Equals(uuid2), IsTrue)

	uuid3 := New()
	assert(uuid.Equals(&uuid3), IsFalse)
}
