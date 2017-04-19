package mux

import (
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
)

type Writer struct {
	id       uint16
	dest     net.Destination
	writer   buf.Writer
	followup bool
}

func NewWriter(id uint16, dest net.Destination, writer buf.Writer) *Writer {
	return &Writer{
		id:       id,
		dest:     dest,
		writer:   writer,
		followup: false,
	}
}

func NewResponseWriter(id uint16, writer buf.Writer) *Writer {
	return &Writer{
		id:       id,
		writer:   writer,
		followup: true,
	}
}

func (w *Writer) writeInternal(mb buf.MultiBuffer) error {
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

	hasData := !mb.IsEmpty()

	if hasData {
		meta.Option.Add(OptionData)
	}

	frame := buf.New()
	frame.AppendSupplier(meta.AsSupplier())

	mb2 := buf.NewMultiBuffer()
	mb2.Append(frame)

	if hasData {
		frame.AppendSupplier(serial.WriteUint16(uint16(mb.Len())))
		mb2.AppendMulti(mb)
	}
	return w.writer.Write(mb2)
}

func (w *Writer) Write(mb buf.MultiBuffer) error {
	const chunkSize = 8 * 1024
	for {
		slice := mb.SliceBySize(chunkSize)
		if err := w.writeInternal(slice); err != nil {
			return err
		}
		if mb.IsEmpty() {
			break
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

	w.writer.Write(buf.NewMultiBufferValue(frame))
}
