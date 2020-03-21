// +build illumos

package buf

import "golang.org/x/sys/unix"

type unixReader struct {
	iovs [][]byte
}

func (r *unixReader) Init(bs []*Buffer) {
	iovs := r.iovs
	if iovs == nil {
		iovs = make([][]byte, 0, len(bs))
	}
	for _, b := range bs {
		iovs = append(iovs, b.v)
	}
	r.iovs = iovs
}

func (r *unixReader) Read(fd uintptr) int32 {
	n, e := unix.Readv(int(fd), r.iovs)
	if e != nil {
		return -1
	}
	return int32(n)
}

func (r *unixReader) Clear() {
	r.iovs = r.iovs[:0]
}

func newMultiReader() multiReader {
	return &unixReader{}
}
