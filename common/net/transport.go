package net

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

func ReadChunk(reader io.Reader, buffer *alloc.Buffer) (*alloc.Buffer, error) {
	if buffer == nil {
		buffer = alloc.NewBuffer()
	}
	if _, err := io.ReadFull(reader, buffer.Value[:2]); err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	length := serial.BytesLiteral(buffer.Value[:2]).Uint16Value()
	if _, err := io.ReadFull(reader, buffer.Value[:length]); err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	buffer.Slice(0, int(length))
	return buffer, nil
}

func ReadAuthenticatedChunk(reader io.Reader, auth crypto.Authenticator, buffer *alloc.Buffer) (*alloc.Buffer, error) {
	buffer, err := ReadChunk(reader, buffer)
	if err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	authSize := auth.AuthBytes()

	authBytes := auth.Authenticate(nil, buffer.Value[authSize:])

	if !serial.BytesLiteral(authBytes).Equals(serial.BytesLiteral(buffer.Value[:authSize])) {
		alloc.Release(buffer)
		return nil, transport.CorruptedPacket
	}
	buffer.SliceFrom(authSize)

	return buffer, nil
}

// ReaderToChan dumps all content from a given reader to a chan by constantly reading it until EOF.
func ReaderToChan(stream chan<- *alloc.Buffer, reader io.Reader) error {
	allocate := alloc.NewBuffer
	large := false
	for {
		buffer, err := ReadFrom(reader, allocate())
		if buffer.Len() > 0 {
			stream <- buffer
		} else {
			buffer.Release()
		}
		if err != nil {
			return err
		}
		if buffer.IsFull() && !large {
			allocate = alloc.NewLargeBuffer
			large = true
		} else if !buffer.IsFull() {
			allocate = alloc.NewBuffer
			large = false
		}
	}
}

// ChanToWriter dumps all content from a given chan to a writer until the chan is closed.
func ChanToWriter(writer io.Writer, stream <-chan *alloc.Buffer) error {
	for buffer := range stream {
		nBytes, err := writer.Write(buffer.Value)
		if nBytes < buffer.Len() {
			_, err = writer.Write(buffer.Value[nBytes:])
		}
		buffer.Release()
		if err != nil {
			return err
		}
	}
	return nil
}
