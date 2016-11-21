package io

import (
	"io"
	"sync"

	"fmt"
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

func (this *BufferedWriter) ReadFrom(reader io.Reader) (int64, error) {
	this.Lock()
	defer this.Unlock()

	if this.writer == nil {
		return 0, io.ErrClosedPipe
	}

	totalBytes := int64(0)
	for {
		nBytes, err := this.buffer.FillFrom(reader)
		totalBytes += int64(nBytes)
		if err != nil {
			if err == io.EOF {
				return totalBytes, nil
			}
			return totalBytes, err
		}
		this.FlushWithoutLock()
	}
}

func (this *BufferedWriter) Write(b []byte) (int, error) {
	this.Lock()
	defer this.Unlock()

	if this.writer == nil {
		return 0, io.ErrClosedPipe
	}

	fmt.Printf("BufferedWriter writing: %v\n", b)

	if !this.cached {
		return this.writer.Write(b)
	}
	nBytes, _ := this.buffer.Write(b)
	if this.buffer.IsFull() {
		this.FlushWithoutLock()
	}
	fmt.Printf("BufferedWriter content: %v\n", this.buffer.Value)
	return nBytes, nil
}

func (this *BufferedWriter) Flush() error {
	this.Lock()
	defer this.Unlock()

	if this.writer == nil {
		return io.ErrClosedPipe
	}

	return this.FlushWithoutLock()
}

func (this *BufferedWriter) FlushWithoutLock() error {
	fmt.Println("BufferedWriter flushing")
	defer this.buffer.Clear()
	for !this.buffer.IsEmpty() {
		nBytes, err := this.writer.Write(this.buffer.Value)
		fmt.Printf("BufferedWriting flushed %d bytes.\n", nBytes)
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
