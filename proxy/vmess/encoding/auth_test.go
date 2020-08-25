package encoding_test

import (
	"crypto/rand"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/proxy/vmess/encoding"
)

func TestFnvAuth(t *testing.T) {
	fnvAuth := new(FnvAuthenticator)

	expectedText := make([]byte, 256)
	_, err := rand.Read(expectedText)
	common.Must(err)

	buffer := make([]byte, 512)
	b := fnvAuth.Seal(buffer[:0], nil, expectedText, nil)
	b, err = fnvAuth.Open(buffer[:0], nil, b, nil)
	common.Must(err)
	if r := cmp.Diff(b, expectedText); r != "" {
		t.Error(r)
	}
}
