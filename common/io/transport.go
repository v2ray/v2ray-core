package io

func Pipe(reader Reader, writer Writer) error {
	for {
		buffer, err := reader.Read()
		if err != nil {
			return err
		}

		if buffer.IsEmpty() {
			buffer.Release()
			continue
		}

		err = writer.Write(buffer)
		if err != nil {
			return err
		}
	}
}
