package serial_test

import (
	"testing"

	. "v2ray.com/core/common/serial"
)

func TestGetInstance(t *testing.T) {
	p, err := GetInstance("")
	if p != nil {
		t.Error("expected nil instance, but got ", p)
	}
	if err == nil {
		t.Error("expect non-nil error, but got nil")
	}
}

func TestConvertingNilMessage(t *testing.T) {
	x := ToTypedMessage(nil)
	if x != nil {
		t.Error("expect nil, but actually not")
	}
}
