package internet_test

import (
	"testing"

	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/headers/noop"
	"v2ray.com/core/transport/internet/headers/srtp"
	"v2ray.com/core/transport/internet/headers/utp"
)

func TestAllHeadersLoadable(t *testing.T) {
	assert := assert.On(t)

	noopAuth, err := CreatePacketHeader((*noop.Config)(nil))
	assert.Error(err).IsNil()
	assert.Int(noopAuth.Size()).Equals(0)

	srtp, err := CreatePacketHeader((*srtp.Config)(nil))
	assert.Error(err).IsNil()
	assert.Int(srtp.Size()).Equals(4)

	utp, err := CreatePacketHeader((*utp.Config)(nil))
	assert.Error(err).IsNil()
	assert.Int(utp.Size()).Equals(4)
}
