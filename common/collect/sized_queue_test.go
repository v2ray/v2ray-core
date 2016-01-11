package collect_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/collect"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSizedQueue(t *testing.T) {
	v2testing.Current(t)

	queue := collect.NewSizedQueue(2)
	assert.Pointer(queue.Put(1)).IsNil()
	assert.Pointer(queue.Put(2)).IsNil()
	assert.Int(queue.Put(3).(int)).Equals(1)
}
