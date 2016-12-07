package protocol

import (
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
	return (v | option) == option
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
