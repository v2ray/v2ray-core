package net

import (
	"io"

	"github.com/v2ray/v2ray-core/log"
)

const (
	bufferSize = 8192
)

func ReaderToChan(stream chan<- []byte, reader io.Reader) error {
	for {
		buffer := make([]byte, bufferSize)
		nBytes, err := reader.Read(buffer)
		if err != nil {
			return err
		}
		stream <- buffer[:nBytes]
	}
	return nil
}

func ChanToWriter(writer io.Writer, stream <-chan []byte) error {
	for buffer := range stream {
		nBytes, err := writer.Write(buffer)
		log.Debug("Writing %d bytes with error %v", nBytes, err)
		if err != nil {
			return err
		}
	}
	return nil
}
