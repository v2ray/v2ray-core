// Package bufio is a replacement of the standard golang package bufio.
package bufio

import (
	"bufio"
	"io"
)

// OriginalReader invokes bufio.NewReader() from Golang standard library.
func OriginalReader(reader io.Reader) *bufio.Reader {
	return bufio.NewReader(reader)
}

// OriginalReaderSize invokes bufio.NewReaderSize() from Golang standard library.
func OriginalReaderSize(reader io.Reader, size int) *bufio.Reader {
	return bufio.NewReaderSize(reader, size)
}
