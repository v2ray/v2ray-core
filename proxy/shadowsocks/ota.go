package shadowsocks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
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

func (v *Authenticator) Authenticate(data []byte) buf.Supplier {
	hasher := hmac.New(sha1.New, v.key())
	common.Must2(hasher.Write(data))
	res := hasher.Sum(nil)
	return func(b []byte) (int, error) {
		return copy(b, res[:AuthSize]), nil
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
	chunkID := 0
	return func() []byte {
		newKey := make([]byte, 0, len(iv)+4)
		newKey = append(newKey, iv...)
		newKey = serial.IntToBytes(chunkID, newKey)
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
		return nil, newError("failed to read size")
	}
	size += AuthSize

	buffer := buf.NewSize(int32(size))
	if err := buffer.AppendSupplier(buf.ReadFullFrom(v.reader, int32(size))); err != nil {
		buffer.Release()
		return nil, err
	}

	authBytes := buffer.BytesTo(AuthSize)
	payload := buffer.BytesFrom(AuthSize)

	actualAuthBytes := make([]byte, AuthSize)
	v.auth.Authenticate(payload)(actualAuthBytes)
	if !bytes.Equal(authBytes, actualAuthBytes) {
		buffer.Release()
		return nil, newError("invalid auth")
	}
	buffer.Advance(AuthSize)

	return buf.NewMultiBufferValue(buffer), nil
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
	defer mb.Release()

	for {
		payloadLen, _ := mb.Read(w.buffer[2+AuthSize:])
		serial.Uint16ToBytes(uint16(payloadLen), w.buffer[:0])
		w.auth.Authenticate(w.buffer[2+AuthSize : 2+AuthSize+payloadLen])(w.buffer[2:])
		if err := buf.WriteAllBytes(w.writer, w.buffer[:2+AuthSize+payloadLen]); err != nil {
			return err
		}
		if mb.IsEmpty() {
			break
		}
	}
	return nil
}
