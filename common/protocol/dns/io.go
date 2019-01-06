package dns

import (
	"sync"

	"golang.org/x/net/dns/dnsmessage"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
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
