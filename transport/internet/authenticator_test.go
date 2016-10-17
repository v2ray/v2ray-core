package internet_test

import (
	"testing"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/authenticators/noop"
	"v2ray.com/core/transport/internet/authenticators/srtp"
	"v2ray.com/core/transport/internet/authenticators/utp"
)

func TestAllAuthenticatorLoadable(t *testing.T) {
	assert := assert.On(t)

	noopAuth, err := CreateAuthenticator(loader.GetType(new(noop.Config)), nil)
	assert.Error(err).IsNil()
	assert.Int(noopAuth.Overhead()).Equals(0)

	srtp, err := CreateAuthenticator(loader.GetType(new(srtp.Config)), nil)
	assert.Error(err).IsNil()
	assert.Int(srtp.Overhead()).Equals(4)

	utp, err := CreateAuthenticator(loader.GetType(new(utp.Config)), nil)
	assert.Error(err).IsNil()
	assert.Int(utp.Overhead()).Equals(4)
}
