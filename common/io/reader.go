package io // import "github.com/v2ray/v2ray-core/common/io"

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/transport"
)

// ReadFrom reads from a reader and put all content to a buffer.
// If buffer is nil, ReadFrom creates a new normal buffer.
func ReadFrom(reader io.Reader, buffer *alloc.Buffer) (*alloc.Buffer, error) {
	if buffer == nil {
		buffer = alloc.NewBuffer()
	}
	nBytes, err := reader.Read(buffer.Value)
	buffer.Slice(0, nBytes)
	return buffer, err
}

type Reader interface {
	Read() (*alloc.Buffer, error)
}

type AdaptiveReader struct {
	reader   io.Reader
	allocate func() *alloc.Buffer
	isLarge  bool
}

func NewAdaptiveReader(reader io.Reader) *AdaptiveReader {
	return &AdaptiveReader{
		reader:   reader,
		allocate: alloc.NewBuffer,
		isLarge:  false,
	}
}

func (this *AdaptiveReader) Read() (*alloc.Buffer, error) {
	buffer, err := ReadFrom(this.reader, this.allocate())

	if buffer.IsFull() && !this.isLarge {
		this.allocate = alloc.NewLargeBuffer
		this.isLarge = true
	} else if !buffer.IsFull() {
		this.allocate = alloc.NewBuffer
		this.isLarge = false
	}

	if err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	return buffer, nil
}

type ChunkReader struct {
	reader io.Reader
}

func NewChunkReader(reader io.Reader) *ChunkReader {
	return &ChunkReader{
		reader: reader,
	}
}

func (this *ChunkReader) Read() (*alloc.Buffer, error) {
	buffer := alloc.NewLargeBuffer()
	if _, err := io.ReadFull(this.reader, buffer.Value[:2]); err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	length := serial.BytesLiteral(buffer.Value[:2]).Uint16Value()
	if _, err := io.ReadFull(this.reader, buffer.Value[:length]); err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	buffer.Slice(0, int(length))
	return buffer, nil
}

type AuthenticationReader struct {
	reader            Reader
	authenticator     crypto.Authenticator
	authBeforePayload bool
}

func NewAuthenticationReader(reader io.Reader, auth crypto.Authenticator, authBeforePayload bool) *AuthenticationReader {
	return &AuthenticationReader{
		reader:            NewChunkReader(reader),
		authenticator:     auth,
		authBeforePayload: authBeforePayload,
	}
}

func (this *AuthenticationReader) Read() (*alloc.Buffer, error) {
	buffer, err := this.reader.Read()
	if err != nil {
		alloc.Release(buffer)
		return nil, err
	}

	authSize := this.authenticator.AuthSize()
	var authBytes, payloadBytes []byte
	if this.authBeforePayload {
		authBytes = buffer.Value[:authSize]
		payloadBytes = buffer.Value[authSize:]
	} else {
		payloadBytes = buffer.Value[:authSize]
		authBytes = buffer.Value[authSize:]
	}

	actualAuthBytes := this.authenticator.Authenticate(nil, payloadBytes)
	if !serial.BytesLiteral(authBytes).Equals(serial.BytesLiteral(actualAuthBytes)) {
		alloc.Release(buffer)
		return nil, transport.CorruptedPacket
	}
	buffer.Value = payloadBytes
	return buffer, nil
}
