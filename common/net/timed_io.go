package net

import (
	"io"
	"net"
	"time"
)

var (
	emptyTime time.Time
)

type TimeOutReader struct {
	timeout    int
	connection net.Conn
	worker     io.Reader
}

func NewTimeOutReader(timeout int, connection net.Conn) *TimeOutReader {
	reader := &TimeOutReader{
		connection: connection,
		timeout:    -100,
	}
	reader.SetTimeOut(timeout)
	return reader
}

func (reader *TimeOutReader) Read(p []byte) (int, error) {
	return reader.worker.Read(p)
}

func (reader *TimeOutReader) GetTimeOut() int {
	return reader.timeout
}

func (reader *TimeOutReader) SetTimeOut(value int) {
	if value == reader.timeout {
		return
	}
	reader.timeout = value
	if value > 0 {
		reader.worker = &timedReaderWorker{
			timeout:    value,
			connection: reader.connection,
		}
	} else {
		reader.worker = &noOpReaderWorker{
			connection: reader.connection,
		}
	}
}

type timedReaderWorker struct {
	timeout    int
	connection net.Conn
}

func (this *timedReaderWorker) Read(p []byte) (int, error) {
	deadline := time.Duration(this.timeout) * time.Second
	this.connection.SetReadDeadline(time.Now().Add(deadline))
	nBytes, err := this.connection.Read(p)
	this.connection.SetReadDeadline(emptyTime)
	return nBytes, err
}

type noOpReaderWorker struct {
	connection net.Conn
}

func (this *noOpReaderWorker) Read(p []byte) (int, error) {
	return this.connection.Read(p)
}
