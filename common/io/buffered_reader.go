package io

import (
	"io"
	"sync"

	"v2ray.com/core/common/buf"
)

type BufferedReader struct {
	sync.Mutex
	reader io.Reader
	buffer *buf.Buffer
	cached bool
}

func NewBufferedReader(rawReader io.Reader) *BufferedReader {
	return &BufferedReader{
		reader: rawReader,
		buffer: buf.New(),
		cached: true,
	}
}

func (v *BufferedReader) Release() {
	v.Lock()
	defer v.Unlock()

	v.buffer.Release()
	v.buffer = nil
	v.reader = nil
}

func (v *BufferedReader) Cached() bool {
	return v.cached
}

func (v *BufferedReader) SetCached(cached bool) {
	v.cached = cached
}

func (v *BufferedReader) Read(b []byte) (int, error) {
	v.Lock()
	defer v.Unlock()

	if v.reader == nil {
		return 0, io.EOF
	}

	if !v.cached {
		if !v.buffer.IsEmpty() {
			return v.buffer.Read(b)
		}
		return v.reader.Read(b)
	}
	if v.buffer.IsEmpty() {
		err := v.buffer.AppendSupplier(buf.ReadFrom(v.reader))
		if err != nil {
			return 0, err
		}
	}

	if v.buffer.IsEmpty() {
		return 0, nil
	}

	return v.buffer.Read(b)
}
