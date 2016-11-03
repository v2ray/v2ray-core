package http

import (
	"bytes"
	"io"
	"net"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/transport/internet"
)

const (
	CRLF   = "\r\n"
	ENDING = CRLF + CRLF
)

type HeaderReader struct {
}

func (*HeaderReader) Read(reader io.Reader) (*alloc.Buffer, error) {
	buffer := alloc.NewLocalBuffer(2048)
	for {
		_, err := buffer.FillFrom(reader)
		if err != nil {
			return nil, err
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
	return buffer, nil
}

type HeaderWriter struct {
	header *alloc.Buffer
}

func (this *HeaderWriter) Write(writer io.Writer) error {
	if this.header == nil {
		return nil
	}
	_, err := writer.Write(this.header.Value)
	this.header.Release()
	this.header = nil
	return err
}

type HttpConn struct {
	net.Conn

	readBuffer    *alloc.Buffer
	oneTimeReader *HeaderReader
	oneTimeWriter *HeaderWriter
}

func NewHttpConn(conn net.Conn, reader *HeaderReader, writer *HeaderWriter) *HttpConn {
	return &HttpConn{
		Conn:          conn,
		oneTimeReader: reader,
		oneTimeWriter: writer,
	}
}

func (this *HttpConn) Read(b []byte) (int, error) {
	if this.oneTimeReader != nil {
		buffer, err := this.oneTimeReader.Read(this.Conn)
		if err != nil {
			return 0, err
		}
		this.readBuffer = buffer
		this.oneTimeReader = nil
	}

	if this.readBuffer.Len() > 0 {
		nBytes, err := this.readBuffer.Read(b)
		if nBytes == this.readBuffer.Len() {
			this.readBuffer.Release()
			this.readBuffer = nil
		}
		return nBytes, err
	}

	return this.Conn.Read(b)
}

func (this *HttpConn) Write(b []byte) (int, error) {
	if this.oneTimeWriter != nil {
		err := this.oneTimeWriter.Write(this.Conn)
		this.oneTimeWriter = nil
		if err != nil {
			return 0, err
		}
	}

	return this.Conn.Write(b)
}

type HttpAuthenticator struct {
	config *Config
}

func (this HttpAuthenticator) GetClientWriter() *HeaderWriter {
	header := alloc.NewLocalBuffer(2048)
	config := this.config.Request
	header.AppendString(config.Method.GetValue()).AppendString(" ").AppendString(config.PickUri()).AppendString(" ").AppendString(config.GetFullVersion()).AppendString(CRLF)

	headers := config.PickHeaders()
	for _, h := range headers {
		header.AppendString(h).AppendString(CRLF)
	}
	header.AppendString(CRLF)
	return &HeaderWriter{
		header: header,
	}
}

func (this HttpAuthenticator) GetServerWriter() *HeaderWriter {
	header := alloc.NewLocalBuffer(2048)
	config := this.config.Response
	header.AppendString(config.GetFullVersion()).AppendString(" ").AppendString(config.Status.GetCode()).AppendString(" ").AppendString(config.Status.GetReason()).AppendString(CRLF)

	headers := config.PickHeaders()
	for _, h := range headers {
		header.AppendString(h).AppendString(CRLF)
	}
	header.AppendString(CRLF)
	return &HeaderWriter{
		header: header,
	}
}

func (this HttpAuthenticator) Client(conn net.Conn) net.Conn {
	if this.config.Request == nil && this.config.Response == nil {
		return conn
	}
	return NewHttpConn(conn, new(HeaderReader), this.GetClientWriter())
}

func (this HttpAuthenticator) Server(conn net.Conn) net.Conn {
	if this.config.Request == nil && this.config.Response == nil {
		return conn
	}
	return NewHttpConn(conn, new(HeaderReader), this.GetServerWriter())
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
