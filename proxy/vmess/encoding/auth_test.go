package encoding_test

import (
	"crypto/rand"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/proxy/vmess/encoding"
	. "v2ray.com/ext/assert"
)

func TestFnvAuth(t *testing.T) {
	assert := With(t)
	fnvAuth := new(FnvAuthenticator)

	expectedText := make([]byte, 256)
	_, err := rand.Read(expectedText)
	common.Must(err)

	buffer := make([]byte, 512)
	b := fnvAuth.Seal(buffer[:0], nil, expectedText, nil)
	b, err = fnvAuth.Open(buffer[:0], nil, b, nil)
	assert(err, IsNil)
	assert(len(b), Equals, 256)
	assert(b, Equals, expectedText)
}
