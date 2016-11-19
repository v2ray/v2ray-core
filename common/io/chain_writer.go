package io

import (
	"io"
	"sync"

	"v2ray.com/core/common/alloc"
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

func (this *ChainWriter) Write(payload []byte) (int, error) {
	this.Lock()
	defer this.Unlock()
	if this.writer == nil {
		return 0, io.ErrClosedPipe
	}

	size := len(payload)
	for size > 0 {
		buffer := alloc.NewBuffer().Clear()
		if size > alloc.BufferSize {
			buffer.Append(payload[:alloc.BufferSize])
			size -= alloc.BufferSize
		} else {
			buffer.Append(payload)
			size = 0
		}
		err := this.writer.Write(buffer)
		if err != nil {
			return 0, err
		}
	}

	return size, nil
}

func (this *ChainWriter) Release() {
	this.Lock()
	this.writer.Release()
	this.writer = nil
	this.Unlock()
}
