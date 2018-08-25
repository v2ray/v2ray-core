package buf

import (
	"io"
	"runtime"
	"syscall"

	"v2ray.com/core/common/platform"
)

type allocStrategy struct {
	current uint32
}

func (s *allocStrategy) Current() uint32 {
	return s.current
}

func (s *allocStrategy) Adjust(n uint32) {
	if n >= s.current {
		s.current *= 4
	} else {
		s.current = n
	}

	if s.current > 32 {
		s.current = 32
	}

	if s.current == 0 {
		s.current = 1
	}
}

func (s *allocStrategy) Alloc() []*Buffer {
	bs := make([]*Buffer, s.current)
	for i := range bs {
		bs[i] = New()
	}
	return bs
}

type multiReader interface {
	Init([]*Buffer)
	Read(fd uintptr) int32
	Clear()
}

type ReadVReader struct {
	io.Reader
	rawConn syscall.RawConn
	mr      multiReader
	alloc   allocStrategy
}

func NewReadVReader(reader io.Reader, rawConn syscall.RawConn) *ReadVReader {
	return &ReadVReader{
		Reader:  reader,
		rawConn: rawConn,
		alloc: allocStrategy{
			current: 1,
		},
		mr: newMultiReader(),
	}
}

func (r *ReadVReader) readMulti() (MultiBuffer, error) {
	bs := r.alloc.Alloc()

	r.mr.Init(bs)
	var nBytes int32
	err := r.rawConn.Read(func(fd uintptr) bool {
		n := r.mr.Read(fd)
		if n < 0 {
			return false
		}

		nBytes = n
		return true
	})
	r.mr.Clear()

	if err != nil {
		mb := MultiBuffer(bs)
		mb.Release()
		return nil, err
	}

	if nBytes == 0 {
		mb := MultiBuffer(bs)
		mb.Release()
		return nil, io.EOF
	}

	nBuf := 0
	for nBuf < len(bs) {
		if nBytes <= 0 {
			break
		}
		end := int32(nBytes)
		if end > Size {
			end = Size
		}
		bs[nBuf].end = end
		nBytes -= end
		nBuf++
	}

	for i := nBuf; i < len(bs); i++ {
		bs[i].Release()
		bs[i] = nil
	}

	return MultiBuffer(bs[:nBuf]), nil
}

// ReadMultiBuffer implements Reader.
func (r *ReadVReader) ReadMultiBuffer() (MultiBuffer, error) {
	if r.alloc.Current() == 1 {
		b, err := readOne(r.Reader)
		if err != nil {
			return nil, err
		}
		if b.IsFull() {
			r.alloc.Adjust(1)
		}
		return NewMultiBufferValue(b), nil
	}

	mb, err := r.readMulti()
	if err != nil {
		return nil, err
	}
	r.alloc.Adjust(uint32(len(mb)))
	return mb, nil
}

var useReadv = false

func init() {
	const defaultFlagValue = "NOT_DEFINED_AT_ALL"
	value := platform.NewEnvFlag("v2ray.buf.readv").GetValue(func() string { return defaultFlagValue })
	switch value {
	case defaultFlagValue, "auto":
		if (runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "s390x") && (runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
			useReadv = true
		}
	case "enable":
		useReadv = true
	}
}
