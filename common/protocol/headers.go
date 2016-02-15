package protocol

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type RequestCommand byte

const (
	RequestCommandTCP = RequestCommand(0x01)
	RequestCommandUDP = RequestCommand(0x02)
)

type RequestOption byte

const (
	RequestOptionChunkStream = RequestOption(0x01)
)

type RequestHeader struct {
	Version byte
	User    *User
	Command RequestCommand
	Option  RequestOption
	Address v2net.Address
	Port    v2net.Port
}

type ResponseCommand interface{}

type ResponseHeader struct {
	Command ResponseCommand
}
