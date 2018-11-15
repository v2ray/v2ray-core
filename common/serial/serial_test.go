package serial_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

func TestUint32Serial(t *testing.T) {
	b := buf.New()
	defer b.Release()

	n, err := serial.WriteUint32(b, 10)
	common.Must(err)
	if n != 4 {
		t.Error("expect 4 bytes writtng, but actually ", n)
	}
	if diff := cmp.Diff(b.Bytes(), []byte{0, 0, 0, 10}); diff != "" {
		t.Error(diff)
	}
}
