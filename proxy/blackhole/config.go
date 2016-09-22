package blackhole

import (
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"strings"
	"v2ray.com/core/common/loader"
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

func (this *Response) GetInternalResponse() (ResponseConfig, error) {
	if this == nil {
		return new(NoneResponse), nil
	}

	var r ResponseConfig
	switch this.Type {
	case Response_None:
		r = new(NoneResponse)
	case Response_HTTP:
		r = new(HTTPResponse)
	}
	err := ptypes.UnmarshalAny(this.Settings, r.(proto.Message))
	if err != nil {
		return nil, err
	}
	return r, nil
}

var (
	cache = loader.ConfigCreatorCache{}
)

func init() {
	cache.RegisterCreator(strings.ToLower(Response_Type_name[int32(Response_None)]), func() interface{} { return new(NoneResponse) })
	cache.RegisterCreator(strings.ToLower(Response_Type_name[int32(Response_HTTP)]), func() interface{} { return new(HTTPResponse) })
}
