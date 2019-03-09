package bitmask_test

import (
	"testing"

	. "v2ray.com/core/common/bitmask"
)

func TestBitmaskByte(t *testing.T) {
	b := Byte(0)
	b.Set(Byte(1))
	if !b.Has(1) {
		t.Fatal("expected ", b, " to contain 1, but actually not")
	}

	b.Set(Byte(2))
	if !b.Has(2) {
		t.Fatal("expected ", b, " to contain 2, but actually not")
	}
	if !b.Has(1) {
		t.Fatal("expected ", b, " to contain 1, but actually not")
	}

	b.Clear(Byte(1))
	if !b.Has(2) {
		t.Fatal("expected ", b, " to contain 2, but actually not")
	}
	if b.Has(1) {
		t.Fatal("expected ", b, " to not contain 1, but actually did")
	}

	b.Toggle(Byte(2))
	if b.Has(2) {
		t.Fatal("expected ", b, " to not contain 2, but actually did")
	}
}
