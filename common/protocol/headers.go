package protocol

import (
	"runtime"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/uuid"
)

type RequestCommand byte

const (
	RequestCommandTCP = RequestCommand(0x01)
	RequestCommandUDP = RequestCommand(0x02)
)

type RequestOption byte

const (
	RequestOptionChunkStream     = RequestOption(0x01)
	RequestOptionConnectionReuse = RequestOption(0x02)
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
	User     *User
	Command  RequestCommand
	Option   RequestOption
	Security Security
	Address  v2net.Address
	Port     v2net.Port
}

func (v *RequestHeader) Destination() v2net.Destination {
	if v.Command == RequestCommandUDP {
		return v2net.UDPDestination(v.Address, v.Port)
	}
	return v2net.TCPDestination(v.Address, v.Port)
}

type ResponseOption byte

const (
	ResponseOptionConnectionReuse = ResponseOption(1)
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
	Host     v2net.Address
	Port     v2net.Port
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
