package user

import (
	"testing"
	"time"

  v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestGenerateRandomInt64InRange(t *testing.T) {
	v2testing.Current(t)
	base := time.Now().Unix()
	delta := 100

	for i := 0; i < 100; i++ {
		v := GenerateRandomInt64InRange(base, delta)
		assert.Int64(v).AtMost(base + int64(delta))
		assert.Int64(v).AtLeast(base - int64(delta))
	}
}
