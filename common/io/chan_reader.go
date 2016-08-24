package io

import (
	"io"
	"sync"

	"v2ray.com/core/common/alloc"
)

type ChanReader struct {
	sync.Mutex
	stream  Reader
	current *alloc.Buffer
	eof     bool
}

func NewChanReader(stream Reader) *ChanReader {
	this := &ChanReader{
		stream: stream,
	}
	this.Fill()
	return this
}

// Private: Visible for testing.
func (this *ChanReader) Fill() {
	b, err := this.stream.Read()
	this.current = b
	if err != nil {
		this.eof = true
		this.current = nil
	}
}

func (this *ChanReader) Read(b []byte) (int, error) {
	if this.eof {
		return 0, io.EOF
	}

	this.Lock()
	defer this.Unlock()
	if this.current == nil {
		this.Fill()
		if this.eof {
			return 0, io.EOF
		}
	}
	nBytes, err := this.current.Read(b)
	if this.current.IsEmpty() {
		this.current.Release()
		this.current = nil
	}
	return nBytes, err
}

func (this *ChanReader) Release() {
	this.Lock()
	defer this.Unlock()

	this.eof = true
	this.current.Release()
	this.current = nil
	this.stream = nil
}
