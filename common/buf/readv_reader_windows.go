// +build windows
package buf

import (
	"io"
	"syscall"
)

var useReadv = false

func NewReadVReader(reader io.Reader, rawConn syscall.RawConn) *ReadVReader {
	return nil
}
