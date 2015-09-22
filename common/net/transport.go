package net

import (
	"io"
)

const (
	bufferSize = 4 * 1024
)

// ReaderToChan dumps all content from a given reader to a chan by constantly reading it until EOF.
func ReaderToChan(stream chan<- []byte, reader io.Reader) error {
	for {
		buffer := make([]byte, bufferSize)
		nBytes, err := reader.Read(buffer)
		if nBytes > 0 {
			stream <- buffer[:nBytes]
		}
		if err != nil {
			return err
		}
	}
}

// ChanToWriter dumps all content from a given chan to a writer until the chan is closed.
func ChanToWriter(writer io.Writer, stream <-chan []byte) error {
	for buffer := range stream {
		_, err := writer.Write(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}
