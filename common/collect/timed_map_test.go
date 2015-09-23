package collect

import (
	"testing"
	"time"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestTimedStringMap(t *testing.T) {
	assert := unit.Assert(t)

	nowSec := time.Now().UTC().Unix()
	m := NewTimedStringMap(2)
	m.Set("Key1", "Value1", nowSec)
	m.Set("Key2", "Value2", nowSec+5)

	v1, ok := m.Get("Key1")
	assert.Bool(ok).IsTrue()
	assert.String(v1.(string)).Equals("Value1")

	v2, ok := m.Get("Key2")
	assert.Bool(ok).IsTrue()
	assert.String(v2.(string)).Equals("Value2")

	tick := time.Tick(4 * time.Second)
	<-tick

	v1, ok = m.Get("Key1")
	assert.Bool(ok).IsFalse()

	v2, ok = m.Get("Key2")
	assert.Bool(ok).IsTrue()
	assert.String(v2.(string)).Equals("Value2")

	<-tick
	v2, ok = m.Get("Key2")
	assert.Bool(ok).IsFalse()

	<-tick
	v2, ok = m.Get("Key2")
	assert.Bool(ok).IsFalse()

	m.Set("Key1", "Value1", time.Now().UTC().Unix()+10)
	v1, ok = m.Get("Key1")
	assert.Bool(ok).IsTrue()
	assert.String(v1.(string)).Equals("Value1")
}
