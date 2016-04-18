package io

func Pipe(reader Reader, writer Writer) error {
	for {
		buffer, err := reader.Read()
		if buffer.Len() > 0 {
			err = writer.Write(buffer)
		} else {
			buffer.Release()
		}

		if err != nil {
			return nil
		}
	}
}
