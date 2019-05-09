package crypto_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/crypto"
)

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	common.Must(err)
	return b
}

func TestChaCha20Stream(t *testing.T) {
	var cases = []struct {
		key    []byte
		iv     []byte
		output []byte
	}{
		{
			key: mustDecodeHex("0000000000000000000000000000000000000000000000000000000000000000"),
			iv:  mustDecodeHex("0000000000000000"),
			output: mustDecodeHex("76b8e0ada0f13d90405d6ae55386bd28bdd219b8a08ded1aa836efcc8b770dc7" +
				"da41597c5157488d7724e03fb8d84a376a43b8f41518a11cc387b669b2ee6586" +
				"9f07e7be5551387a98ba977c732d080dcb0f29a048e3656912c6533e32ee7aed" +
				"29b721769ce64e43d57133b074d839d531ed1f28510afb45ace10a1f4b794d6f"),
		},
		{
			key: mustDecodeHex("5555555555555555555555555555555555555555555555555555555555555555"),
			iv:  mustDecodeHex("5555555555555555"),
			output: mustDecodeHex("bea9411aa453c5434a5ae8c92862f564396855a9ea6e22d6d3b50ae1b3663311" +
				"a4a3606c671d605ce16c3aece8e61ea145c59775017bee2fa6f88afc758069f7" +
				"e0b8f676e644216f4d2a3422d7fa36c6c4931aca950e9da42788e6d0b6d1cd83" +
				"8ef652e97b145b14871eae6c6804c7004db5ac2fce4c68c726d004b10fcaba86"),
		},
		{
			key:    mustDecodeHex("0000000000000000000000000000000000000000000000000000000000000000"),
			iv:     mustDecodeHex("000000000000000000000000"),
			output: mustDecodeHex("76b8e0ada0f13d90405d6ae55386bd28bdd219b8a08ded1aa836efcc8b770dc7da41597c5157488d7724e03fb8d84a376a43b8f41518a11cc387b669b2ee6586"),
		},
	}
	for _, c := range cases {
		s := NewChaCha20Stream(c.key, c.iv)
		input := make([]byte, len(c.output))
		actualOutout := make([]byte, len(c.output))
		s.XORKeyStream(actualOutout, input)
		if r := cmp.Diff(c.output, actualOutout); r != "" {
			t.Fatal(r)
		}
	}
}

func TestChaCha20Decoding(t *testing.T) {
	key := make([]byte, 32)
	common.Must2(rand.Read(key))
	iv := make([]byte, 8)
	common.Must2(rand.Read(iv))
	stream := NewChaCha20Stream(key, iv)

	payload := make([]byte, 1024)
	common.Must2(rand.Read(payload))

	x := make([]byte, len(payload))
	stream.XORKeyStream(x, payload)

	stream2 := NewChaCha20Stream(key, iv)
	stream2.XORKeyStream(x, x)
	if r := cmp.Diff(x, payload); r != "" {
		t.Fatal(r)
	}
}
