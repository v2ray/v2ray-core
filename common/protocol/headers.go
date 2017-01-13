package protocol

import (
	"runtime"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/uuid"
)

type RequestCommand byte

const (
	RequestCommandTCP = RequestCommand(0x01)
	RequestCommandUDP = RequestCommand(0x02)
)

// RequestOption is the options of a request.
type RequestOption byte

const (
	// RequestOptionChunkStream indicates request payload is chunked. Each chunk consists of length, authentication and payload.
	RequestOptionChunkStream = RequestOption(0x01)

	// RequestOptionConnectionReuse indicates client side expects to reuse the connection.
	RequestOptionConnectionReuse = RequestOption(0x02)

	// RequestOptionCompressedStream indicates request payload is compressed.
	RequestOptionCompressedStream = RequestOption(0x04)
)

func (v RequestOption) Has(option RequestOption) bool {
	return (v & option) == option
}

func (v *RequestOption) Set(option RequestOption) {
	*v = (*v | option)
}

func (v *RequestOption) Clear(option RequestOption) {
	*v = (*v & (^option))
}

type Security byte

func (v Security) Is(t SecurityType) bool {
	return v == Security(t)
}

func NormSecurity(s Security) Security {
	if s.Is(SecurityType_UNKNOWN) {
		return Security(SecurityType_LEGACY)
	}
	return s
}

type RequestHeader struct {
	Version  byte
	Command  RequestCommand
	Option   RequestOption
	Security Security
	Port     net.Port
	Address  net.Address
	User     *User
}

func (v *RequestHeader) Destination() net.Destination {
	if v.Command == RequestCommandUDP {
		return net.UDPDestination(v.Address, v.Port)
	}
	return net.TCPDestination(v.Address, v.Port)
}

type ResponseOption byte

const (
	ResponseOptionConnectionReuse = ResponseOption(0x01)
)

func (v *ResponseOption) Set(option ResponseOption) {
	*v = (*v | option)
}

func (v ResponseOption) Has(option ResponseOption) bool {
	return (v & option) == option
}

func (v *ResponseOption) Clear(option ResponseOption) {
	*v = (*v & (^option))
}

type ResponseCommand interface{}

type ResponseHeader struct {
	Option  ResponseOption
	Command ResponseCommand
}

type CommandSwitchAccount struct {
	Host     net.Address
	Port     net.Port
	ID       *uuid.UUID
	AlterIds uint16
	Level    uint32
	ValidMin byte
}

func (v *SecurityConfig) AsSecurity() Security {
	if v == nil {
		return Security(SecurityType_LEGACY)
	}
	if v.Type == SecurityType_AUTO {
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "s390x" {
			return Security(SecurityType_AES128_GCM)
		}
		return Security(SecurityType_CHACHA20_POLY1305)
	}
	return NormSecurity(Security(v.Type))
}
