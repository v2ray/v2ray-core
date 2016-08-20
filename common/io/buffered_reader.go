package io

import (
	"io"
	"sync"

	"v2ray.com/core/common/alloc"
)

type BufferedReader struct {
	sync.Mutex
	reader io.Reader
	buffer *alloc.Buffer
	cached bool
}

func NewBufferedReader(rawReader io.Reader) *BufferedReader {
	return &BufferedReader{
		reader: rawReader,
		buffer: alloc.NewBuffer().Clear(),
		cached: true,
	}
}

func (this *BufferedReader) Release() {
	this.Lock()
	defer this.Unlock()

	this.buffer.Release()
	this.buffer = nil
	this.reader = nil
}

func (this *BufferedReader) Cached() bool {
	return this.cached
}

func (this *BufferedReader) SetCached(cached bool) {
	this.cached = cached
}

func (this *BufferedReader) Read(b []byte) (int, error) {
	this.Lock()
	defer this.Unlock()

	if this.reader == nil {
		return 0, io.EOF
	}

	if !this.cached {
		if !this.buffer.IsEmpty() {
			return this.buffer.Read(b)
		}
		return this.reader.Read(b)
	}
	if this.buffer.IsEmpty() {
		_, err := this.buffer.FillFrom(this.reader)
		if err != nil {
			return 0, err
		}
	}

	if this.buffer.IsEmpty() {
		return 0, nil
	}

	return this.buffer.Read(b)
}
