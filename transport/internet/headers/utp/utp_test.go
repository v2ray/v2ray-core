package utp_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/headers/utp"
)

func TestUTPWrite(t *testing.T) {
	assert := assert.On(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	utp := UTP{}

	payload := alloc.NewLocalBuffer(2048)
	payload.AppendFunc(utp.Write)
	payload.Append(content)

	assert.Int(payload.Len()).Equals(len(content) + utp.Size())
}
