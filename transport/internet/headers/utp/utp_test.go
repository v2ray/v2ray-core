package utp_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
	. "v2ray.com/core/transport/internet/headers/utp"
)

func TestUTPWrite(t *testing.T) {
	assert := With(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	utp := UTP{}

	payload := buf.NewLocal(2048)
	payload.AppendSupplier(utp.Write)
	payload.Append(content)

	assert(payload.Len(), Equals, len(content) + utp.Size())
}
