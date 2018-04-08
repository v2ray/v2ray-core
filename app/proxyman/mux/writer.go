package mux

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
)

type Writer struct {
	dest         net.Destination
	writer       buf.Writer
	id           uint16
	followup     bool
	hasError     bool
	transferType protocol.TransferType
}

func NewWriter(id uint16, dest net.Destination, writer buf.Writer, transferType protocol.TransferType) *Writer {
	return &Writer{
		id:           id,
		dest:         dest,
		writer:       writer,
		followup:     false,
		transferType: transferType,
	}
}

func NewResponseWriter(id uint16, writer buf.Writer, transferType protocol.TransferType) *Writer {
	return &Writer{
		id:           id,
		writer:       writer,
		followup:     true,
		transferType: transferType,
	}
}

func (w *Writer) getNextFrameMeta() FrameMetadata {
	meta := FrameMetadata{
		SessionID: w.id,
		Target:    w.dest,
	}

	if w.followup {
		meta.SessionStatus = SessionStatusKeep
	} else {
		w.followup = true
		meta.SessionStatus = SessionStatusNew
	}

	return meta
}

func (w *Writer) writeMetaOnly() error {
	meta := w.getNextFrameMeta()
	b := buf.New()
	if err := meta.WriteTo(b); err != nil {
		return err
	}
	return w.writer.WriteMultiBuffer(buf.NewMultiBufferValue(b))
}

func (w *Writer) writeData(mb buf.MultiBuffer) error {
	meta := w.getNextFrameMeta()
	meta.Option.Set(OptionData)

	frame := buf.New()
	if err := meta.WriteTo(frame); err != nil {
		return err
	}
	if err := frame.AppendSupplier(serial.WriteUint16(uint16(mb.Len()))); err != nil {
		return err
	}

	mb2 := buf.NewMultiBufferCap(int32(len(mb)) + 1)
	mb2.Append(frame)
	mb2.AppendMulti(mb)
	return w.writer.WriteMultiBuffer(mb2)
}

// WriteMultiBuffer implements buf.Writer.
func (w *Writer) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer mb.Release()

	if mb.IsEmpty() {
		return w.writeMetaOnly()
	}

	for !mb.IsEmpty() {
		var chunk buf.MultiBuffer
		if w.transferType == protocol.TransferTypeStream {
			chunk = mb.SliceBySize(8 * 1024)
		} else {
			chunk = buf.NewMultiBufferValue(mb.SplitFirst())
		}
		if err := w.writeData(chunk); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) Close() error {
	meta := FrameMetadata{
		SessionID:     w.id,
		SessionStatus: SessionStatusEnd,
	}
	if w.hasError {
		meta.Option.Set(OptionError)
	}

	frame := buf.New()
	common.Must(meta.WriteTo(frame))

	w.writer.WriteMultiBuffer(buf.NewMultiBufferValue(frame))
	return nil
}
