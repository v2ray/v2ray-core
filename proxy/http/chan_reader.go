package http

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
)

type ChanReader struct {
	stream  v2io.Reader
	current *alloc.Buffer
	eof     bool
}

func NewChanReader(stream v2io.Reader) *ChanReader {
	this := &ChanReader{
		stream: stream,
	}
	this.fill()
	return this
}

func (this *ChanReader) fill() {
	b, err := this.stream.Read()
	this.current = b
	if err != nil {
		this.eof = true
		this.current = nil
	}
}

func (this *ChanReader) Read(b []byte) (int, error) {
	if this.current == nil {
		this.fill()
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
