// +build !windows

package buf

import (
	"io"
	"runtime"
	"syscall"
	"unsafe"

	"v2ray.com/core/common/platform"
)

type ReadVReader struct {
	io.Reader
	rawConn syscall.RawConn
	iovects []syscall.Iovec
	nBuf    int32
}

func NewReadVReader(reader io.Reader, rawConn syscall.RawConn) *ReadVReader {
	return &ReadVReader{
		Reader:  reader,
		rawConn: rawConn,
		nBuf:    1,
	}
}

func allocN(n int32) []*Buffer {
	bs := make([]*Buffer, 0, n)
	for i := int32(0); i < n; i++ {
		bs = append(bs, New())
	}
	return bs
}

func (r *ReadVReader) ReadMultiBuffer() (MultiBuffer, error) {
	bs := allocN(r.nBuf)

	var iovecs []syscall.Iovec
	if r.iovects != nil {
		iovecs = r.iovects
	}
	for idx, b := range bs {
		iovecs = append(iovecs, syscall.Iovec{
			Base: &(b.v[0]),
		})
		iovecs[idx].SetLen(int(Size))
	}
	r.iovects = iovecs[:0]

	var nBytes int

	err := r.rawConn.Read(func(fd uintptr) bool {
		n, _, e := syscall.Syscall(syscall.SYS_READV, fd, uintptr(unsafe.Pointer(&iovecs[0])), uintptr(len(iovecs)))
		if e != 0 {
			return false
		}
		nBytes = int(n)
		return true
	})

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

	var isFull bool = (nBytes == int(r.nBuf)*Size)

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
		nBytes -= int(end)
		nBuf++
	}

	for i := nBuf; i < len(bs); i++ {
		bs[i].Release()
		bs[i] = nil
	}

	if isFull && nBuf < 8 {
		r.nBuf *= 4
	} else {
		r.nBuf = int32(nBuf)
	}

	return MultiBuffer(bs[:nBuf]), nil
}

var useReadv = false

func init() {
	const defaultFlagValue = "NOT_DEFINED_AT_ALL"
	value := platform.NewEnvFlag("v2ray.buf.readv").GetValue(func() string { return defaultFlagValue })
	if value != defaultFlagValue && (runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
		useReadv = true
	}
}
