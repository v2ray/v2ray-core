package buf_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
)

func TestMultiBufferRead(t *testing.T) {
	assert := With(t)

	b1 := New()
	b1.WriteString("ab")

	b2 := New()
	b2.WriteString("cd")
	mb := MultiBuffer{b1, b2}

	bs := make([]byte, 32)
	_, nBytes := SplitBytes(mb, bs)
	assert(nBytes, Equals, 4)
	assert(bs[:nBytes], Equals, []byte("abcd"))
}

func TestMultiBufferAppend(t *testing.T) {
	assert := With(t)

	var mb MultiBuffer
	b := New()
	b.WriteString("ab")
	mb = append(mb, b)
	assert(mb.Len(), Equals, int32(2))
}

func TestMultiBufferSliceBySizeLarge(t *testing.T) {
	lb := make([]byte, 8*1024)
	common.Must2(io.ReadFull(rand.Reader, lb))

	mb := MergeBytes(nil, lb)

	mb, mb2 := SplitSize(mb, 1024)
	if mb2.Len() != 1024 {
		t.Error("expect length 1024, but got ", mb2.Len())
	}
	if mb.Len() != 7*1024 {
		t.Error("expect length 7*1024, but got ", mb.Len())
	}

	mb, mb3 := SplitSize(mb, 7*1024)
	if mb3.Len() != 7*1024 {
		t.Error("expect length 7*1024, but got", mb.Len())
	}

	if !mb.IsEmpty() {
		t.Error("expect empty buffer, but got ", mb.Len())
	}
}

func TestMultiBufferSplitFirst(t *testing.T) {
	b1 := New()
	b1.WriteString("b1")

	b2 := New()
	b2.WriteString("b2")

	b3 := New()
	b3.WriteString("b3")

	var mb MultiBuffer
	mb = append(mb, b1, b2, b3)

	mb, c1 := SplitFirst(mb)
	if diff := cmp.Diff(b1.String(), c1.String()); diff != "" {
		t.Error(diff)
	}

	mb, c2 := SplitFirst(mb)
	if diff := cmp.Diff(b2.String(), c2.String()); diff != "" {
		t.Error(diff)
	}

	mb, c3 := SplitFirst(mb)
	if diff := cmp.Diff(b3.String(), c3.String()); diff != "" {
		t.Error(diff)
	}

	if !mb.IsEmpty() {
		t.Error("expect empty buffer, but got ", mb.String())
	}
}
