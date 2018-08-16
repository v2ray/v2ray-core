package mtproto_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/proxy/mtproto"
	. "v2ray.com/ext/assert"
)

func TestInverse(t *testing.T) {
	assert := With(t)

	b := make([]byte, 64)
	common.Must2(rand.Read(b))

	bi := Inverse(b)
	assert(b[0], NotEquals, bi[0])

	bii := Inverse(bi)
	assert(bii, Equals, b)
}

func TestAuthenticationReadWrite(t *testing.T) {
	assert := With(t)

	a := NewAuthentication()
	b := bytes.NewReader(a.Header[:])
	a2, err := ReadAuthentication(b)
	assert(err, IsNil)

	assert(a.EncodingKey[:], Equals, a2.DecodingKey[:])
	assert(a.EncodingNonce[:], Equals, a2.DecodingNonce[:])
	assert(a.DecodingKey[:], Equals, a2.EncodingKey[:])
	assert(a.DecodingNonce[:], Equals, a2.EncodingNonce[:])
}
