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
		id:     id,
		dest:   dest,
		writer: writer,
	}
}

func (w *Writer) writeInternal(b *buf.Buffer) error {
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

	if b.Len() > 0 {
		meta.Option.Add(OptionData)
	}

	frame := buf.New()
	frame.AppendSupplier(meta.AsSupplier())

	if b.Len() > 0 {
		frame.AppendSupplier(serial.WriteUint16(0))
		lengthBytes := frame.BytesFrom(-2)

		nBytes, err := frame.Write(b.Bytes())
		if err != nil {
			frame.Release()
			return err
		}

		serial.Uint16ToBytes(uint16(nBytes), lengthBytes[:0])
		b.SliceFrom(nBytes)
	}

	return w.writer.Write(frame)
}

func (w *Writer) Write(b *buf.Buffer) error {
	defer b.Release()

	if err := w.writeInternal(b); err != nil {
		return err
	}
	for !b.IsEmpty() {
		if err := w.writeInternal(b); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Close() {
	meta := FrameMetadata{
		SessionID:     w.id,
		Target:        w.dest,
		SessionStatus: SessionStatusEnd,
	}

	frame := buf.New()
	frame.AppendSupplier(meta.AsSupplier())

	w.writer.Write(frame)
}
