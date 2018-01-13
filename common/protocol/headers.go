package protocol

import (
	"runtime"

	"v2ray.com/core/common/bitmask"
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

const (
	// RequestOptionChunkStream indicates request payload is chunked. Each chunk consists of length, authentication and payload.
	RequestOptionChunkStream bitmask.Byte = 0x01

	// RequestOptionConnectionReuse indicates client side expects to reuse the connection.
	RequestOptionConnectionReuse bitmask.Byte = 0x02

	RequestOptionChunkMasking bitmask.Byte = 0x04
)

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
	Option   bitmask.Byte
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

const (
	ResponseOptionConnectionReuse bitmask.Byte = 0x01
)

type ResponseCommand interface{}

type ResponseHeader struct {
	Option  bitmask.Byte
	Command ResponseCommand
}

type CommandSwitchAccount struct {
	Host     net.Address
	Port     net.Port
	ID       *uuid.UUID
	Level    uint32
	AlterIds uint16
	ValidMin byte
}

func (sc *SecurityConfig) AsSecurity() Security {
	if sc == nil || sc.Type == SecurityType_AUTO {
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64" || runtime.GOARCH == "s390x" {
			return Security(SecurityType_AES128_GCM)
		}
		return Security(SecurityType_CHACHA20_POLY1305)
	}
	return NormSecurity(Security(sc.Type))
}

func IsDomainTooLong(domain string) bool {
	return len(domain) > 256
}
