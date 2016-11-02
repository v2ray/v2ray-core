package http

import (
	"bytes"
	"net"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/transport/internet"
)

const (
	CRLF   = "\r\n"
	ENDING = CRLF + CRLF
)

type HttpConn struct {
	net.Conn

	buffer     *alloc.Buffer
	readHeader bool

	writeHeaderContent *alloc.Buffer
	writeHeader        bool
}

func NewHttpConn(conn net.Conn, writeHeaderContent *alloc.Buffer) *HttpConn {
	return &HttpConn{
		Conn:               conn,
		readHeader:         true,
		writeHeader:        true,
		writeHeaderContent: writeHeaderContent,
	}
}

func (this *HttpConn) Read(b []byte) (int, error) {
	if this.readHeader {
		buffer := alloc.NewLocalBuffer(2048)
		for {
			_, err := buffer.FillFrom(this.Conn)
			if err != nil {
				return 0, err
			}
			if n := bytes.Index(buffer.Value, []byte(ENDING)); n != -1 {
				buffer.SliceFrom(n + len(ENDING))
				break
			}
			if buffer.Len() >= len(ENDING) {
				copy(buffer.Value, buffer.Value[buffer.Len()-len(ENDING):])
				buffer.Slice(0, len(ENDING))
			}
		}
		this.buffer = buffer
		this.readHeader = false
	}

	if this.buffer.Len() > 0 {
		nBytes, err := this.buffer.Read(b)
		if nBytes == this.buffer.Len() {
			this.buffer.Release()
			this.buffer = nil
		}
		return nBytes, err
	}

	return this.Conn.Read(b)
}

func (this *HttpConn) Write(b []byte) (int, error) {
	if this.writeHeader {
		_, err := this.Conn.Write(this.writeHeaderContent.Value)
		this.writeHeaderContent.Release()
		if err != nil {
			return 0, err
		}
		this.writeHeader = false
	}

	return this.Conn.Write(b)
}

type HttpAuthenticator struct {
	config *Config
}

func (this HttpAuthenticator) GetClientWriteHeader() *alloc.Buffer {
	header := alloc.NewLocalBuffer(2048)
	config := this.config.Request
	header.AppendString(config.Method.GetValue()).AppendString(" ").AppendString(config.PickUri()).AppendString(" ").AppendString(config.GetFullVersion()).AppendString(CRLF)

	headers := config.PickHeaders()
	for _, h := range headers {
		header.AppendString(h).AppendString(CRLF)
	}
	header.AppendString(CRLF)
	return header
}

func (this HttpAuthenticator) GetServerWriteHeader() *alloc.Buffer {
	header := alloc.NewLocalBuffer(2048)
	config := this.config.Response
	header.AppendString(config.GetFullVersion()).AppendString(" ").AppendString(config.Status.GetCode()).AppendString(" ").AppendString(config.Status.GetReason()).AppendString(CRLF)

	headers := config.PickHeaders()
	for _, h := range headers {
		header.AppendString(h).AppendString(CRLF)
	}
	header.AppendString(CRLF)
	return header
}

func (this HttpAuthenticator) Client(conn net.Conn) net.Conn {
	return NewHttpConn(conn, this.GetClientWriteHeader())
}

func (this HttpAuthenticator) Server(conn net.Conn) net.Conn {
	return NewHttpConn(conn, this.GetServerWriteHeader())
}

type HttpAuthenticatorFactory struct{}

func (HttpAuthenticatorFactory) Create(config interface{}) internet.ConnectionAuthenticator {
	return HttpAuthenticator{
		config: config.(*Config),
	}
}

func init() {
	internet.RegisterConnectionAuthenticator(loader.GetType(new(Config)), HttpAuthenticatorFactory{})
}
