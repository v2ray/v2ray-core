package shadowsocks_test

import (
	"crypto/rand"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/proxy/shadowsocks"
)

func TestAEADCipherUDP(t *testing.T) {
	rawAccount := &shadowsocks.Account{
		CipherType: shadowsocks.CipherType_AES_128_GCM,
		Password:   "test",
	}
	account, err := rawAccount.AsAccount()
	common.Must(err)

	cipher := account.(*shadowsocks.MemoryAccount).Cipher

	key := make([]byte, cipher.KeySize())
	common.Must2(rand.Read(key))

	payload := make([]byte, 1024)
	common.Must2(rand.Read(payload))

	b1 := buf.New()
	common.Must2(b1.ReadFullFrom(rand.Reader, cipher.IVSize()))
	common.Must2(b1.Write(payload))
	common.Must(cipher.EncodePacket(key, b1))

	common.Must(cipher.DecodePacket(key, b1))
	if diff := cmp.Diff(b1.Bytes(), payload); diff != "" {
		t.Error(diff)
	}
}
