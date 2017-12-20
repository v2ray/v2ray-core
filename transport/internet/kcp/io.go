package kcp

import (
	"crypto/cipher"
	"crypto/rand"
	"io"

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

	buffer [2048]byte
}

func (w *KCPPacketWriter) Overhead() int {
	overhead := 0
	if w.Header != nil {
		overhead += w.Header.Size()
	}
	if w.Security != nil {
		overhead += w.Security.Overhead()
	}
	return overhead
}

func (w *KCPPacketWriter) Write(b []byte) (int, error) {
	x := w.buffer[:]
	size := 0
	if w.Header != nil {
		nBytes, _ := w.Header.Write(x)
		size += nBytes
		x = x[nBytes:]
	}
	if w.Security != nil {
		nonceSize := w.Security.NonceSize()
		var nonce []byte
		if nonceSize > 0 {
			nonce = x[:nonceSize]
			rand.Read(nonce)
			x = x[nonceSize:]
		}
		x = w.Security.Seal(x[:0], nonce, b, nil)
		size += nonceSize + len(x)
	} else {
		size += copy(x, b)
	}

	_, err := w.Writer.Write(w.buffer[:size])
	return len(b), err
}
