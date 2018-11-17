package shadowsocks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bytespool"
	"v2ray.com/core/common/serial"
)

const (
	// AuthSize is the number of extra bytes for Shadowsocks OTA.
	AuthSize = 10
)

type KeyGenerator func() []byte

type Authenticator struct {
	key KeyGenerator
}

func NewAuthenticator(keygen KeyGenerator) *Authenticator {
	return &Authenticator{
		key: keygen,
	}
}

func (v *Authenticator) Authenticate(data []byte, dest []byte) {
	hasher := hmac.New(sha1.New, v.key())
	common.Must2(hasher.Write(data))
	res := hasher.Sum(nil)
	copy(dest, res[:AuthSize])
}

func HeaderKeyGenerator(key []byte, iv []byte) func() []byte {
	return func() []byte {
		newKey := make([]byte, 0, len(key)+len(iv))
		newKey = append(newKey, iv...)
		newKey = append(newKey, key...)
		return newKey
	}
}

func ChunkKeyGenerator(iv []byte) func() []byte {
	chunkID := uint32(0)
	return func() []byte {
		newKey := make([]byte, len(iv)+4)
		copy(newKey, iv)
		binary.BigEndian.PutUint32(newKey[len(iv):], chunkID)
		chunkID++
		return newKey
	}
}

type ChunkReader struct {
	reader io.Reader
	auth   *Authenticator
}

func NewChunkReader(reader io.Reader, auth *Authenticator) *ChunkReader {
	return &ChunkReader{
		reader: reader,
		auth:   auth,
	}
}

func (v *ChunkReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	size, err := serial.ReadUint16(v.reader)
	if err != nil {
		return nil, newError("failed to read size").Base(err)
	}
	size += AuthSize

	buffer := bytespool.Alloc(int32(size))
	defer bytespool.Free(buffer)

	if _, err := io.ReadFull(v.reader, buffer[:size]); err != nil {
		return nil, err
	}

	authBytes := buffer[:AuthSize]
	payload := buffer[AuthSize:size]

	actualAuthBytes := make([]byte, AuthSize)
	v.auth.Authenticate(payload, actualAuthBytes)
	if !bytes.Equal(authBytes, actualAuthBytes) {
		return nil, newError("invalid auth")
	}

	var mb buf.MultiBuffer
	common.Must2(mb.Write(payload))

	return mb, nil
}

type ChunkWriter struct {
	writer io.Writer
	auth   *Authenticator
	buffer []byte
}

func NewChunkWriter(writer io.Writer, auth *Authenticator) *ChunkWriter {
	return &ChunkWriter{
		writer: writer,
		auth:   auth,
		buffer: make([]byte, 32*1024),
	}
}

// WriteMultiBuffer implements buf.Writer.
func (w *ChunkWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer buf.ReleaseMulti(mb)

	for {
		payloadLen, _ := mb.Read(w.buffer[2+AuthSize:])
		binary.BigEndian.PutUint16(w.buffer, uint16(payloadLen))
		w.auth.Authenticate(w.buffer[2+AuthSize:2+AuthSize+payloadLen], w.buffer[2:])
		if err := buf.WriteAllBytes(w.writer, w.buffer[:2+AuthSize+payloadLen]); err != nil {
			return err
		}
		if mb.IsEmpty() {
			break
		}
	}
	return nil
}
