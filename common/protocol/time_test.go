package protocol_test

import (
	"testing"
	"time"

	. "v2ray.com/core/common/protocol"
)

func TestGenerateRandomInt64InRange(t *testing.T) {
	base := time.Now().Unix()
	delta := 100
	generator := NewTimestampGenerator(Timestamp(base), delta)

	for i := 0; i < 100; i++ {
		val := int64(generator())
		if val > base+int64(delta) || val < base-int64(delta) {
			t.Error(val, " not between ", base-int64(delta), " and ", base+int64(delta))
		}
	}
}
