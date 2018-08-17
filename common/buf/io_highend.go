// +build 386 amd64 s390 s390x

package buf

import (
	"io"
	"syscall"
)

func newReaderPlatform(reader io.Reader) Reader {
	if useReadv {
		if sc, ok := reader.(syscall.Conn); ok {
			rawConn, err := sc.SyscallConn()
			if err != nil {
				newError("failed to get sysconn").Base(err).WriteToLog()
			} else {
				return NewReadVReader(reader, rawConn)
			}
		}
	}

	return NewBytesToBufferReader(reader)
}
