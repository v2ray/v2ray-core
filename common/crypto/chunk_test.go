package crypto_test

import (
	"bytes"
	"io"
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/crypto"
)

func TestChunkStreamIO(t *testing.T) {
	cache := bytes.NewBuffer(make([]byte, 0, 8192))

	writer := NewChunkStreamWriter(PlainChunkSizeParser{}, cache)
	reader := NewChunkStreamReader(PlainChunkSizeParser{}, cache)

	b := buf.New()
	b.WriteString("abcd")
	common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{b}))

	b = buf.New()
	b.WriteString("efg")
	common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{b}))

	common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{}))

	if cache.Len() != 13 {
		t.Fatalf("Cache length is %d, want 13", cache.Len())
	}

	mb, err := reader.ReadMultiBuffer()
	common.Must(err)

	if s := mb.String(); s != "abcd" {
		t.Error("content: ", s)
	}

	mb, err = reader.ReadMultiBuffer()
	common.Must(err)

	if s := mb.String(); s != "efg" {
		t.Error("content: ", s)
	}

	_, err = reader.ReadMultiBuffer()
	if err != io.EOF {
		t.Error("error: ", err)
	}
}
