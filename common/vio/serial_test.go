package vio_test

import (
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/compare"
	"v2ray.com/core/common/vio"
)

func TestUint32Serial(t *testing.T) {
	b := buf.New()
	defer b.Release()

	n, err := vio.WriteUint32(b, 10)
	common.Must(err)
	if n != 4 {
		t.Error("expect 4 bytes writtng, but actually ", n)
	}
	if err := compare.BytesEqualWithDetail(b.Bytes(), []byte{0, 0, 0, 10}); err != nil {
		t.Error(err)
	}
}
