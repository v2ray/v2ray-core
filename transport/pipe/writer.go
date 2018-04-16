package pipe

import (
	"v2ray.com/core/common/buf"
)

type Writer struct {
	pipe *pipe
}

func (w *Writer) WriteMultiBuffer(mb buf.MultiBuffer) error {
	return w.pipe.WriteMultiBuffer(mb)
}

func (w *Writer) Close() error {
	return w.pipe.Close()
}

func (w *Writer) CloseError() {
	w.pipe.CloseError()
}
