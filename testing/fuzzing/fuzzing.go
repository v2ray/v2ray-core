package fuzzing

import (
	"bytes"
	"crypto/rand"
	"io"
)

func RandomBytes() []byte {
	buffer := make([]byte, 256)
	rand.Read(buffer)
	return buffer[1 : 1+int(buffer[0])]
}

func RandomReader() io.Reader {
	return bytes.NewReader(RandomBytes())
}
