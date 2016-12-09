package io

import (
	"io"
	"sync"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
)

type BufferedWriter struct {
	sync.Mutex
	writer io.Writer
	buffer *buf.Buffer
	cached bool
}

func NewBufferedWriter(rawWriter io.Writer) *BufferedWriter {
	return &BufferedWriter{
		writer: rawWriter,
		buffer: buf.NewSmallBuffer(),
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
			if errors.Cause(err) == io.EOF {
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
	nBytes, err := v.buffer.Write(b)
	if err != nil {
		return 0, err
	}
	if v.buffer.IsFull() {
		err := v.FlushWithoutLock()
		if err != nil {
			return 0, err
		}
		if nBytes < len(b) {
			if _, err := v.writer.Write(b[nBytes:]); err != nil {
				return nBytes, err
			}
		}
	}
	return len(b), nil
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
		nBytes, err := v.writer.Write(v.buffer.Bytes())
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
