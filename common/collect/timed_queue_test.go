package collect

import (
	"testing"
	"time"

	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestTimedQueue(t *testing.T) {
	v2testing.Current(t)

	removed := make(map[string]bool)

	nowSec := time.Now().Unix()
	q := NewTimedQueue(2, func(v interface{}) {
		removed[v.(string)] = true
	})

	q.Add("Value1", nowSec)
	q.Add("Value2", nowSec+5)

	v1, ok := removed["Value1"]
	assert.Bool(ok).IsFalse()

	v2, ok := removed["Value2"]
	assert.Bool(ok).IsFalse()

	tick := time.Tick(4 * time.Second)
	<-tick

	v1, ok = removed["Value1"]
	assert.Bool(ok).IsTrue()
	assert.Bool(v1).IsTrue()
	removed["Value1"] = false

	v2, ok = removed["Value2"]
	assert.Bool(ok).IsFalse()

	<-tick
	v2, ok = removed["Value2"]
	assert.Bool(ok).IsTrue()
	assert.Bool(v2).IsTrue()
	removed["Value2"] = false

	<-tick
	assert.Bool(removed["Values"]).IsFalse()

	q.Add("Value1", time.Now().Unix()+10)

	<-tick
	v1, ok = removed["Value1"]
	assert.Bool(ok).IsTrue()
	assert.Bool(v1).IsFalse()
}
