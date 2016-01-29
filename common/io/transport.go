package io

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
)

func RawReaderToChan(stream chan<- *alloc.Buffer, reader io.Reader) error {
	return ReaderToChan(stream, NewAdaptiveReader(reader))
}

// ReaderToChan dumps all content from a given reader to a chan by constantly reading it until EOF.
func ReaderToChan(stream chan<- *alloc.Buffer, reader Reader) error {
	for {
		buffer, err := reader.Read()
		if alloc.Len(buffer) > 0 {
			stream <- buffer
		} else {
			alloc.Release(buffer)
		}

		if err != nil {
			return err
		}
	}
}

// ChanToWriter dumps all content from a given chan to a writer until the chan is closed.
func ChanToWriter(writer io.Writer, stream <-chan *alloc.Buffer) error {
	for buffer := range stream {
		nBytes, err := writer.Write(buffer.Value)
		if nBytes < buffer.Len() {
			_, err = writer.Write(buffer.Value[nBytes:])
		}
		buffer.Release()
		if err != nil {
			return err
		}
	}
	return nil
}
