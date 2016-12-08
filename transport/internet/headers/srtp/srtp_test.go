package srtp_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/headers/srtp"
)

func TestSRTPWrite(t *testing.T) {
	assert := assert.On(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	srtp := SRTP{}

	payload := alloc.NewLocalBuffer(2048)
	payload.AppendFunc(srtp.Write)
	payload.Append(content)

	assert.Int(payload.Len()).Equals(len(content) + srtp.Size())
}
