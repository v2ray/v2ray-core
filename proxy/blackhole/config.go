package blackhole

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
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
	WriteTo(buf.Writer) int32
}

// WriteTo implements ResponseConfig.WriteTo().
func (*NoneResponse) WriteTo(buf.Writer) int32 { return 0 }

// WriteTo implements ResponseConfig.WriteTo().
func (*HTTPResponse) WriteTo(writer buf.Writer) int32 {
	b := buf.New()
	common.Must2(b.WriteString(http403response))
	n := b.Len()
	writer.WriteMultiBuffer(buf.MultiBuffer{b})
	return n
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
