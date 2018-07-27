// +build windows

package buf

import (
	"io"
	"syscall"
)

const useReadv = false

func NewReadVReader(reader io.Reader, rawConn syscall.RawConn) Reader {
	panic("Shoud not happen")
}
