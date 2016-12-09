package io

import (
	"io"
	"sync"

	"v2ray.com/core/common/buf"
)

type ChainWriter struct {
	sync.Mutex
	writer Writer
}

func NewChainWriter(writer Writer) *ChainWriter {
	return &ChainWriter{
		writer: writer,
	}
}

func (v *ChainWriter) Write(payload []byte) (int, error) {
	v.Lock()
	defer v.Unlock()
	if v.writer == nil {
		return 0, io.ErrClosedPipe
	}

	bytesWritten := 0
	size := len(payload)
	for size > 0 {
		buffer := buf.NewBuffer()
		nBytes, _ := buffer.Write(payload)
		size -= nBytes
		payload = payload[nBytes:]
		bytesWritten += nBytes
		err := v.writer.Write(buffer)
		if err != nil {
			return bytesWritten, err
		}
	}

	return bytesWritten, nil
}

func (v *ChainWriter) Release() {
	v.Lock()
	v.writer.Release()
	v.writer = nil
	v.Unlock()
}
