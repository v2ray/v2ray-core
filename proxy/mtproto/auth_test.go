package mtproto_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/proxy/mtproto"
	. "v2ray.com/ext/assert"
)

func TestInverse(t *testing.T) {
	const size = 64
	b := make([]byte, 64)
	for b[0] == b[size-1] {
		common.Must2(rand.Read(b))
	}

	bi := Inverse(b)
	if b[0] == bi[0] {
		t.Fatal("seems bytes are not inversed: ", b[0], "vs", bi[0])
	}

	bii := Inverse(bi)
	if r := cmp.Diff(bii, b); r != "" {
		t.Fatal(r)
	}
}

func TestAuthenticationReadWrite(t *testing.T) {
	assert := With(t)

	a := NewAuthentication(DefaultSessionContext())
	b := bytes.NewReader(a.Header[:])
	a2, err := ReadAuthentication(b)
	assert(err, IsNil)

	assert(a.EncodingKey[:], Equals, a2.DecodingKey[:])
	assert(a.EncodingNonce[:], Equals, a2.DecodingNonce[:])
	assert(a.DecodingKey[:], Equals, a2.EncodingKey[:])
	assert(a.DecodingNonce[:], Equals, a2.EncodingNonce[:])
}
