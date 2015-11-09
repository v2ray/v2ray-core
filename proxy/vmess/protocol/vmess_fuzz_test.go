package protocol

import (
	"crypto/rand"
	"testing"

	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user/testing/mocks"
)

func TestVMessRequestReader(t *testing.T) {
	reader := NewVMessRequestReader(&mocks.StaticUserSet{})
	for i := 0; i < 10000000; i++ {
		reader.Read(rand.Reader)
	}
}
