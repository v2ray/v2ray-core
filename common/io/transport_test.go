package io_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	. "github.com/v2ray/v2ray-core/common/io"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestReaderAndWrite(t *testing.T) {
	v2testing.Current(t)

	size := 1024 * 1024
	buffer := make([]byte, size)
	nBytes, err := rand.Read(buffer)
	assert.Int(nBytes).Equals(len(buffer))
	assert.Error(err).IsNil()

	readerBuffer := bytes.NewReader(buffer)
	writerBuffer := bytes.NewBuffer(make([]byte, 0, size))

	transportChan := make(chan *alloc.Buffer, 1024)

	err = ReaderToChan(transportChan, NewAdaptiveReader(readerBuffer))
	assert.Error(err).Equals(io.EOF)
	close(transportChan)

	err = ChanToWriter(writerBuffer, transportChan)
	assert.Error(err).IsNil()

	assert.Bytes(buffer).Equals(writerBuffer.Bytes())
}
