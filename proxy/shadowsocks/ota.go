package shadowsocks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"io"
	"v2ray.com/core/common/alloc"
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

func (v *Authenticator) Authenticate(auth []byte, data []byte) []byte {
	hasher := hmac.New(sha1.New, v.key())
	hasher.Write(data)
	res := hasher.Sum(nil)
	return append(auth, res[:AuthSize]...)
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

func (v *ChunkReader) Read() (*alloc.Buffer, error) {
	buffer := alloc.NewBuffer()
	if _, err := io.ReadFull(v.reader, buffer.Value[:2]); err != nil {
		buffer.Release()
		return nil, err
	}
	// There is a potential buffer overflow here. Large buffer is 64K bytes,
	// while uin16 + 10 will be more than that
	length := serial.BytesToUint16(buffer.Value[:2]) + AuthSize
	if length > alloc.BufferSize {
		// Theoretically the size of a chunk is 64K, but most Shadowsocks implementations used <4K buffer.
		buffer.Release()
		buffer = alloc.NewLocalBuffer(int(length) + 128)
	}
	if _, err := io.ReadFull(v.reader, buffer.Value[:length]); err != nil {
		buffer.Release()
		return nil, err
	}
	buffer.Slice(0, int(length))

	authBytes := buffer.Value[:AuthSize]
	payload := buffer.Value[AuthSize:]

	actualAuthBytes := v.auth.Authenticate(nil, payload)
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
}

func NewChunkWriter(writer io.Writer, auth *Authenticator) *ChunkWriter {
	return &ChunkWriter{
		writer: writer,
		auth:   auth,
	}
}

func (v *ChunkWriter) Release() {
	v.writer = nil
	v.auth = nil
}

func (v *ChunkWriter) Write(payload *alloc.Buffer) error {
	totalLength := payload.Len()
	payload.SliceBack(AuthSize)
	v.auth.Authenticate(payload.Value[:0], payload.Value[AuthSize:])
	payload.PrependUint16(uint16(totalLength))
	_, err := v.writer.Write(payload.Bytes())
	return err
}
