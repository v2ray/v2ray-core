package srtp_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
	. "v2ray.com/core/transport/internet/headers/srtp"
)

func TestSRTPWrite(t *testing.T) {
	assert := With(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	srtp := SRTP{}

	payload := buf.NewLocal(2048)
	payload.AppendSupplier(srtp.Write)
	payload.Append(content)

	assert(payload.Len(), Equals, len(content) + srtp.Size())
}
