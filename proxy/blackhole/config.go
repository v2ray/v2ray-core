package blackhole

import (
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

const (
	http403response = `HTTP/1.1 403 Forbidden
Connection: close
Cache-Control: max-age=3600, public
Content-Length: 0


`
)

type ResponseConfig interface {
	AsAny() *any.Any
	WriteTo(v2io.Writer)
}

func (this *NoneResponse) WriteTo(v2io.Writer) {}

func (this *NoneResponse) AsAny() *any.Any {
	r, _ := ptypes.MarshalAny(this)
	return r
}

func (this *HTTPResponse) WriteTo(writer v2io.Writer) {
	writer.Write(alloc.NewLocalBuffer(512).Clear().AppendString(http403response))
}

func (this *HTTPResponse) AsAny() *any.Any {
	r, _ := ptypes.MarshalAny(this)
	return r
}

func (this *Config) GetInternalResponse() (ResponseConfig, error) {
	if this.GetResponse() == nil {
		return new(NoneResponse), nil
	}

	config, err := this.GetResponse().GetInstance()
	if err != nil {
		return nil, err
	}
	return config.(ResponseConfig), nil
}
