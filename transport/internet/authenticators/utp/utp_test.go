package utp_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/authenticators/utp"
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
