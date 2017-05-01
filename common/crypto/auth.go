package crypto

import (
	"crypto/cipher"
	"io"

	"v2ray.com/core/common/buf"
)

type BytesGenerator interface {
	Next() []byte
}

type NoOpBytesGenerator struct {
	buffer [1]byte
}

func (v NoOpBytesGenerator) Next() []byte {
	return v.buffer[:0]
}

type StaticBytesGenerator struct {
	Content []byte
}

func (v StaticBytesGenerator) Next() []byte {
	return v.Content
}

type Authenticator interface {
	NonceSize() int
	Overhead() int
	Open(dst, cipherText []byte) ([]byte, error)
	Seal(dst, plainText []byte) ([]byte, error)
}

type AEADAuthenticator struct {
	cipher.AEAD
	NonceGenerator          BytesGenerator
	AdditionalDataGenerator BytesGenerator
}

func (v *AEADAuthenticator) Open(dst, cipherText []byte) ([]byte, error) {
	iv := v.NonceGenerator.Next()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, newError("invalid AEAD nonce size: ", len(iv))
	}

	additionalData := v.AdditionalDataGenerator.Next()
	return v.AEAD.Open(dst, iv, cipherText, additionalData)
}

func (v *AEADAuthenticator) Seal(dst, plainText []byte) ([]byte, error) {
	iv := v.NonceGenerator.Next()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, newError("invalid AEAD nonce size: ", len(iv))
	}

	additionalData := v.AdditionalDataGenerator.Next()
	return v.AEAD.Seal(dst, iv, plainText, additionalData), nil
}

type StreamMode int

const (
	ModeStream StreamMode = iota
	ModePacket
)

type AuthenticationReader struct {
	auth       Authenticator
	buffer     *buf.Buffer
	reader     io.Reader
	sizeParser ChunkSizeDecoder
	size       int
	mode       StreamMode
}

const (
	readerBufferSize = 32 * 1024
)

func NewAuthenticationReader(auth Authenticator, sizeParser ChunkSizeDecoder, reader io.Reader, mode StreamMode) *AuthenticationReader {
	return &AuthenticationReader{
		auth:       auth,
		buffer:     buf.NewLocal(readerBufferSize),
		reader:     reader,
		sizeParser: sizeParser,
		size:       -1,
		mode:       mode,
	}
}

func (r *AuthenticationReader) readSize() error {
	if r.size >= 0 {
		return nil
	}

	sizeBytes := r.sizeParser.SizeBytes()
	if r.buffer.Len() < sizeBytes {
		r.buffer.Reset(buf.ReadFrom(r.buffer))
		delta := sizeBytes - r.buffer.Len()
		if err := r.buffer.AppendSupplier(buf.ReadAtLeastFrom(r.reader, delta)); err != nil {
			return err
		}
	}
	size, err := r.sizeParser.Decode(r.buffer.BytesTo(sizeBytes))
	if err != nil {
		return err
	}
	r.size = int(size)
	r.buffer.SliceFrom(sizeBytes)
	return nil
}

func (r *AuthenticationReader) readChunk(waitForData bool) ([]byte, error) {
	if err := r.readSize(); err != nil {
		return nil, err
	}
	if r.size > readerBufferSize-r.sizeParser.SizeBytes() {
		return nil, newError("size too large ", r.size).AtWarning()
	}

	if r.size == r.auth.Overhead() {
		return nil, io.EOF
	}

	if r.buffer.Len() < r.size {
		if !waitForData {
			return nil, io.ErrNoProgress
		}
		r.buffer.Reset(buf.ReadFrom(r.buffer))

		delta := r.size - r.buffer.Len()
		if err := r.buffer.AppendSupplier(buf.ReadAtLeastFrom(r.reader, delta)); err != nil {
			return nil, err
		}
	}

	b, err := r.auth.Open(r.buffer.BytesTo(0), r.buffer.BytesTo(r.size))
	if err != nil {
		return nil, err
	}
	r.buffer.SliceFrom(r.size)
	r.size = -1
	return b, nil
}

func (r *AuthenticationReader) Read() (buf.MultiBuffer, error) {
	b, err := r.readChunk(true)
	if err != nil {
		return nil, err
	}

	mb := buf.NewMultiBuffer()
	if r.mode == ModeStream {
		mb.Write(b)
	} else {
		var bb *buf.Buffer
		if len(b) < buf.Size {
			bb = buf.New()
		} else {
			bb = buf.NewLocal(len(b))
		}
		bb.Append(b)
		mb.Append(bb)
	}

	for r.buffer.Len() >= r.sizeParser.SizeBytes() {
		b, err := r.readChunk(false)
		if err != nil {
			break
		}
		if r.mode == ModeStream {
			mb.Write(b)
		} else {
			var bb *buf.Buffer
			if len(b) < buf.Size {
				bb = buf.New()
			} else {
				bb = buf.NewLocal(len(b))
			}
			bb.Append(b)
			mb.Append(bb)
		}
	}

	return mb, nil
}

type AuthenticationWriter struct {
	auth       Authenticator
	payload    []byte
	buffer     *buf.Buffer
	writer     io.Writer
	sizeParser ChunkSizeEncoder
	mode       StreamMode
}

func NewAuthenticationWriter(auth Authenticator, sizeParser ChunkSizeEncoder, writer io.Writer, mode StreamMode) *AuthenticationWriter {
	return &AuthenticationWriter{
		auth:       auth,
		payload:    make([]byte, 1024),
		buffer:     buf.NewLocal(readerBufferSize),
		writer:     writer,
		sizeParser: sizeParser,
		mode:       mode,
	}
}

func (w *AuthenticationWriter) append(b []byte) {
	encryptedSize := len(b) + w.auth.Overhead()

	w.buffer.AppendSupplier(func(bb []byte) (int, error) {
		w.sizeParser.Encode(uint16(encryptedSize), bb[:0])
		return w.sizeParser.SizeBytes(), nil
	})

	w.buffer.AppendSupplier(func(bb []byte) (int, error) {
		w.auth.Seal(bb[:0], b)
		return encryptedSize, nil
	})
}

func (w *AuthenticationWriter) flush() error {
	_, err := w.writer.Write(w.buffer.Bytes())
	w.buffer.Clear()
	return err
}

func (w *AuthenticationWriter) writeStream(mb buf.MultiBuffer) error {
	defer mb.Release()

	for {
		n, _ := mb.Read(w.payload)
		w.append(w.payload[:n])
		if w.buffer.Len() > readerBufferSize-2*1024 {
			if err := w.flush(); err != nil {
				return err
			}
		}
		if mb.IsEmpty() {
			break
		}
	}

	if !w.buffer.IsEmpty() {
		return w.flush()
	}
	return nil
}

func (w *AuthenticationWriter) writePacket(mb buf.MultiBuffer) error {
	defer mb.Release()

	for {
		b := mb.SplitFirst()
		if b == nil {
			b = buf.New()
		}
		if w.buffer.Len() > readerBufferSize-b.Len()-128 {
			if err := w.flush(); err != nil {
				b.Release()
				return err
			}
		}
		w.append(b.Bytes())
		b.Release()
		if mb.IsEmpty() {
			break
		}
	}

	if !w.buffer.IsEmpty() {
		return w.flush()
	}

	return nil
}

func (w *AuthenticationWriter) Write(mb buf.MultiBuffer) error {
	if w.mode == ModeStream {
		return w.writeStream(mb)
	}

	return w.writePacket(mb)
}
