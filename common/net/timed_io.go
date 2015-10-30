package net

import (
	"net"
	"time"
)

var (
	emptyTime time.Time
)

type TimeOutReader struct {
	timeout    int
	connection net.Conn
}

func NewTimeOutReader(timeout int, connection net.Conn) *TimeOutReader {
	return &TimeOutReader{
		timeout:    timeout,
		connection: connection,
	}
}

func (reader *TimeOutReader) Read(p []byte) (n int, err error) {
	if reader.timeout > 0 {
		deadline := time.Duration(reader.timeout) * time.Second
		reader.connection.SetReadDeadline(time.Now().Add(deadline))
	}
	n, err = reader.connection.Read(p)
	if reader.timeout > 0 {
		reader.connection.SetReadDeadline(emptyTime)
	}
	return
}

func (reader *TimeOutReader) GetTimeOut() int {
	return reader.timeout
}

func (reader *TimeOutReader) SetTimeOut(value int) {
	reader.timeout = value
}
