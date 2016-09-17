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

func (this RequestOption) Has(option RequestOption) bool {
	return (this & option) == option
}

func (this *RequestOption) Set(option RequestOption) {
	*this = (*this | option)
}

func (this *RequestOption) Clear(option RequestOption) {
	*this = (*this & (^option))
}

type RequestHeader struct {
	Version byte
	User    *User
	Command RequestCommand
	Option  RequestOption
	Address v2net.Address
	Port    v2net.Port
}

func (this *RequestHeader) Destination() v2net.Destination {
	if this.Command == RequestCommandUDP {
		return v2net.UDPDestination(this.Address, this.Port)
	}
	return v2net.TCPDestination(this.Address, this.Port)
}

type ResponseOption byte

const (
	ResponseOptionConnectionReuse = ResponseOption(1)
)

func (this *ResponseOption) Set(option ResponseOption) {
	*this = (*this | option)
}

func (this ResponseOption) Has(option ResponseOption) bool {
	return (this | option) == option
}

func (this *ResponseOption) Clear(option ResponseOption) {
	*this = (*this & (^option))
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
