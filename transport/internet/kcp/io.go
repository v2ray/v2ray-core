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

func (v *KCPPacketReader) Read(b []byte) []Segment {
	if v.Header != nil {
		b = b[v.Header.Size():]
	}
	if v.Security != nil {
		nonceSize := v.Security.NonceSize()
		out, err := v.Security.Open(b[nonceSize:nonceSize], b[:nonceSize], b[nonceSize:], nil)
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

	buffer [32 * 1024]byte
}

func (v *KCPPacketWriter) Overhead() int {
	overhead := 0
	if v.Header != nil {
		overhead += v.Header.Size()
	}
	if v.Security != nil {
		overhead += v.Security.Overhead()
	}
	return overhead
}

func (v *KCPPacketWriter) Write(b []byte) (int, error) {
	x := v.buffer[:]
	size := 0
	if v.Header != nil {
		nBytes := v.Header.Write(x)
		size += nBytes
		x = x[nBytes:]
	}
	if v.Security != nil {
		nonceSize := v.Security.NonceSize()
		var nonce []byte
		if nonceSize > 0 {
			nonce = x[:nonceSize]
			rand.Read(nonce)
			x = x[nonceSize:]
		}
		x = v.Security.Seal(x[:0], nonce, b, nil)
		size += nonceSize + len(x)
	} else {
		size += copy(x, b)
	}

	_, err := v.Writer.Write(v.buffer[:size])
	return len(b), err
}
