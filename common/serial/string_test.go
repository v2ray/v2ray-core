package serial_test

import (
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
)

func TestUint32(t *testing.T) {
	assert := assert.On(t)

	x := uint32(458634234)
	s1 := Uint32ToBytes(x, []byte{})
	s2 := buf.New()
	common.Must(s2.AppendSupplier(WriteUint32(x)))
	assert.Bytes(s1).Equals(s2.Bytes())
}
