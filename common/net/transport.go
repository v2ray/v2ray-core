package net

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
)

func ReadFrom(reader io.Reader, buffer *alloc.Buffer) (*alloc.Buffer, error) {
	if buffer == nil {
		buffer = alloc.NewBuffer()
	}
	nBytes, err := reader.Read(buffer.Value)
	buffer.Slice(0, nBytes)
	return buffer, err
}

// ReaderToChan dumps all content from a given reader to a chan by constantly reading it until EOF.
func ReaderToChan(stream chan<- *alloc.Buffer, reader io.Reader) error {
	for {
		buffer, err := ReadFrom(reader, nil)
		if buffer.Len() > 0 {
			stream <- buffer
		} else {
			buffer.Release()
		}
		if err != nil {
			return err
		}
	}
}

// ChanToWriter dumps all content from a given chan to a writer until the chan is closed.
func ChanToWriter(writer io.Writer, stream <-chan *alloc.Buffer) error {
	for buffer := range stream {
		_, err := writer.Write(buffer.Value)
		buffer.Release()
		if err != nil {
			return err
		}
	}
	return nil
}
