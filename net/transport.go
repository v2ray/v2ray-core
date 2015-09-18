package net

import (
	"io"
)

const (
	bufferSize = 32 * 1024
)

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
	return nil
}

func ChanToWriter(writer io.Writer, stream <-chan []byte) error {
	for buffer := range stream {
		_, err := writer.Write(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}
