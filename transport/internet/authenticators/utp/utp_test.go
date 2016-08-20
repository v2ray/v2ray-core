package utp_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/authenticators/utp"
)

func TestUTPOpenSeal(t *testing.T) {
	assert := assert.On(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	payload := alloc.NewLocalBuffer(2048).Clear().Append(content)
	utp := UTP{}
	utp.Seal(payload)
	assert.Int(payload.Len()).GreaterThan(len(content))
	assert.Bool(utp.Open(payload)).IsTrue()
	assert.Bytes(content).Equals(payload.Bytes())
}
