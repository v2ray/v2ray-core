package net

import (
	"io"
	"net"
	"time"

	"v2ray.com/core/common"
)

var (
	emptyTime time.Time
)

type TimeOutReader struct {
	timeout    uint32
	connection net.Conn
	worker     io.Reader
}

func NewTimeOutReader(timeout uint32 /* seconds */, connection net.Conn) *TimeOutReader {
	reader := &TimeOutReader{
		connection: connection,
		timeout:    0,
	}
	reader.SetTimeOut(timeout)
	return reader
}

func (reader *TimeOutReader) Read(p []byte) (int, error) {
	return reader.worker.Read(p)
}

func (reader *TimeOutReader) GetTimeOut() uint32 {
	return reader.timeout
}

func (reader *TimeOutReader) SetTimeOut(value uint32) {
	if reader.worker != nil && value == reader.timeout {
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

func (reader *TimeOutReader) Release() {
	common.Release(reader.connection)
	common.Release(reader.worker)
}

type timedReaderWorker struct {
	timeout    uint32
	connection net.Conn
}

func (v *timedReaderWorker) Read(p []byte) (int, error) {
	deadline := time.Duration(v.timeout) * time.Second
	v.connection.SetReadDeadline(time.Now().Add(deadline))
	nBytes, err := v.connection.Read(p)
	v.connection.SetReadDeadline(emptyTime)
	return nBytes, err
}

type noOpReaderWorker struct {
	connection net.Conn
}

func (v *noOpReaderWorker) Read(p []byte) (int, error) {
	return v.connection.Read(p)
}
