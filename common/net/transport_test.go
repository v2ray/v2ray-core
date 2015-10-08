package net

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
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

	transportChan := make(chan *alloc.Buffer, 1024)

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
	for i := 0; i < size; i++ {
		b[i] = byte(i)
	}
	//rand.Read(b[:size])
	reader.current += size
	if reader.current == reader.total {
		err = io.EOF
	}
	return
}

func BenchmarkTransport1K(b *testing.B) {
	size := 1 * 1024

	for i := 0; i < b.N; i++ {
		runBenchmarkTransport(size)
	}
}

func BenchmarkTransport2K(b *testing.B) {
	size := 2 * 1024

	for i := 0; i < b.N; i++ {
		runBenchmarkTransport(size)
	}
}

func BenchmarkTransport4K(b *testing.B) {
	size := 4 * 1024

	for i := 0; i < b.N; i++ {
		runBenchmarkTransport(size)
	}
}

func BenchmarkTransport10K(b *testing.B) {
	size := 10 * 1024

	for i := 0; i < b.N; i++ {
		runBenchmarkTransport(size)
	}
}

func BenchmarkTransport100K(b *testing.B) {
	size := 100 * 1024

	for i := 0; i < b.N; i++ {
		runBenchmarkTransport(size)
	}
}

func BenchmarkTransport1M(b *testing.B) {
	size := 1024 * 1024

	for i := 0; i < b.N; i++ {
		runBenchmarkTransport(size)
	}
}

func BenchmarkTransport10M(b *testing.B) {
	size := 10 * 1024 * 1024

	for i := 0; i < b.N; i++ {
		runBenchmarkTransport(size)
	}
}

func runBenchmarkTransport(size int) {

	transportChanA := make(chan *alloc.Buffer, 16)
	transportChanB := make(chan *alloc.Buffer, 16)

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
