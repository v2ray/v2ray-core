package shadowsocks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/serial"
)

const (
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

func (v *Authenticator) Authenticate(data []byte) buf.BytesWriter {
	hasher := hmac.New(sha1.New, v.key())
	hasher.Write(data)
	res := hasher.Sum(nil)
	return func(b []byte) int {
		copy(b, res[:AuthSize])
		return AuthSize
	}
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
	chunkId := 0
	return func() []byte {
		newKey := make([]byte, 0, len(iv)+4)
		newKey = append(newKey, iv...)
		newKey = serial.IntToBytes(chunkId, newKey)
		chunkId++
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

func (v *ChunkReader) Release() {
	v.reader = nil
	v.auth = nil
}

func (v *ChunkReader) Read() (*buf.Buffer, error) {
	buffer := buf.NewBuffer()
	if _, err := buffer.FillFullFrom(v.reader, 2); err != nil {
		buffer.Release()
		return nil, err
	}
	// There is a potential buffer overflow here. Large buffer is 64K bytes,
	// while uin16 + 10 will be more than that
	length := serial.BytesToUint16(buffer.BytesTo(2)) + AuthSize
	if length > buf.BufferSize {
		// Theoretically the size of a chunk is 64K, but most Shadowsocks implementations used <4K buffer.
		buffer.Release()
		buffer = buf.NewLocalBuffer(int(length) + 128)
	}

	buffer.Clear()
	if _, err := buffer.FillFullFrom(v.reader, int(length)); err != nil {
		buffer.Release()
		return nil, err
	}

	authBytes := buffer.BytesTo(AuthSize)
	payload := buffer.BytesFrom(AuthSize)

	actualAuthBytes := make([]byte, AuthSize)
	v.auth.Authenticate(payload)(actualAuthBytes)
	if !bytes.Equal(authBytes, actualAuthBytes) {
		buffer.Release()
		return nil, errors.New("Shadowsocks|AuthenticationReader: Invalid auth.")
	}
	buffer.SliceFrom(AuthSize)

	return buffer, nil
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

func (v *ChunkWriter) Release() {
	v.writer = nil
	v.auth = nil
}

func (v *ChunkWriter) Write(payload *buf.Buffer) error {
	totalLength := payload.Len()
	serial.Uint16ToBytes(uint16(totalLength), v.buffer[:0])
	v.auth.Authenticate(payload.Bytes())(v.buffer[2:])
	copy(v.buffer[2+AuthSize:], payload.Bytes())
	_, err := v.writer.Write(v.buffer[:2+AuthSize+payload.Len()])
	return err
}
