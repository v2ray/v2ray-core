package blackhole

import (
	"v2ray.com/core/common"
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

// ResponseConfig is the configuration for blackhole responses.
type ResponseConfig interface {
	// WriteTo writes predefined response to the give buffer.
	WriteTo(buf.Writer)
}

// WriteTo implements ResponseConfig.WriteTo().
func (*NoneResponse) WriteTo(buf.Writer) {}

// WriteTo implements ResponseConfig.WriteTo().
func (*HTTPResponse) WriteTo(writer buf.Writer) {
	b := buf.New()
	common.Must(b.AppendSupplier(serial.WriteString(http403response)))
	writer.WriteMultiBuffer(buf.NewMultiBufferValue(b))
}

// GetInternalResponse converts response settings from proto to internal data structure.
func (c *Config) GetInternalResponse() (ResponseConfig, error) {
	if c.GetResponse() == nil {
		return new(NoneResponse), nil
	}

	config, err := c.GetResponse().GetInstance()
	if err != nil {
		return nil, err
	}
	return config.(ResponseConfig), nil
}
