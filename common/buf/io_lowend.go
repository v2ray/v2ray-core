// +build mips mipsle mips64 mips64le arm arm64

package buf

import (
	"io"
)

func newReaderPlatform(reader io.Reader) Reader {
	return &SingleReader{
		Reader: reader,
	}
}
