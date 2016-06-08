package io

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
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

// @Private
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
	nBytes := copy(b, this.current.Value)
	if nBytes == this.current.Len() {
		this.current.Release()
		this.current = nil
	} else {
		this.current.SliceFrom(nBytes)
	}
	return nBytes, nil
}

func (this *ChanReader) Release() {
	this.Lock()
	defer this.Unlock()

	this.eof = true
	this.current.Release()
	this.current = nil
	this.stream = nil
}
