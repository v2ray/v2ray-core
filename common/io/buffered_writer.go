package io

import (
	"io"
	"sync"
	"v2ray.com/core/common/alloc"
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

func (v *BufferedWriter) ReadFrom(reader io.Reader) (int64, error) {
	v.Lock()
	defer v.Unlock()

	if v.writer == nil {
		return 0, io.ErrClosedPipe
	}

	totalBytes := int64(0)
	for {
		nBytes, err := v.buffer.FillFrom(reader)
		totalBytes += int64(nBytes)
		if err != nil {
			if err == io.EOF {
				return totalBytes, nil
			}
			return totalBytes, err
		}
		v.FlushWithoutLock()
	}
}

func (v *BufferedWriter) Write(b []byte) (int, error) {
	v.Lock()
	defer v.Unlock()

	if v.writer == nil {
		return 0, io.ErrClosedPipe
	}

	if !v.cached {
		return v.writer.Write(b)
	}
	nBytes, _ := v.buffer.Write(b)
	if v.buffer.IsFull() {
		v.FlushWithoutLock()
	}
	return nBytes, nil
}

func (v *BufferedWriter) Flush() error {
	v.Lock()
	defer v.Unlock()

	if v.writer == nil {
		return io.ErrClosedPipe
	}

	return v.FlushWithoutLock()
}

func (v *BufferedWriter) FlushWithoutLock() error {
	defer v.buffer.Clear()
	for !v.buffer.IsEmpty() {
		nBytes, err := v.writer.Write(v.buffer.Value)
		if err != nil {
			return err
		}
		v.buffer.SliceFrom(nBytes)
	}
	return nil
}

func (v *BufferedWriter) Cached() bool {
	return v.cached
}

func (v *BufferedWriter) SetCached(cached bool) {
	v.cached = cached
	if !cached && !v.buffer.IsEmpty() {
		v.Flush()
	}
}

func (v *BufferedWriter) Release() {
	v.Flush()

	v.Lock()
	defer v.Unlock()

	v.buffer.Release()
	v.buffer = nil
	v.writer = nil
}
