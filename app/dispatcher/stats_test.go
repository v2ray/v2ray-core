package dispatcher_test

import (
	"testing"

	. "v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

type TestCounter int64

func (c *TestCounter) Value() int64 {
	return int64(*c)
}

func (c *TestCounter) Add(v int64) int64 {
	x := int64(*c) + v
	*c = TestCounter(x)
	return x
}

func (c *TestCounter) Set(v int64) int64 {
	*c = TestCounter(v)
	return v
}

func TestStatsWriter(t *testing.T) {
	var c TestCounter
	writer := &SizeStatWriter{
		Counter: &c,
		Writer:  buf.Discard,
	}

	mb := buf.MergeBytes(nil, []byte("abcd"))
	common.Must(writer.WriteMultiBuffer(mb))

	mb = buf.MergeBytes(nil, []byte("efg"))
	common.Must(writer.WriteMultiBuffer(mb))

	if c.Value() != 7 {
		t.Fatal("unexpected counter value. want 7, but got ", c.Value())
	}
}
