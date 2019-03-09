// +build !confonly

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
		if int32(len(b)) <= r.Header.Size() {
			return nil
		}
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
	bb := buf.StackNew()
	defer bb.Release()

	if w.Header != nil {
		w.Header.Serialize(bb.Extend(w.Header.Size()))
	}
	if w.Security != nil {
		nonceSize := w.Security.NonceSize()
		common.Must2(bb.ReadFullFrom(rand.Reader, int32(nonceSize)))
		nonce := bb.BytesFrom(int32(-nonceSize))

		encrypted := bb.Extend(int32(w.Security.Overhead() + len(b)))
		w.Security.Seal(encrypted[:0], nonce, b, nil)
	} else {
		bb.Write(b)
	}

	_, err := w.Writer.Write(bb.Bytes())
	return len(b), err
}
