package utp_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/headers/utp"
)

func TestUTPWrite(t *testing.T) {
	assert := assert.On(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	utp := UTP{}

	payload := buf.NewLocal(2048)
	payload.AppendSupplier(utp.Write)
	payload.Append(content)

	assert.Int(payload.Len()).Equals(len(content) + utp.Size())
}
