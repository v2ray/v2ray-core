package http

import (
	"bytes"
	"io"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/transport/internet"
)

const (
	CRLF   = "\r\n"
	ENDING = CRLF + CRLF
)

type RequestAuthenticator struct {
	config *RequestConfig
}

func (this *RequestAuthenticator) Seal(writer io.Writer) io.Writer {
	header := alloc.NewLocalBuffer(2048)
	header.AppendString(this.config.Method).AppendString(" ").AppendString(this.config.PickUri()).AppendString(" ").AppendString(this.config.GetVersion()).AppendString(CRLF)

	headers := this.config.PickHeaders()
	for _, h := range headers {
		header.AppendString(h).AppendString(CRLF)
	}
	header.AppendString(CRLF)

	writer.Write(header.Value)
	header.Release()

	return writer
}

func (this *RequestAuthenticator) Open(reader io.Reader) (io.Reader, error) {
	buffer := alloc.NewLocalBuffer(2048)
	for {
		_, err := buffer.FillFrom(reader)
		if err != nil {
			return nil, err
		}
		if n := bytes.Index(buffer.Value, []byte(ENDING)); n != -1 {
			buffer.SliceFrom(n + len(ENDING))
			return &BufferAndReader{
				buffer: buffer,
				reader: reader,
			}, nil
		}
		if buffer.Len() >= len(ENDING) {
			copy(buffer.Value, buffer.Value[buffer.Len()-len(ENDING):])
			buffer.Slice(0, len(ENDING))
		}
	}
}

type BufferAndReader struct {
	buffer *alloc.Buffer
	reader io.Reader
}

func (this *BufferAndReader) Read(b []byte) (int, error) {
	if this.buffer.Len() == 0 {
		return this.reader.Read(b)
	}
	n, err := this.buffer.Read(b)
	if n == this.buffer.Len() {
		this.buffer.Release()
		this.buffer = nil
	}
	return n, err
}

type RequestAuthenticatorFactory struct{}

func (RequestAuthenticatorFactory) Create(config interface{}) internet.ConnectionAuthenticator {
	return &RequestAuthenticator{
		config: config.(*RequestConfig),
	}
}

type ResponseAuthenticator struct {
	config *ResponseConfig
}

func (this *ResponseAuthenticator) Seal(writer io.Writer) io.Writer {
	header := alloc.NewLocalBuffer(2048)
	header.AppendString(this.config.GetVersion()).AppendString(" ").AppendString(this.config.Status).AppendString(" ").AppendString(this.config.Reason).AppendString(CRLF)

	headers := this.config.PickHeaders()
	for _, h := range headers {
		header.AppendString(h).AppendString(CRLF)
	}
	header.AppendString(CRLF)

	writer.Write(header.Value)
	header.Release()

	return writer
}

func (this *ResponseAuthenticator) Open(reader io.Reader) (io.Reader, error) {
	buffer := alloc.NewLocalBuffer(2048)
	for {
		_, err := buffer.FillFrom(reader)
		if err != nil {
			return nil, err
		}
		if n := bytes.Index(buffer.Value, []byte(ENDING)); n != -1 {
			buffer.SliceFrom(n + len(ENDING))
			return &BufferAndReader{
				buffer: buffer,
				reader: reader,
			}, nil
		}
		if buffer.Len() >= len(ENDING) {
			copy(buffer.Value, buffer.Value[buffer.Len()-len(ENDING):])
			buffer.Slice(0, len(ENDING))
		}
	}
}

type ResponseAuthenticatorFactory struct{}

func (ResponseAuthenticatorFactory) Create(config interface{}) internet.ConnectionAuthenticator {
	return &ResponseAuthenticator{
		config: config.(*ResponseConfig),
	}
}

func init() {
	internet.RegisterConnectionAuthenticator(loader.GetType(new(RequestConfig)), RequestAuthenticatorFactory{})
	internet.RegisterConnectionAuthenticator(loader.GetType(new(ResponseConfig)), ResponseAuthenticatorFactory{})
}
