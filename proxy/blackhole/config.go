package blackhole

import (
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

const (
	http403response = `HTTP/1.1 403 Forbidden
Connection: close
Cache-Control: max-age=3600, public
Content-Length: 0


`
)

type ResponseConfig interface {
	WriteTo(buf.Writer)
}

func (v *NoneResponse) WriteTo(buf.Writer) {}

func (v *HTTPResponse) WriteTo(writer buf.Writer) {
	b := buf.NewLocal(512)
	b.AppendSupplier(serial.WriteString(http403response))
	writer.Write(b)
}

func (v *Config) GetInternalResponse() (ResponseConfig, error) {
	if v.GetResponse() == nil {
		return new(NoneResponse), nil
	}

	config, err := v.GetResponse().GetInstance()
	if err != nil {
		return nil, err
	}
	return config.(ResponseConfig), nil
}
