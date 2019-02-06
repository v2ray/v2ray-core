package dns

import (
	"encoding/binary"
	"sync"

	"golang.org/x/net/dns/dnsmessage"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

func PackMessage(msg *dnsmessage.Message) (*buf.Buffer, error) {
	buffer := buf.New()
	rawBytes := buffer.Extend(buf.Size)
	packed, err := msg.AppendPack(rawBytes[:0])
	if err != nil {
		buffer.Release()
		return nil, err
	}
	buffer.Resize(0, int32(len(packed)))
	return buffer, nil
}

type MessageReader interface {
	ReadMessage() (*buf.Buffer, error)
}

type UDPReader struct {
	buf.Reader

	access sync.Mutex
	cache  buf.MultiBuffer
}

func (r *UDPReader) readCache() *buf.Buffer {
	r.access.Lock()
	defer r.access.Unlock()

	mb, b := buf.SplitFirst(r.cache)
	r.cache = mb
	return b
}

func (r *UDPReader) refill() error {
	mb, err := r.Reader.ReadMultiBuffer()
	if err != nil {
		return err
	}
	r.access.Lock()
	r.cache = mb
	r.access.Unlock()
	return nil
}

// ReadMessage implements MessageReader.
func (r *UDPReader) ReadMessage() (*buf.Buffer, error) {
	for {
		b := r.readCache()
		if b != nil {
			return b, nil
		}
		if err := r.refill(); err != nil {
			return nil, err
		}
	}
}

// Close implements common.Closable.
func (r *UDPReader) Close() error {
	defer func() {
		r.access.Lock()
		buf.ReleaseMulti(r.cache)
		r.cache = nil
		r.access.Unlock()
	}()

	return common.Close(r.Reader)
}

type TCPReader struct {
	reader *buf.BufferedReader
}

func NewTCPReader(reader buf.Reader) *TCPReader {
	return &TCPReader{
		reader: &buf.BufferedReader{
			Reader: reader,
		},
	}
}

func (r *TCPReader) ReadMessage() (*buf.Buffer, error) {
	size, err := serial.ReadUint16(r.reader)
	if err != nil {
		return nil, err
	}
	if size > buf.Size {
		return nil, newError("message size too large: ", size)
	}
	b := buf.New()
	if _, err := b.ReadFullFrom(r.reader, int32(size)); err != nil {
		return nil, err
	}
	return b, nil
}

func (r *TCPReader) Interrupt() {
	common.Interrupt(r.reader)
}

func (r *TCPReader) Close() error {
	return common.Close(r.reader)
}

type MessageWriter interface {
	WriteMessage(msg *buf.Buffer) error
}

type UDPWriter struct {
	buf.Writer
}

func (w *UDPWriter) WriteMessage(b *buf.Buffer) error {
	return w.WriteMultiBuffer(buf.MultiBuffer{b})
}

type TCPWriter struct {
	buf.Writer
}

func (w *TCPWriter) WriteMessage(b *buf.Buffer) error {
	if b.IsEmpty() {
		return nil
	}

	mb := make(buf.MultiBuffer, 0, 2)

	size := buf.New()
	binary.BigEndian.PutUint16(size.Extend(2), uint16(b.Len()))
	mb = append(mb, size, b)
	return w.WriteMultiBuffer(mb)
}
