package srtp_test

import (
	"context"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/internet/headers/srtp"
	. "v2ray.com/ext/assert"
)

func TestSRTPWrite(t *testing.T) {
	assert := With(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	srtpRaw, err := New(context.Background(), &Config{})
	assert(err, IsNil)

	srtp := srtpRaw.(*SRTP)

	payload := buf.New()
	payload.AppendSupplier(srtp.Write)
	payload.Write(content)

	assert(payload.Len(), Equals, int32(len(content))+srtp.Size())
}
