package mtproto_test

import (
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
