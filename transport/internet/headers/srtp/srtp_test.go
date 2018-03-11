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

	payload := buf.NewSize(2048)
	payload.AppendSupplier(srtp.Write)
	payload.Append(content)

	assert(payload.Len(), Equals, len(content)+srtp.Size())
}
