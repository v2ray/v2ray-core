package io

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
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
	if this.writer == nil {
		return 0, io.EOF
	}

	size := len(payload)
	buffer := alloc.NewBufferWithSize(size).Clear()
	buffer.Append(payload)

	this.Lock()
	defer this.Unlock()
	if this.writer == nil {
		return 0, io.EOF
	}

	err := this.writer.Write(buffer)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (this *ChainWriter) Release() {
	this.Lock()
	this.writer.Release()
	this.writer = nil
	this.Unlock()
}
