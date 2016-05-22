package io

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
)

type BufferedWriter struct {
	sync.Mutex
	writer io.Writer
	buffer *alloc.Buffer
	cached bool
}

func NewBufferedWriter(rawWriter io.Writer) *BufferedWriter {
	return &BufferedWriter{
		writer: rawWriter,
		buffer: alloc.NewBuffer().Clear(),
		cached: true,
	}
}

func (this *BufferedWriter) Write(b []byte) (int, error) {
	this.Lock()
	defer this.Unlock()

	if this.writer == nil {
		return 0, io.EOF
	}

	if !this.cached {
		return this.writer.Write(b)
	}
	nBytes, _ := this.buffer.Write(b)
	if this.buffer.IsFull() {
		go this.Flush()
	}
	return nBytes, nil
}

func (this *BufferedWriter) Flush() error {
	this.Lock()
	defer this.Unlock()
	if this.writer == nil {
		return io.EOF
	}

	defer this.buffer.Clear()
	for !this.buffer.IsEmpty() {
		nBytes, err := this.writer.Write(this.buffer.Value)
		if err != nil {
			return err
		}
		this.buffer.SliceFrom(nBytes)
	}
	return nil
}

func (this *BufferedWriter) Cached() bool {
	return this.cached
}

func (this *BufferedWriter) SetCached(cached bool) {
	this.cached = cached
	if !cached && !this.buffer.IsEmpty() {
		this.Flush()
	}
}

func (this *BufferedWriter) Release() {
	this.Flush()

	this.Lock()
	defer this.Unlock()

	this.buffer.Release()
	this.buffer = nil
	this.writer = nil
}
