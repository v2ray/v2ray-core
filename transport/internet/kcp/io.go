package kcp

import (
	"crypto/cipher"
	"crypto/rand"
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/transport/internet"
)

type PacketReader interface {
	Read([]byte) []Segment
}

type PacketWriter interface {
	Overhead() int
	io.Writer
}

type KCPPacketReader struct {
	Security cipher.AEAD
	Header   internet.PacketHeader
}

func (r *KCPPacketReader) Read(b []byte) []Segment {
	if r.Header != nil {
		b = b[r.Header.Size():]
	}
	if r.Security != nil {
		nonceSize := r.Security.NonceSize()
		overhead := r.Security.Overhead()
		if len(b) <= nonceSize+overhead {
			return nil
		}
		out, err := r.Security.Open(b[nonceSize:nonceSize], b[:nonceSize], b[nonceSize:], nil)
		if err != nil {
			return nil
		}
		b = out
	}
	var result []Segment
	for len(b) > 0 {
		seg, x := ReadSegment(b)
		if seg == nil {
			break
		}
		result = append(result, seg)
		b = x
	}
	return result
}

type KCPPacketWriter struct {
	Header   internet.PacketHeader
	Security cipher.AEAD
	Writer   io.Writer
}

func (w *KCPPacketWriter) Overhead() int {
	overhead := 0
	if w.Header != nil {
		overhead += int(w.Header.Size())
	}
	if w.Security != nil {
		overhead += w.Security.Overhead()
	}
	return overhead
}

func (w *KCPPacketWriter) Write(b []byte) (int, error) {
	bb := buf.New()
	defer bb.Release()

	if w.Header != nil {
		common.Must(bb.AppendSupplier(func(x []byte) (int, error) {
			return w.Header.Write(x)
		}))
	}
	if w.Security != nil {
		nonceSize := w.Security.NonceSize()
		common.Must(bb.AppendSupplier(func(x []byte) (int, error) {
			return rand.Read(x[:nonceSize])
		}))
		nonce := bb.BytesFrom(int32(-nonceSize))
		common.Must(bb.AppendSupplier(func(x []byte) (int, error) {
			eb := w.Security.Seal(x[:0], nonce, b, nil)
			return len(eb), nil
		}))
	} else {
		bb.Write(b)
	}

	_, err := w.Writer.Write(bb.Bytes())
	return len(b), err
}
