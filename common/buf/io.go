package buf

import "io"

func ReadFrom(reader io.Reader) Supplier {
	return func(b []byte) (int, error) {
		return reader.Read(b)
	}
}

func ReadFullFrom(reader io.Reader, size int) Supplier {
	return func(b []byte) (int, error) {
		return io.ReadFull(reader, b[:size])
	}
}
