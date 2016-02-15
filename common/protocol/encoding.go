package protocol

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
)

type RequestEncoder interface {
	EncodeHeader(*RequestHeader) *alloc.Buffer
	EncodeBody(io.Writer) io.Writer
}

type RequestDecoder interface {
	DecodeHeader(io.Reader) *RequestHeader
	DecodeBody(io.Reader) io.Reader
}

type ResponseEncoder interface {
	EncodeHeader(*ResponseHeader) *alloc.Buffer
	EncodeBody(io.Writer) io.Writer
}

type ResponseDecoder interface {
	DecodeHeader(io.Reader) *ResponseHeader
	DecodeBody(io.Reader) io.Reader
}
