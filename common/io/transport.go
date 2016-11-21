package io

import (
	"io"
	"v2ray.com/core/common/log"
)

func Pipe(reader Reader, writer Writer) error {
	for {
		buffer, err := reader.Read()
		if err != nil {
			log.Debug("IO: Pipe exits as ", err)
			return err
		}

		if buffer.IsEmpty() {
			buffer.Release()
			continue
		}

		err = writer.Write(buffer)
		if err != nil {
			log.Debug("IO: Pipe exits as ", err)
			buffer.Release()
			return err
		}
	}
}

func PipeUntilEOF(reader Reader, writer Writer) error {
	err := Pipe(reader, writer)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
