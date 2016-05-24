package protocol_test

import (
	"testing"
	"time"

	. "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestGenerateRandomInt64InRange(t *testing.T) {
	assert := assert.On(t)

	base := time.Now().Unix()
	delta := 100
	generator := NewTimestampGenerator(Timestamp(base), delta)

	for i := 0; i < 100; i++ {
		v := int64(generator())
		assert.Int64(v).AtMost(base + int64(delta))
		assert.Int64(v).AtLeast(base - int64(delta))
	}
}
