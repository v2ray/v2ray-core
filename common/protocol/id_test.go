package protocol_test

import (
	"testing"

	. "v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
)

func TestIdEquals(t *testing.T) {
	id1 := NewID(uuid.New())
	id2 := NewID(id1.UUID())

	if !id1.Equals(id2) {
		t.Error("expected id1 to equal id2, but actually not")
	}

	if !id1.String() != id2.String() {
		t.Error(id1.String(), " != ", id2.String())
	}
}
