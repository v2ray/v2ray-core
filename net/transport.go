package net

import (
	"io"
)

const (
	bufferSize = 32 * 1024
)

var (
	dirtyBuffers = make(chan []byte, 1024)
)

func getBuffer() []byte {
	var buffer []byte
	select {
	case buffer = <-dirtyBuffers:
	default:
		buffer = make([]byte, bufferSize)
	}
	return buffer
}

func putBuffer(buffer []byte) {
	select {
	case dirtyBuffers <- buffer:
	default:
	}
}

func ReaderToChan(stream chan<- []byte, reader io.Reader) error {
	for {
    buffer := make([]byte, bufferSize)
		//buffer := getBuffer()
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
		//putBuffer(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}
