package protocol

import (
	"testing"

	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/fuzzing"
)

func TestVMessRequestReader(t *testing.T) {
	reader := NewVMessRequestReader(&mocks.StaticUserSet{})
	for i := 0; i < 1000000; i++ {
		reader.Read(fuzzing.RandomReader())
	}
}
