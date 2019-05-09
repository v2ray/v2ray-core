package uuid_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/uuid"
)

func TestParseBytes(t *testing.T) {
	str := "2418d087-648d-4990-86e8-19dca1d006d3"
	bytes := []byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseBytes(bytes)
	common.Must(err)
	if diff := cmp.Diff(uuid.String(), str); diff != "" {
		t.Error(diff)
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
	if r := cmp.Diff(expectedBytes, uuid.Bytes()); r != "" {
		t.Fatal(r)
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
	uuid := New()
	uuid2, err := ParseString(uuid.String())

	common.Must(err)
	if uuid.String() != uuid2.String() {
		t.Error("uuid string: ", uuid.String(), " != ", uuid2.String())
	}
	if r := cmp.Diff(uuid.Bytes(), uuid2.Bytes()); r != "" {
		t.Error(r)
	}
}

func TestRandom(t *testing.T) {
	uuid := New()
	uuid2 := New()

	if uuid.String() == uuid2.String() {
		t.Error("duplicated uuid")
	}
}

func TestEquals(t *testing.T) {
	var uuid *UUID
	var uuid2 *UUID
	if !uuid.Equals(uuid2) {
		t.Error("empty uuid should equal")
	}

	uuid3 := New()
	if uuid.Equals(&uuid3) {
		t.Error("nil uuid equals non-nil uuid")
	}
}
