package utp_test

import (
	"context"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/internet/headers/utp"
	. "v2ray.com/ext/assert"
)

func TestUTPWrite(t *testing.T) {
	assert := With(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	utpRaw, err := New(context.Background(), &Config{})
	assert(err, IsNil)

	utp := utpRaw.(*UTP)

	payload := buf.NewLocal(2048)
	payload.AppendSupplier(utp.Write)
	payload.Append(content)

	assert(payload.Len(), Equals, len(content)+utp.Size())
}
