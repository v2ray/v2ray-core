package io

import (
	"io"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
)

// Pipe dumps all content from reader to writer, until an error happens.
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

// PipeUntilEOF behaves the same as Pipe(). The only difference is PipeUntilEOF returns nil on EOF.
func PipeUntilEOF(reader Reader, writer Writer) error {
	err := Pipe(reader, writer)
	if err != nil && errors.Cause(err) != io.EOF {
		return err
	}
	return nil
}
