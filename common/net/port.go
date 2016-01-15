package net

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

type Port serial.Uint16Literal

func PortFromBytes(port []byte) Port {
	return Port(uint16(port[0])<<8 + uint16(port[1]))
}

func (this Port) Value() uint16 {
	return uint16(this)
}

func (this Port) Bytes() []byte {
	return []byte{byte(this >> 8), byte(this)}
}

func (this Port) String() string {
	return serial.Uint16Literal(this).String()
}

type PortRange struct {
	From Port
	To   Port
}
