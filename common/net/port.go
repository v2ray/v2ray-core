package net

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

type Port serial.Uint16Literal

func PortFromBytes(port []byte) Port {
	return Port(serial.ParseUint16(port))
}

func (this Port) Value() uint16 {
	return uint16(this)
}

func (this Port) Bytes() []byte {
	return serial.Uint16Literal(this).Bytes()
}

func (this Port) String() string {
	return serial.Uint16Literal(this).String()
}

type PortRange struct {
	From Port
	To   Port
}

func (this PortRange) Contains(port Port) bool {
	return this.From <= port && port <= this.To
}
