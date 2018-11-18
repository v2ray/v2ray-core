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
	return w.writer.WriteMultiBuffer(buf.MultiBuffer{b})
}

func writeMetaWithFrame(writer buf.Writer, meta FrameMetadata, data buf.MultiBuffer) error {
	frame := buf.New()
	if err := meta.WriteTo(frame); err != nil {
		return err
	}
	if _, err := serial.WriteUint16(frame, uint16(data.Len())); err != nil {
		return err
	}

	mb2 := make(buf.MultiBuffer, 0, len(data)+1)
	mb2 = append(mb2, frame)
	mb2 = append(mb2, data...)
	return writer.WriteMultiBuffer(mb2)
}

func (w *Writer) writeData(mb buf.MultiBuffer) error {
	meta := w.getNextFrameMeta()
	meta.Option.Set(OptionData)

	return writeMetaWithFrame(w.writer, meta, mb)
}

// WriteMultiBuffer implements buf.Writer.
func (w *Writer) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer buf.ReleaseMulti(mb)

	if mb.IsEmpty() {
		return w.writeMetaOnly()
	}

	for !mb.IsEmpty() {
		var chunk buf.MultiBuffer
		if w.transferType == protocol.TransferTypeStream {
			mb, chunk = buf.SplitSize(mb, 8*1024)
		} else {
			mb2, b := buf.SplitFirst(mb)
			mb = mb2
			chunk = buf.MultiBuffer{b}
		}
		if err := w.writeData(chunk); err != nil {
			return err
		}
	}

	return nil
}

// Close implements common.Closable.
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

	w.writer.WriteMultiBuffer(buf.MultiBuffer{frame}) // nolint: errcheck
	return nil
}
