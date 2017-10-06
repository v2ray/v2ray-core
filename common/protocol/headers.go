package protocol

import (
	"runtime"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/uuid"
)

// RequestCommand is a custom command in a proxy request.
type RequestCommand byte

const (
	RequestCommandTCP = RequestCommand(0x01)
	RequestCommandUDP = RequestCommand(0x02)
	RequestCommandMux = RequestCommand(0x03)
)

func (c RequestCommand) TransferType() TransferType {
	if c == RequestCommandTCP {
		return TransferTypeStream
	}

	return TransferTypePacket
}

// RequestOption is the options of a request.
type RequestOption byte

const (
	// RequestOptionChunkStream indicates request payload is chunked. Each chunk consists of length, authentication and payload.
	RequestOptionChunkStream = RequestOption(0x01)

	// RequestOptionConnectionReuse indicates client side expects to reuse the connection.
	RequestOptionConnectionReuse = RequestOption(0x02)

	RequestOptionChunkMasking = RequestOption(0x04)
)

func (o RequestOption) Has(option RequestOption) bool {
	return (o & option) == option
}

func (o *RequestOption) Set(option RequestOption) {
	*o = (*o | option)
}

func (o *RequestOption) Clear(option RequestOption) {
	*o = (*o & (^option))
}

type Security byte

func (s Security) Is(t SecurityType) bool {
	return s == Security(t)
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

func (h *RequestHeader) Destination() net.Destination {
	if h.Command == RequestCommandUDP {
		return net.UDPDestination(h.Address, h.Port)
	}
	return net.TCPDestination(h.Address, h.Port)
}

type ResponseOption byte

const (
	ResponseOptionConnectionReuse = ResponseOption(0x01)
)

func (o *ResponseOption) Set(option ResponseOption) {
	*o = (*o | option)
}

func (o ResponseOption) Has(option ResponseOption) bool {
	return (o & option) == option
}

func (o *ResponseOption) Clear(option ResponseOption) {
	*o = (*o & (^option))
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

func (sc *SecurityConfig) AsSecurity() Security {
	if sc == nil {
		return Security(SecurityType_LEGACY)
	}
	if sc.Type == SecurityType_AUTO {
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "s390x" {
			return Security(SecurityType_AES128_GCM)
		}
		return Security(SecurityType_CHACHA20_POLY1305)
	}
	return NormSecurity(Security(sc.Type))
}
