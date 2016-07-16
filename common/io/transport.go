package io

import (
	"github.com/v2ray/v2ray-core/common/log"
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
