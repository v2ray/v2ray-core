package net

import (
	"io"
)

func ReadFrom(reader io.Reader, sizeInKilo int) ([]byte, error) {
	buffer := make([]byte, sizeInKilo<<10)
	nBytes, err := reader.Read(buffer)
	if nBytes == 0 {
		return nil, err
	}
	return buffer[:nBytes], err
}

// ReaderToChan dumps all content from a given reader to a chan by constantly reading it until EOF.
func ReaderToChan(stream chan<- []byte, reader io.Reader) error {
	bufferSizeKilo := 4
	for {
		data, err := ReadFrom(reader, bufferSizeKilo)
		if len(data) > 0 {
			stream <- data
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
