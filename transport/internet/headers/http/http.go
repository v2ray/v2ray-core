package http

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

const (
	CRLF   = "\r\n"
	ENDING = CRLF + CRLF
)

var (
	writeCRLF = serial.WriteString(CRLF)
)

type Reader interface {
	Read(io.Reader) (*buf.Buffer, error)
}

type Writer interface {
	Write(io.Writer) error
}

type NoOpReader struct{}

func (v *NoOpReader) Read(io.Reader) (*buf.Buffer, error) {
	return nil, nil
}

type NoOpWriter struct{}

func (v *NoOpWriter) Write(io.Writer) error {
	return nil
}

type HeaderReader struct {
}

func (*HeaderReader) Read(reader io.Reader) (*buf.Buffer, error) {
	buffer := buf.NewSmall()
	for {
		err := buffer.AppendSupplier(buf.ReadFrom(reader))
		if err != nil {
			return nil, err
		}
		if n := bytes.Index(buffer.Bytes(), []byte(ENDING)); n != -1 {
			buffer.SliceFrom(n + len(ENDING))
			break
		}
		if buffer.Len() >= len(ENDING) {
			copy(buffer.Bytes(), buffer.BytesFrom(buffer.Len()-len(ENDING)))
			buffer.Slice(0, len(ENDING))
		}
	}
	if buffer.IsEmpty() {
		buffer.Release()
		return nil, nil
	}
	return buffer, nil
}

type HeaderWriter struct {
	header *buf.Buffer
}

func NewHeaderWriter(header *buf.Buffer) *HeaderWriter {
	return &HeaderWriter{
		header: header,
	}
}

func (v *HeaderWriter) Write(writer io.Writer) error {
	if v.header == nil {
		return nil
	}
	_, err := writer.Write(v.header.Bytes())
	v.header.Release()
	v.header = nil
	return err
}

type HttpConn struct {
	net.Conn

	readBuffer    *buf.Buffer
	oneTimeReader Reader
	oneTimeWriter Writer
}

func NewHttpConn(conn net.Conn, reader Reader, writer Writer) *HttpConn {
	return &HttpConn{
		Conn:          conn,
		oneTimeReader: reader,
		oneTimeWriter: writer,
	}
}

func (v *HttpConn) Read(b []byte) (int, error) {
	if v.oneTimeReader != nil {
		buffer, err := v.oneTimeReader.Read(v.Conn)
		if err != nil {
			return 0, err
		}
		v.readBuffer = buffer
		v.oneTimeReader = nil
	}

	if v.readBuffer.Len() > 0 {
		nBytes, err := v.readBuffer.Read(b)
		if nBytes == v.readBuffer.Len() {
			v.readBuffer.Release()
			v.readBuffer = nil
		}
		return nBytes, err
	}

	return v.Conn.Read(b)
}

func (v *HttpConn) Write(b []byte) (int, error) {
	if v.oneTimeWriter != nil {
		err := v.oneTimeWriter.Write(v.Conn)
		v.oneTimeWriter = nil
		if err != nil {
			return 0, err
		}
	}

	return v.Conn.Write(b)
}

type HttpAuthenticator struct {
	config *Config
}

func (v HttpAuthenticator) GetClientWriter() *HeaderWriter {
	header := buf.NewSmall()
	config := v.config.Request
	header.AppendSupplier(serial.WriteString(strings.Join([]string{config.Method.GetValue(), config.PickUri(), config.GetFullVersion()}, " ")))
	header.AppendSupplier(writeCRLF)

	headers := config.PickHeaders()
	for _, h := range headers {
		header.AppendSupplier(serial.WriteString(h))
		header.AppendSupplier(writeCRLF)
	}
	header.AppendSupplier(writeCRLF)
	return &HeaderWriter{
		header: header,
	}
}

func (v HttpAuthenticator) GetServerWriter() *HeaderWriter {
	header := buf.NewSmall()
	config := v.config.Response
	header.AppendSupplier(serial.WriteString(strings.Join([]string{config.GetFullVersion(), config.Status.GetCode(), config.Status.GetReason()}, " ")))
	header.AppendSupplier(writeCRLF)

	headers := config.PickHeaders()
	for _, h := range headers {
		header.AppendSupplier(serial.WriteString(h))
		header.AppendSupplier(writeCRLF)
	}
	if !config.HasHeader("Date") {
		header.AppendSupplier(serial.WriteString("Date: "))
		header.AppendSupplier(serial.WriteString(time.Now().Format(http.TimeFormat)))
		header.AppendSupplier(writeCRLF)
	}
	header.AppendSupplier(writeCRLF)
	return &HeaderWriter{
		header: header,
	}
}

func (v HttpAuthenticator) Client(conn net.Conn) net.Conn {
	if v.config.Request == nil && v.config.Response == nil {
		return conn
	}
	var reader Reader = new(NoOpReader)
	if v.config.Request != nil {
		reader = new(HeaderReader)
	}

	var writer Writer = new(NoOpWriter)
	if v.config.Response != nil {
		writer = v.GetClientWriter()
	}
	return NewHttpConn(conn, reader, writer)
}

func (v HttpAuthenticator) Server(conn net.Conn) net.Conn {
	if v.config.Request == nil && v.config.Response == nil {
		return conn
	}
	return NewHttpConn(conn, new(HeaderReader), v.GetServerWriter())
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
