package net

import (
	"strconv"
)

type Port uint16

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
	return strconv.Itoa(int(this))
}
