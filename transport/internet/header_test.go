package internet_test

import (
	"testing"

	. "v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/headers/noop"
	"v2ray.com/core/transport/internet/headers/srtp"
	"v2ray.com/core/transport/internet/headers/utp"
	. "v2ray.com/ext/assert"
)

func TestAllHeadersLoadable(t *testing.T) {
	assert := With(t)

	noopAuth, err := CreatePacketHeader((*noop.Config)(nil))
	assert(err, IsNil)
	assert(noopAuth.Size(), Equals, 0)

	srtp, err := CreatePacketHeader((*srtp.Config)(nil))
	assert(err, IsNil)
	assert(srtp.Size(), Equals, 4)

	utp, err := CreatePacketHeader((*utp.Config)(nil))
	assert(err, IsNil)
	assert(utp.Size(), Equals, 4)
}
