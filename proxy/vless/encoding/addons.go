// +build !confonly

package encoding

import (
	"io"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/vless"
)

func EncodeHeaderAddons(buffer *buf.Buffer, addons *Addons) error {

	switch addons.Flow {
	case vless.XRO, vless.XRD:

		bytes, err := proto.Marshal(addons)
		if err != nil {
			return newError("failed to marshal addons protobuf value").Base(err)
		}
		if err := buffer.WriteByte(byte(len(bytes))); err != nil {
			return newError("failed to write addons protobuf length").Base(err)
		}
		if _, err := buffer.Write(bytes); err != nil {
			return newError("failed to write addons protobuf value").Base(err)
		}

	default:

		if err := buffer.WriteByte(0); err != nil {
			return newError("failed to write addons protobuf length").Base(err)
		}

	}

	return nil

}

func DecodeHeaderAddons(buffer *buf.Buffer, reader io.Reader) (*Addons, error) {

	addons := new(Addons)

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		return nil, newError("failed to read addons protobuf length").Base(err)
	}

	if length := int32(buffer.Byte(0)); length != 0 {

		buffer.Clear()
		if _, err := buffer.ReadFullFrom(reader, length); err != nil {
			return nil, newError("failed to read addons protobuf value").Base(err)
		}

		if err := proto.Unmarshal(buffer.Bytes(), addons); err != nil {
			return nil, newError("failed to unmarshal addons protobuf value").Base(err)
		}

		// Verification.
		switch addons.Flow {
		default:

		}

	}

	return addons, nil

}

// EncodeBodyAddons returns a Writer that auto-encrypt content written by caller.
func EncodeBodyAddons(writer io.Writer, request *protocol.RequestHeader, addons *Addons) buf.Writer {

	switch addons.Flow {
	default:

		if request.Command == protocol.RequestCommandUDP {
			return NewMultiLengthPacketWriter(writer.(buf.Writer))
		}

	}

	return buf.NewWriter(writer)

}

// DecodeBodyAddons returns a Reader from which caller can fetch decrypted body.
func DecodeBodyAddons(reader io.Reader, request *protocol.RequestHeader, addons *Addons) buf.Reader {

	switch addons.Flow {
	default:

		if request.Command == protocol.RequestCommandUDP {
			return NewLengthPacketReader(reader)
		}

	}

	return buf.NewReader(reader)

}

func NewMultiLengthPacketWriter(writer buf.Writer) *MultiLengthPacketWriter {
	return &MultiLengthPacketWriter{
		Writer: writer,
	}
}

type MultiLengthPacketWriter struct {
	buf.Writer
}

func (w *MultiLengthPacketWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer buf.ReleaseMulti(mb)
	mb2Write := make(buf.MultiBuffer, 0, len(mb)+1)
	for _, b := range mb {
		length := b.Len()
		if length == 0 || length+2 > buf.Size {
			continue
		}
		eb := buf.New()
		if err := eb.WriteByte(byte(length >> 8)); err != nil {
			eb.Release()
			continue
		}
		if err := eb.WriteByte(byte(length)); err != nil {
			eb.Release()
			continue
		}
		if _, err := eb.Write(b.Bytes()); err != nil {
			eb.Release()
			continue
		}
		mb2Write = append(mb2Write, eb)
	}
	if mb2Write.IsEmpty() {
		return nil
	}
	return w.Writer.WriteMultiBuffer(mb2Write)
}

func NewLengthPacketWriter(writer io.Writer) *LengthPacketWriter {
	return &LengthPacketWriter{
		Writer: writer,
		cache:  make([]byte, 0, 65536),
	}
}

type LengthPacketWriter struct {
	io.Writer
	cache []byte
}

func (w *LengthPacketWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	length := mb.Len() // none of mb is nil
	//fmt.Println("Write", length)
	if length == 0 {
		return nil
	}
	defer func() {
		w.cache = w.cache[:0]
	}()
	w.cache = append(w.cache, byte(length>>8), byte(length))
	for i, b := range mb {
		w.cache = append(w.cache, b.Bytes()...)
		b.Release()
		mb[i] = nil
	}
	if _, err := w.Write(w.cache); err != nil {
		return newError("failed to write a packet").Base(err)
	}
	return nil
}

func NewLengthPacketReader(reader io.Reader) *LengthPacketReader {
	return &LengthPacketReader{
		Reader: reader,
		cache:  make([]byte, 2),
	}
}

type LengthPacketReader struct {
	io.Reader
	cache []byte
}

func (r *LengthPacketReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if _, err := io.ReadFull(r.Reader, r.cache); err != nil { // maybe EOF
		return nil, newError("failed to read packet length").Base(err)
	}
	length := int32(r.cache[0])<<8 | int32(r.cache[1])
	//fmt.Println("Read", length)
	mb := make(buf.MultiBuffer, 0, length/buf.Size+1)
	for length > 0 {
		size := length
		if size > buf.Size {
			size = buf.Size
		}
		length -= size
		b := buf.New()
		if _, err := b.ReadFullFrom(r.Reader, size); err != nil {
			return nil, newError("failed to read packet payload").Base(err)
		}
		mb = append(mb, b)
	}
	return mb, nil
}
