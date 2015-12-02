package net_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
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

	err = v2net.ReaderToChan(transportChan, readerBuffer)
	assert.Error(err).Equals(io.EOF)
	close(transportChan)

	err = v2net.ChanToWriter(writerBuffer, transportChan)
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
		v2net.ChanToWriter(writerA, transportChanA)
		close(finishA)
	}()

	go func() {
		v2net.ReaderToChan(transportChanA, readerA)
		close(transportChanA)
	}()

	go func() {
		v2net.ChanToWriter(writerB, transportChanB)
		close(finishB)
	}()

	go func() {
		v2net.ReaderToChan(transportChanB, readerB)
		close(transportChanB)
	}()

	<-transportChanA
	<-transportChanB
}
