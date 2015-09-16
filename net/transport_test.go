package net

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestReaderAndWrite(t *testing.T) {
	assert := unit.Assert(t)

	size := 1024 * 1024
	buffer := make([]byte, size)
	nBytes, err := rand.Read(buffer)
	assert.Int(nBytes).Equals(len(buffer))
	assert.Error(err).IsNil()

	readerBuffer := bytes.NewReader(buffer)
	writerBuffer := bytes.NewBuffer(make([]byte, 0, size))

	transportChan := make(chan []byte, size/bufferSize*10)

	err = ReaderToChan(transportChan, readerBuffer)
	assert.Error(err).Equals(io.EOF)
	close(transportChan)

	err = ChanToWriter(writerBuffer, transportChan)
	assert.Error(err).IsNil()

	assert.Bytes(buffer).Equals(writerBuffer.Bytes())
}

type StaticReader struct {
	total   int
	current int
}

func (reader *StaticReader) Read(b []byte) (size int, err error) {
	size = len(b)
	if size > reader.total-reader.current {
		size = reader.total - reader.current
	}
	//rand.Read(b[:size])
	reader.current += size
	if reader.current == reader.total {
		err = io.EOF
	}
	return
}

func BenchmarkTransport(b *testing.B) {
	size := 1024 * 1024

	for i := 0; i < b.N; i++ {
		transportChanA := make(chan []byte, 128)
		transportChanB := make(chan []byte, 128)

		readerA := &StaticReader{size, 0}
		readerB := &StaticReader{size, 0}

		writerA := ioutil.Discard
		writerB := ioutil.Discard

		finishA := make(chan bool)
		finishB := make(chan bool)

		go func() {
			ChanToWriter(writerA, transportChanA)
			close(finishA)
		}()

		go func() {
			ReaderToChan(transportChanA, readerA)
			close(transportChanA)
		}()

		go func() {
			ChanToWriter(writerB, transportChanB)
			close(finishB)
		}()

		go func() {
			ReaderToChan(transportChanB, readerB)
			close(transportChanB)
		}()

		<-transportChanA
		<-transportChanB
	}
}
