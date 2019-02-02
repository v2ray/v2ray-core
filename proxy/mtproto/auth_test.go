package mtproto_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/proxy/mtproto"
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
	a := NewAuthentication(DefaultSessionContext())
	b := bytes.NewReader(a.Header[:])
	a2, err := ReadAuthentication(b)
	common.Must(err)

	if r := cmp.Diff(a.EncodingKey[:], a2.DecodingKey[:]); r != "" {
		t.Error("decoding key: ", r)
	}

	if r := cmp.Diff(a.EncodingNonce[:], a2.DecodingNonce[:]); r != "" {
		t.Error("decoding nonce: ", r)
	}

	if r := cmp.Diff(a.DecodingKey[:], a2.EncodingKey[:]); r != "" {
		t.Error("encoding key: ", r)
	}

	if r := cmp.Diff(a.DecodingNonce[:], a2.EncodingNonce[:]); r != "" {
		t.Error("encoding nonce: ", r)
	}
}
