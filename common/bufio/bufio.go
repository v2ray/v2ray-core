// Package bufio is a replacement of the standard golang package bufio.
package bufio

import (
	"bufio"
	"io"
)

func OriginalReader(reader io.Reader) *bufio.Reader {
	return bufio.NewReader(reader)
}

func OriginalReaderSize(reader io.Reader, size int) *bufio.Reader {
	return bufio.NewReaderSize(reader, size)
}
