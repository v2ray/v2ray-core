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
	return &ChanReader{
		stream: stream,
	}
}

// Private: Visible for testing.
func (v *ChanReader) Fill() {
	b, err := v.stream.Read()
	v.current = b
	if err != nil {
		v.eof = true
		v.current = nil
	}
}

func (v *ChanReader) Read(b []byte) (int, error) {
	if v.eof {
		return 0, io.EOF
	}

	v.Lock()
	defer v.Unlock()
	if v.current == nil {
		v.Fill()
		if v.eof {
			return 0, io.EOF
		}
	}
	nBytes, err := v.current.Read(b)
	if v.current.IsEmpty() {
		v.current.Release()
		v.current = nil
	}
	return nBytes, err
}

func (v *ChanReader) Release() {
	v.Lock()
	defer v.Unlock()

	v.eof = true
	v.current.Release()
	v.current = nil
	v.stream = nil
}
