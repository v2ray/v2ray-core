package protocol_test

import (
	"testing"
	"time"

	. "github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestGenerateRandomInt64InRange(t *testing.T) {
	v2testing.Current(t)

	base := time.Now().Unix()
	delta := 100
	generator := NewRandomTimestampGenerator(Timestamp(base), delta)

	for i := 0; i < 100; i++ {
		v := int64(generator.Next())
		assert.Int64(v).AtMost(base + int64(delta))
		assert.Int64(v).AtLeast(base - int64(delta))
	}
}
