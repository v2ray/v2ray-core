package mux

import (
	"runtime"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
)

type Writer struct {
	id           uint16
	dest         net.Destination
	writer       buf.Writer
	followup     bool
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
	if err := b.AppendSupplier(meta.AsSupplier()); err != nil {
		return err
	}
	runtime.KeepAlive(meta)
	return w.writer.Write(buf.NewMultiBufferValue(b))
}

func (w *Writer) writeData(mb buf.MultiBuffer) error {
	meta := w.getNextFrameMeta()
	meta.Option.Add(OptionData)

	frame := buf.New()
	if err := frame.AppendSupplier(meta.AsSupplier()); err != nil {
		return err
	}
	runtime.KeepAlive(meta)
	if err := frame.AppendSupplier(serial.WriteUint16(uint16(mb.Len()))); err != nil {
		return err
	}

	mb2 := buf.NewMultiBuffer()
	mb2.Append(frame)
	mb2.AppendMulti(mb)
	return w.writer.Write(mb2)
}

// Write implements buf.MultiBufferWriter.
func (w *Writer) Write(mb buf.MultiBuffer) error {
	if mb.IsEmpty() {
		return w.writeMetaOnly()
	}

	if w.transferType == protocol.TransferTypeStream {
		const chunkSize = 8 * 1024
		for !mb.IsEmpty() {
			slice := mb.SliceBySize(chunkSize)
			if err := w.writeData(slice); err != nil {
				return err
			}
		}
	} else {
		for _, b := range mb {
			if err := w.writeData(buf.NewMultiBufferValue(b)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Writer) Close() {
	meta := FrameMetadata{
		SessionID:     w.id,
		SessionStatus: SessionStatusEnd,
	}

	frame := buf.New()
	frame.AppendSupplier(meta.AsSupplier())
	runtime.KeepAlive(meta)

	w.writer.Write(buf.NewMultiBufferValue(frame))
}
