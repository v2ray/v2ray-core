package mux

import "v2ray.com/core/common/buf"
import "v2ray.com/core/common/serial"

type muxWriter struct {
	meta   *FrameMetadata
	writer buf.Writer
}

func (w *muxWriter) Write(b *buf.Buffer) error {
	frame := buf.New()
	frame.AppendSupplier(w.meta.AsSupplier())
	if w.meta.SessionStatus == SessionStatusNew {
		w.meta.SessionStatus = SessionStatusKeep
	}

	frame.AppendSupplier(serial.WriteUint16(0))
	lengthBytes := frame.BytesFrom(-2)

	nBytes, err := frame.Write(b.Bytes())
	if err != nil {
		return err
	}

	serial.Uint16ToBytes(uint16(nBytes), lengthBytes[:0])
	if err := w.writer.Write(frame); err != nil {
		frame.Release()
		b.Release()
		return err
	}

	b.SliceFrom(nBytes)
	if !b.IsEmpty() {
		return w.Write(b)
	}
	b.Release()

	return nil
}
